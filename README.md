# Book2Ru 📚

`book2ru` is a powerful and efficient CLI tool for translating large text files from English to Russian using the  OpenRouter API. Available in both Go and Ruby implementations, it processes text through a standard STDIN/STDOUT  pipeline, making it a versatile component in any text processing workflow.

## Version Comparison 📊

| Feature                     | Go Version         | Ruby Version            |
|-----------------------------|--------------------|-------------------------|
| **Performance**             | ⚡ Faster           | 🐢 Slower               |
| **Deployment**              | 📦 Single binary   | 💎 Requires Ruby + gems |
| **Resume after failures**   | ✅ `--start-batch`  | ❌ Not supported         |
| **Rate limiting**           | ❌ Not supported    | ✅ `--rate-limits`       |
| **Configuration file**      | `.book2ru-go.yml`  | `.book2ru.yml`          |
| **Retry attempts**          | ✅ 3 (configurable) | ✅ 3 (configurable)      |
| **All error types retried** | ✅ Yes              | ✅ Yes                   |
| **Ecosystem**               | 🔧 Minimal         | 🌟 Rich Ruby gems       |

### Which Version to Choose? 🤔

- **Choose Go version** if you need:
    - Maximum performance and speed
    - Single binary deployment
    - Resume after failures (`--start-batch`)
    - Minimal dependencies

- **Choose Ruby version** if you need:
    - Rate limiting features
    - Ruby ecosystem integration
    - More readable/customizable code
    - Don't mind Ruby dependencies

## Go Version (Prototype) 🚀

`book2ru-go` is a Go implementation prototype that provides the same functionality as the Ruby version with improved  performance and easier deployment.

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

2. **Custom Settings (Optional)**: Copy the sample configuration file and customize:
   ```bash
   cp .book2ru-go.yml.sample .book2ru-go.yml
   # Edit .book2ru-go.yml with your preferences
   ```

   Example configuration:
   ```yaml
   # .book2ru-go.yml
   model: "anthropic/claude-3-5-sonnet-20240620"
   prompt: "Translate the following English text to Russian. Return only the translated text, preserving all original formatting and line breaks:"
   batch_size: 15000  # Larger batches for faster processing
   retry_attempts: 5  # More retries for unstable connections
   ```

### Usage for Go Version 📖

```bash
# Basic translation
./book2ru-go < my_book.txt > my_book_ru.txt

# Using a different model
./book2ru-go --model "mistralai/mistral-large-latest" < input.txt > output.txt

# Resume translation from batch 10 (after a failure)
./book2ru-go --start-batch 10 < my_book.txt > my_book_ru.txt

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

| Option                   | Alias | Description                                              |
|--------------------------|-------|----------------------------------------------------------|
| `--help`                 | `-h`  | Show the help message.                                   |
| `--version`              |       | Display version and model information.                   |
| `--model <MODEL>`        | `-m`  | Specify the LLM model to use.                            |
| `--openrouter_key <KEY>` | `-o`  | Provide the OpenRouter API key directly.                 |
| `--start-batch <N>`      |       | Resume translation from batch N (useful after failures). |
| `--batch-size <SIZE>`    |       | Set batch size in bytes (default: 10000).               |

### Go Version Configuration File Options

The `.book2ru-go.yml` file supports these options:

```yaml
# Model configuration
model: "google/gemini-flash-1.5"                    # LLM model to use
prompt: "Translate this text to Russian..."         # Translation prompt

# Processing options  
start_batch: 1                                      # Starting batch number
batch_size: 10000                                   # Batch size in bytes
retry_attempts: 3                                   # Number of retry attempts
metadata_footer: true                               # Show progress information
```

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

2. **Custom Settings (Optional)**: Copy the sample configuration file and customize:

   ```bash
   cp .book2ru.yml.sample .book2ru.yml
   # Edit .book2ru.yml with your preferences
   ```

   Example configuration:
   ```yaml
   # .book2ru.yml
   model: "anthropic/claude-3-5-sonnet-20240620"
   prompt: "Translate the following English text to Russian. Return only the translated text, preserving all original formatting and line breaks:"
   batch_size: 15000      # Larger batches for faster processing
   retry_attempts: 5      # More retries for unstable connections
   rate_limits: 30        # Limit to 30 requests per minute
   ```

### Usage for Ruby Version 📖

The tool reads from STDIN and writes the translated output to STDOUT.

```bash
# Basic translation
./book2ru.rb < my_book.txt > my_book_ru.txt

# Using a different model
./book2ru.rb --model "mistralai/mistral-large-latest" < input.txt > output.txt

# Set rate limits (requests per minute)
./book2ru.rb --rate-limits 10 < input.txt > output.txt

# Provide API key directly
./book2ru.rb --openrouter_key "sk-or-v1-..." < input.txt > output.txt

# Show version
./book2ru.rb --version

# Show help
./book2ru.rb --help
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

| Option                   | Alias | Description                              |
|--------------------------|-------|------------------------------------------|
| `--help`                 | `-h`  | Show the help message.                   |
| `--version`              |       | Display version and model information.   |
| `--model <MODEL>`        | `-m`  | Specify the LLM model to use.            |
| `--openrouter_key <KEY>` | `-o`  | Provide the OpenRouter API key directly. |
| `--rate-limits <RPM>`    | `-r`  | Set rate limits (requests per minute).   |
| `--batch-size <SIZE>`    | `-b`  | Set batch size in bytes (default: 10000). |

### Ruby Version Configuration File Options

The `.book2ru.yml` file supports these options:

