# Book2Ru 📚

`book2ru` is a powerful and efficient CLI tool for translating large text files from English to Russian using the
OpenRouter API. Available in both Go and Ruby implementations, it processes text through a standard STDIN/STDOUT
pipeline, making it a versatile component in any text processing workflow.

## Go Version (Prototype) 🚀

`book2ru-go` is a Go implementation prototype that provides the same functionality as the Ruby version with improved
performance and easier deployment.

### Prerequisites for Go Version

- **Go** (version 1.21 or later)

### Installation and Compilation 🔧

1. **Clone the repository:**
   ```bash
   git clone <repository_url>
   cd book2ru
   ```

2. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

3. **Compile the application:**
   ```bash
   go build -o book2ru-go .
   ```

### Configuration for Go Version ⚙️

1. **API Key**: Create a `.env` file in the project root and add your OpenRouter API key:
   ```env
   # .env
   OPENROUTER_KEY='sk-or-v1-...'
   ```

2. **Custom Settings (Optional)**: Create a `.book2ru-go.yml` file to override default settings:
   ```yaml
   # .book2ru-go.yml
   model: "anthropic/claude-3-5-sonnet-20240620"
   prompt: "Translate the following English text to Russian. Return only the translated text, preserving all original formatting and line breaks:"
   ```

### Usage for Go Version 📖

```bash
# Basic translation
./book2ru-go < my_book.txt > my_book_ru.txt

# Using a different model
./book2ru-go --model "mistralai/mistral-large-latest" < input.txt > output.txt

# Show version
./book2ru-go --version

# Show help
./book2ru-go --help
```

### Testing Go Version 🧪

Run the comprehensive test suite:

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run specific test
go test -run TestCreateBatchesFromContent
```

### Go Version Command-Line Options

| Option               | Alias | Description                               |
| -------------------- | ----- | ----------------------------------------- |
| `--help`             | `-h`  | Show the help message.                    |
| `--version`          |       | Display version and model information.    |
| `--model <MODEL>`    | `-m`  | Specify the LLM model to use.             |
| `--openrouter_key <KEY>` | `-o`  | Provide the OpenRouter API key directly.  |

---

## Ruby Version (Original) 💎

The original Ruby implementation provides the same functionality with Ruby's expressiveness and ecosystem.

### Prerequisites for Ruby Version

- **Ruby** (version `3.4.2` is recommended)
- **Bundler** for managing gems

### Installation for Ruby Version 🚀

1. **Install dependencies:**
   ```bash
   bundle install
   ```

2. **Make the script executable:**
   ```bash
   chmod +x book2ru.rb
   ```

### Configuration for Ruby Version ⚙️

1. **API Key**: Create a `.env` file in the project root and add your OpenRouter API key:

   ```env
   # .env
   OPENROUTER_KEY='sk-or-v1-...'
   ```

2. **Custom Settings (Optional)**: Create a `.book2ru.yml` file to override default settings like the model or
   translation prompt.

   ```yaml
   # .book2ru.yml
   model: "anthropic/claude-3-5-sonnet-20240620"
   prompt: "Translate the following English text to Russian. Return only the translated text, preserving all original formatting and line breaks:"
   ```

### Usage for Ruby Version 📖

The tool reads from STDIN and writes the translated output to STDOUT.

#### Basic Translation

```bash
./book2ru.rb < my_book.txt > my_book_ru.txt
```

#### Using a Different Model

You can specify a different model via the command line:

```bash
./book2ru.rb --model "mistralai/mistral-large-latest" < input.txt > output.txt
```

### Testing Ruby Version 🧪

To ensure the application is working correctly, you can run the test suite:

1. Make sure all development dependencies are installed:
   ```bash
   bundle install
   ```

2. Run the RSpec tests:
   ```bash
   bundle exec rspec
   ```

### Ruby Version Command-Line Options

| Option               |  Alias | Description                                    |
| -------------------- | ------------------------------------------------------- |
| `--help`             |  `-h`  | Show the help      message.                    |
| `--version`          |  | Display versio           n and model information.    |
| `--model <MODEL>`       | `-m`  | Specify the LL   M model to use.             |
| `--openrouter_key <KEY> ` | `-o`  | Provide the Op enRouter API key directly.  |

---

## Key Features ✨

- **Large File Processing**: Handles large text files (e.g., books) with ease using a stream-based approach.
- **Intelligent Batching**: Automatically splits input text into 10KB chunks to work efficiently with API limits.
- **Robust Error Handling**: Features an automatic retry mechanism with exponential backoff for network and API errors.
- **Flexible Configuration**: Configure API keys, models, and prompts via a `.env` file, configuration files, or
  command-line arguments.
- **Metadata Output**: Outputs progress and debugging information to STDERR, keeping STDOUT clean for the translated
  text.
- **Cross-Platform**: Available in both Go (fast, single binary) and Ruby (expressive, ecosystem-rich) implementations.

## Changing the Target Language 🌐

By default, both tools are configured to translate to Russian. You can change the target language by modifying the
translation prompt.

1. Create or open the configuration file (`.book2ru-go.yml` for Go version, `.book2ru.yml` for Ruby version) in the
   project root.
2. Set the `prompt` key with instructions for the new target language. For example, to translate to Spanish:

   ```yaml
   prompt: "Translate this text to Spanish. Only return the translated text, nothing else:"
   ```

The model will now follow the new instruction for translation.

## How It Works

1. The application reads the entire text content from STDIN.
2. It splits the content into batches of approximately 10KB, ensuring that lines are not broken in the middle.
3. Each batch is sent to the OpenRouter API for translation.
4. The translated text is immediately written to STDOUT, preserving the original structure.
5. All progress, status, and API responses are logged to STDERR, so they don't interfere with the output file.
