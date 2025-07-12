require 'spec_helper'
require_relative '../book2ru'
require 'tempfile'

RSpec.describe Book2ru do
  # Helper to run the app with args and capture output
  def run_app(args, stdin_str = '')
    app = described_class.new
    original_stdin = $stdin
    $stdin = StringIO.new(stdin_str)
    stdout, stderr = capture_output { app.run(args) }
    $stdin = original_stdin
    [stdout, stderr]
  end

  describe 'CLI options' do
    it 'shows version with --version' do
      stdout, = run_app(['--version'])
      expect(stdout).to include("book2ru v#{Book2ru::VERSION}")
    end

    it 'shows help with -h' do
      stdout, = run_app(['-h'])
      expect(stdout).to include('Usage: book2ru [options]')
    end

    it 'shows help with --help' do
      stdout, = run_app(['--help'])
      expect(stdout).to include('Usage: book2ru [options]')
    end
  end

  describe 'configuration' do
    it 'loads config from .book2ru.yml' do
      config_data = { 'model' => 'test-model', 'prompt' => 'test prompt' }
      temp_config = Tempfile.new('.book2ru.yml')
      temp_config.write(config_data.to_yaml)
      temp_config.close

      # Stub Dir.pwd to make the app find the temp config file
      allow(Dir).to receive(:pwd).and_return(File.dirname(temp_config.path))
      allow(File).to receive(:exist?).and_call_original
      allow(File).to receive(:exist?).with('.book2ru.yml').and_return(true)
      allow(YAML).to receive(:load_file).with('.book2ru.yml').and_return(config_data)

      local_app = described_class.new
      expect(local_app.instance_variable_get(:@config)[:model]).to eq('test-model')
      expect(local_app.instance_variable_get(:@config)[:prompt]).to eq('test prompt')
    ensure
      temp_config.unlink
    end
  end

  describe 'translation flow' do
    let(:api_key) { 'test-key' }
    let(:input_text) { 'Hello world' }
    let(:translated_text) { 'Привет, мир' }

    before do
      ENV['OPENROUTER_KEY'] = api_key
      stub_request(:post, 'https://openrouter.ai/api/v1/chat/completions')
        .to_return(
          status: 200,
          body: { choices: [{ message: { content: translated_text } }] }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
    end

    it 'translates text from stdin' do
      stdout, stderr = run_app([], input_text)
      expect(stdout).to eq(translated_text)
      expect(stderr).to include("starting translation")
      expect(stderr).to include("Translated by book2ru")
    end

    context 'when API key is missing' do
      before do
        ENV.delete('OPENROUTER_KEY')
      end

      it 'returns an error and exits' do
        stdout, = run_app([], input_text)
        expect(stdout).to include('[ERROR] No API key provided.')
      end
    end

    context 'when API call fails' do
      it 'retries and then fails' do
        stub_request(:post, 'https://openrouter.ai/api/v1/chat/completions')
          .to_return({ status: 500 }, { status: 500 }, { status: 500 })

        stdout, stderr = run_app([], input_text)

        expect(stdout).to include('Failed to translate after 3 attempts')
        expect(stderr).to include('Retry attempt 1 failed')
        expect(stderr).to include('Retry attempt 2 failed')
      end
    end
  end

  describe '#create_batches_from_content' do
    let(:instance) { described_class.new }
    let(:long_line) { 'a' * 9000 }

    it 'returns an empty array for empty content' do
      batches = instance.send(:create_batches_from_content, '')
      expect(batches).to be_empty
    end

    it 'creates a single batch for content smaller than BATCH_SIZE' do
      content = 'This is a small text.'
      batches = instance.send(:create_batches_from_content, content)
      expect(batches.size).to eq(1)
      expect(batches.first[:content]).to eq(content)
    end

    it 'splits content larger than BATCH_SIZE into multiple batches' do
      content = "#{long_line}\n#{long_line}"
      batches = instance.send(:create_batches_from_content, content)
      expect(batches.size).to eq(2)
      expect(batches[0][:content]).to eq("#{long_line}\n")
      expect(batches[1][:content]).to eq(long_line)
    end

    it 'correctly handles content ending with a newline' do
      content = "#{long_line}\n#{long_line}\n"
      batches = instance.send(:create_batches_from_content, content)
      expect(batches.size).to eq(2)
      expect(batches[0][:content]).to eq("#{long_line}\n")
      expect(batches[1][:content]).to eq("#{long_line}\n")
    end
  end
end 