```yaml
# Model configuration
model: "google/gemini-flash-1.5"                    # LLM model to use
prompt: "Translate this text to Russian..."         # Translation prompt

# Processing options
batch_size: 10000                                   # Batch size in bytes
retry_attempts: 3                                   # Number of retry attempts
metadata_footer: true                               # Show progress information
rate_limits: null                                   # Rate limits (requests per minute)
```

**Note**: Ruby version does not support `--start-batch` option. Use the Go version for resuming after failures.

---

## Key Features ✨

### Common Features (Both Versions)

- **Large File Processing**: Handles large text files (e.g., books) with ease using a stream-based approach.
- **Intelligent Batching**: Automatically splits input text into 10KB chunks to work efficiently with API limits.
- **Robust Error Handling**: Features an automatic retry mechanism with exponential backoff for network and API errors (
  3 attempts by default).
- **Flexible Configuration**: Configure API keys, models, and prompts via a `.env` file, configuration files, or
  command-line arguments.
- **Metadata Output**: Outputs progress and debugging information to STDERR, keeping STDOUT clean for the translated
  text.

### Go Version Exclusive Features

- **Resume After Failures**: Can resume translation from a specific batch number after interruptions using
  `--start-batch` option.
- **Helpful Error Messages**: When translation fails, shows exact command to resume from the failed batch.
- **Single Binary**: Compiles to a single executable file for easy deployment.
- **Better Performance**: Generally faster execution and lower memory usage.

### Ruby Version Exclusive Features

- **Rate Limiting**: Built-in rate limiting support with `--rate-limits` option.
- **Ruby Ecosystem**: Access to the rich Ruby gem ecosystem for extensions.
- **Expressive Syntax**: More readable and maintainable code for customization.

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

## Working with Different File Formats 📄

### Converting PDF to Text

Use the included `pdf2text_clean.py` script to extract text from PDF files:

```bash
# Install required dependency
pip install PyMuPDF

# Convert PDF to text
cat "book.pdf" | python3 pdf2text_clean.py > book.txt
# or
python3 pdf2text_clean.py < "book.pdf" > "book.txt"

# Then translate
./book2ru-go < book.txt > book_ru.txt
```

### Converting EPUB to Text

For EPUB files, you can use various tools:

```bash
# Using pandoc (recommended)
pandoc book.epub -t plain -o book.txt

# Using calibre
ebook-convert book.epub book.txt

# Then translate
./book2ru-go < book.txt > book_ru.txt
```

### Converting Text Back to EPUB/PDF

After translation, you can convert back to structured formats:

```bash
# Text to EPUB using pandoc
pandoc book_ru.txt -o book_ru.epub --metadata title="Translated Book"

# Text to PDF using pandoc
pandoc book_ru.txt -o book_ru.pdf

# For better formatting, create a markdown file first
echo "# Translated Book" > book_ru.md
echo "" >> book_ru.md
cat book_ru.txt >> book_ru.md
pandoc book_ru.md -o book_ru.epub
```

### Complete Workflow Example

```bash
# 1. Extract text from PDF
python3 pdf2text_clean.py < original_book.pdf > original_text.txt

# 2. Translate text
./book2ru-go < original_text.txt > translated_text.txt

# 3. Convert back to EPUB
pandoc translated_text.txt -o translated_book.epub --metadata title="Переведенная книга"

# 4. Or convert to PDF
pandoc translated_text.txt -o translated_book.pdf
```

## Error Handling and Recovery 🔄

### Automatic Retry (Both Versions)

Both Go and Ruby versions automatically retry failed requests up to 3 times with exponential backoff:

- 1st retry: 2 seconds delay
- 2nd retry: 4 seconds delay
- 3rd retry: 8 seconds delay

All types of errors are retried, including network timeouts, API errors, and JSON parsing failures.

### Resume After Failures (Go Version Only)

If translation fails after all retries, you can resume from where it left off using the `--start-batch` option:

#### Command Line Option

```bash
./book2ru-go --start-batch 13 < input.txt >> output.txt
```

#### Configuration File Option

```yaml
# .book2ru-go.yml
start_batch: 13
```

#### Step-by-Step Recovery Process

1. **Note the batch number** from the error message (the tool will tell you exactly what to do):
   ```
   [ERROR] translating batch 13: decoding API response: read tcp ...
   
   To resume from this batch, use: --start-batch 13 --batch-size 10000
   ```

2. **Resume from that batch using command line**:
   ```bash
   ./book2ru-go --start-batch 13 --batch-size 10000 < input.txt >> output.txt
   ```

   **Important**: Use `>>` (append) instead of `>` (overwrite) to continue adding to your existing output file.

3. **Or set it in configuration file and run normally**:
   ```bash
   echo "start_batch: 13" >> .book2ru-go.yml
   echo "batch_size: 10000" >> .book2ru-go.yml
   ./book2ru-go < input.txt >> output.txt
   ```

4. **For large files**, you can also save progress and resume later:
   ```bash
   # First run (processes batches 1-12, then fails at 13)
   ./book2ru-go --batch-size 15000 < large_book.txt > partial_translation.txt
   
   # Resume from batch 13 with same batch size
   ./book2ru-go --start-batch 13 --batch-size 15000 < large_book.txt >> partial_translation.txt
   ```

**Note**: The Ruby version does not support resuming from a specific batch. If you need this functionality, use the Go
version.

## How It Works

1. The application reads the entire text content from STDIN.
2. It splits the content into batches of approximately 10KB, ensuring that lines are not broken in the middle.
3. Each batch is sent to the OpenRouter API for translation.
4. The translated text is immediately written to STDOUT, preserving the original structure.
5. All progress, status, and API responses are logged to STDERR, so they don't interfere with the output file.
