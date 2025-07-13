# frozen_string_literal: true

source 'https://rubygems.org'

ruby '3.4.2'

# CLI application dependencies
gem 'dotenv'

# Testing dependencies
group :test do
  gem 'rspec', '~> 3.12'
  gem 'rspec-collection_matchers'
  gem 'rspec-core'
  gem 'rspec-expectations'
  gem 'rspec-its', '~> 1.3'
  gem 'rspec-mocks'
  gem 'rspec-support'
  gem 'rubocop-rspec'
  gem 'simplecov', '~> 0.22.0', require: false
  gem 'webmock', '~> 3.19'
end

group :development, :test do
  gem 'brakeman'
  gem 'bundle-audit'
  gem 'rubocop'
end
