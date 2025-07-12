# book2ru Project Structure

## Overview

`book2ru` is a Ruby CLI application designed to translate English text to Russian using OpenRouter API. The project follows a simple, single-file architecture optimized for text processing with intelligent batching and retry mechanisms.

## Core Architecture

### Main Application File

**`book2ru.rb`** - The primary executable file containing the complete application logic:

- **Book2ru Class**: Main application class handling the entire translation workflow
- **Configuration Management**: Loads settings from YAML files and environment variables
- **CLI Interface**: Command-line argument parsing with OptionParser
- **Batch Processing**: Intelligent text batching system (10KB chunks)
- **API Integration**: OpenRouter API client with retry logic and error handling
- **Stream Processing**: STDIN/STDOUT pipeline for large file processing

### Key Methods Structure

- `initialize()` - Sets up default configuration and loads external settings
- `run(argv)` - Main entry point handling command-line arguments
- `translate_from_stdin()` - Core translation pipeline
- `create_batches_from_content()` - Text batching logic
- `translate_text_batch()` - OpenRouter API interaction
- `process_batch()` - Individual batch processing with metadata

## Configuration System

### Primary Configuration File

**`.book2ru.yml`** - YAML configuration file containing:

- `model`: OpenRouter model specification (default: "google/gemini-flash-1.5")
- `prompt`: Translation prompt template
- `metadata_footer`: Debug information output control
- `retry_attempts`: API failure retry count

### Environment Configuration

**`.env`** - Environment variables file (not tracked in git):

- `OPENROUTER_KEY`: Required API key for OpenRouter service

## Dependencies Management

### Gemfile Structure

**`Gemfile`** - Ruby dependencies specification:

- **Production Dependencies**:
  - `dotenv`: Environment variable management
  
- **Development Dependencies**:
  - `rspec`: Testing framework
  - `rubocop`: Code style enforcement

**`Gemfile.lock`** - Locked dependency versions ensuring consistent builds

## Project Data

### Source Data Directory

**`_src/`** - Contains source and output files:

- Input text files for translation
- Generated translation outputs
- Book content in various formats (TXT, EPUB, PDF)

## CLI Interface

### Command-Line Options

The application supports the following command-line interface:

- `--help` / `-h`: Display usage information
- `--version`: Show version and model information
- `--model MODEL` / `-m`: Override default OpenRouter model
- `--openrouter_key KEY` / `-o`: Provide API key via command line
- `--rate-limits RPM` / `-r`: Set API rate limiting

### Usage Pattern

```bash
book2ru [options] < input.txt > output.txt
```

## Processing Pipeline

### Input Processing

1. **Content Reading**: Reads entire STDIN content into memory
2. **Batch Creation**: Splits content into 10KB chunks while preserving line boundaries
3. **Batch Processing**: Processes each batch sequentially through OpenRouter API
4. **Output Generation**: Streams translated content directly to STDOUT

### Error Handling

- **API Failures**: Exponential backoff retry mechanism (3 attempts)
- **Network Timeouts**: Automatic retry with progressive delays
- **Configuration Errors**: Graceful error messages and exit codes

## Development Environment

### Ruby Version

- **Required**: Ruby 3.4.2 (specified in Gemfile)
- **Manager**: rbenv for version management

### Testing Structure

- **Framework**: RSpec for unit testing
- **Configuration**: `.rspec` file for test runner settings
- **Helper**: `spec/spec_helper.rb` for test environment setup

## File Organization

### Root Level Files

- `book2ru.rb` - Main application executable
- `.book2ru.yml` - Application configuration
- `Gemfile` - Ruby dependencies
- `README.md` - Project documentation
- `.gitignore` - Git ignore patterns
- `.rspec` - RSpec configuration

### Directory Structure

```
book2ru/
├── book2ru.rb          # Main application
├── .book2ru.yml        # Configuration
├── Gemfile            # Dependencies
├── _src/              # Source files
└── ai-docs/           # Documentation
```

## Configuration Loading Priority

1. **Default Configuration**: Hardcoded defaults in `Book2ru#initialize`
2. **YAML Configuration**: `.book2ru.yml` file overrides
3. **Environment Variables**: `.env` file loaded via dotenv
4. **Command Line Arguments**: Highest priority overrides

## Key Design Decisions

### Single-File Architecture

The entire application is contained in a single Ruby file (`book2ru.rb`) for simplicity and ease of deployment. This design choice enables:

- Easy distribution as a single executable
- Minimal dependency management
- Straightforward deployment in various environments

### Batch Processing Strategy

Text is processed in 10KB batches to:

- Optimize API usage and costs
- Handle large files efficiently
- Maintain reasonable memory usage
- Preserve text structure and formatting

### Stream-Based I/O

The application uses STDIN/STDOUT for file processing to:

- Support Unix pipeline patterns
- Handle arbitrarily large files
- Maintain low memory footprint
- Enable flexible file processing workflows 