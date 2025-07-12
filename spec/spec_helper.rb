require 'webmock/rspec'

require 'simplecov'
SimpleCov.start do
  add_filter '/spec/' # исключаем тесты из отчета
  add_filter '/vendor/'
  enable_coverage :branch # включаем анализ покрытия веток
  minimum_coverage ENV.fetch('MIN_COVERAGE', 90).to_i

  # Добавляем группы файлов
  add_group 'Library', 'lib'
  add_group 'CLI', 'lib/runner/cli.rb'
  add_group 'Runner', 'lib/runner/runner.rb'
end

# See http://rubydoc.info/gems/rspec-core/RSpec/Core/Configuration
RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups
  config.filter_run_when_matching :focus
  config.example_status_persistence_file_path = "spec/examples.txt"
  config.disable_monkey_patching!
  config.warnings = true

  if config.files_to_run.one?
    config.default_formatter = "doc"
  end

  config.profile_examples = 10
  config.order = :random
  Kernel.srand config.seed

  # Clean up env variables and prevent Dotenv loading
  config.before(:each) do
    allow(Dotenv).to receive(:load)
    ENV.delete('OPENROUTER_KEY')
  end
end

# Helper to capture stdout/stderr
def capture_output
  original_stdout = $stdout
  original_stderr = $stderr
  $stdout = StringIO.new
  $stderr = StringIO.new
  yield
  [$stdout.string, $stderr.string]
ensure
  $stdout = original_stdout
  $stderr = original_stderr
end
