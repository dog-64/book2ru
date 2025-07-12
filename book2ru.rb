#!/usr/bin/env ruby
# CLI-приложение для перевода текстов через OpenRouter API

require 'bundler/setup'
require 'optparse'
require 'yaml'
require 'net/http'
require 'json'
require 'dotenv'

class Book2ru
  VERSION = '0.1.0'.freeze
  BATCH_SIZE = 10_000 # Размер батча в байтах

  def initialize
    @config = {
      model: 'google/gemini-flash-1.5',
      prompt: 'Translate this text to Russian. Only return the translated text, nothing else:',
      metadata_footer: true,
      rate_limits: nil,
      retry_attempts: 3
    }

    load_config
    load_environment_files
  end

  def run(argv)
    options = parse_options(argv)

    case options[:action]
    when :help
      puts option_parser.help
      0
    when :version
      puts "book2ru v#{VERSION}"
      puts "Using model: #{@config[:model]}"
      0
    when :translate
      validate_api_key
      translate_from_stdin
      0
    else
      puts 'Unknown action. Use --help for usage information.'
      1
    end
  rescue => e
    puts "[ERROR] #{e.message}"
    1
  end

  private

  def parse_options(argv)
    options = { action: :translate }

    option_parser.parse!(argv)

    # Проверяем, была ли установлена другая акция через флаги
    if @action
      options[:action] = @action
      @action = nil # Сбрасываем для следующего использования
    end

    options
  end

  def option_parser
    @option_parser ||= OptionParser.new do |opts|
      opts.banner = 'Usage: book2ru [options] < input.txt > output.txt'
      opts.separator ''
      opts.separator 'Options:'

      opts.on('-h', '--help', 'Show this help message') do
        @action = :help
      end

      opts.on('--version', 'Show version information') do
        @action = :version
      end

      opts.on('-m', '--model MODEL', 'Specify LLM model') do |model|
        @config[:model] = model
      end

      opts.on('-o', '--openrouter_key KEY', 'OpenRouter API key') do |key|
        ENV['OPENROUTER_KEY'] = key
      end

      opts.on('-r', '--rate-limits RPM', Integer, 'Rate limits (requests per minute)') do |rpm|
        @config[:rate_limits] = rpm
      end

      opts.separator ''
      opts.separator 'Examples:'
      opts.separator '  book2ru < input.txt > output-ru.txt'
      opts.separator '  book2ru --model claude-3-haiku --rate-limits 10 < input.txt > output.txt'
    end
  end

  def load_config
    config_file = '.book2ru.yml'
    return unless File.exist?(config_file)

    file_config = YAML.load_file(config_file)
    @config.merge!(file_config.transform_keys(&:to_sym)) if file_config
  end

  def load_environment_files
    # Загружаем .env файл если он существует
    Dotenv.load if File.exist?('.env')
  end

  def validate_api_key
    return if ENV['OPENROUTER_KEY']

    raise 'No API key provided. Set OPENROUTER_KEY environment variable.'
  end

  # Основной метод для перевода из stdin
  def translate_from_stdin
    warn "# book2ru v#{VERSION} - starting translation using #{@config[:model]}" if @config[:metadata_footer]

    # Читаем все содержимое и создаем батчи
    content = $stdin.read
    batches = create_batches_from_content(content)

    warn "# Created #{batches.length} batches from #{content.bytesize} bytes" if @config[:metadata_footer]

    # Обрабатываем каждый батч
    batches.each_with_index do |batch, index|
      process_batch(batch, index + 1, batches.length)
    end

    return unless @config[:metadata_footer]

    warn "# Translated by book2ru v#{VERSION} using model #{@config[:model]}"
  end

  # Создает батчи из содержимого, разбивая по 10KB
  def create_batches_from_content(content)
    batches = []
    lines = content.split("\n", -1) # -1 чтобы сохранить пустые строки в конце

    current_batch = []
    current_size = 0

    lines.each_with_index do |line, index|
      line_with_newline = index == lines.length - 1 ? line : "#{line}\n"
      line_size = line_with_newline.bytesize

      # Если добавление этой строки превысит лимит - завершаем текущий батч
      if current_size + line_size > BATCH_SIZE && current_batch.any?
        batches << {
          content: current_batch.join,
          size: current_size,
          lines: current_batch.length
        }
        current_batch = [line_with_newline]
        current_size = line_size
      else
        current_batch << line_with_newline
        current_size += line_size
      end
    end

    # Добавляем последний батч
    if current_batch.any?
      batches << {
        content: current_batch.join,
        size: current_size,
        lines: current_batch.length
      }
    end

    batches
  end

  # Обрабатывает один батч
  def process_batch(batch, batch_number, total_batches)
    if @config[:metadata_footer]
      warn "# Processing batch #{batch_number}/#{total_batches} (#{batch[:size]} bytes, #{batch[:lines]} lines)"
    end

    translated = translate_text_batch(batch[:content])
    print translated
  end

  # Переводит текстовый батч через OpenRouter API
  def translate_text_batch(text)
    attempts = 0

    begin
      attempts += 1

      # Создаём HTTP запрос к OpenRouter API
      uri = URI('https://openrouter.ai/api/v1/chat/completions')
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = true

      request = Net::HTTP::Post.new(uri)
      request['Authorization'] = "Bearer #{ENV.fetch('OPENROUTER_KEY', nil)}"
      request['Content-Type'] = 'application/json'
      request['HTTP-Referer'] = 'https://github.com/book2ru'
      request['X-Title'] = 'book2ru CLI'
      request['OpenAI-Organization'] = 'openrouter'
      request['User-Agent'] = 'book2ru/0.1.0'

      request.body = {
        model: @config[:model],
        messages: [
          {
            role: 'user',
            content: "#{@config[:prompt]}\n\n#{text}"
          }
        ],
        temperature: 0.1,
        max_tokens: nil,
        stream: false
      }.to_json

      response = http.request(request)

      raise "OpenRouter API error: #{response.code} - #{response.body}" unless response.code == '200'

      parsed = JSON.parse(response.body)
      content = parsed['choices'][0]['message']['content']

      # Отладочный вывод в stderr
      if @config[:metadata_footer]
        warn "# API Response: #{content[0..100]}#{'...' if content.length > 100}"
      end

      content
    rescue => e
      unless attempts < @config[:retry_attempts]
        raise "Failed to translate after #{@config[:retry_attempts]} attempts: #{e.message}"
      end

      warn "# Retry attempt #{attempts} failed: #{e.message}"
      sleep(2**attempts) # Экспоненциальная задержка
      retry
    end
  end
end

# Запуск приложения
if __FILE__ == $PROGRAM_NAME
  app = Book2ru.new
  exit_code = app.run(ARGV)
  exit(exit_code)
end
