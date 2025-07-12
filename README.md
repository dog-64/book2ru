# Book2Ru 📚

`book2ru` is a powerful and efficient Ruby CLI tool for translating large text files from English to Russian using the OpenRouter API. Designed for simplicity and performance, it processes text through a standard STDIN/STDOUT pipeline, making it a versatile component in any text processing workflow.

## Key Features ✨

- **Large File Processing**: Handles large text files (e.g., books) with ease using a stream-based approach.
- **Intelligent Batching**: Automatically splits input text into 10KB chunks to work efficiently with API limits.
- **Robust Error Handling**: Features an automatic retry mechanism with exponential backoff for network and API errors.
- **Flexible Configuration**: Configure API keys, models, and prompts via a `.env` file, a `.book2ru.yml` file, or command-line arguments.
- **Metadata Output**: Outputs progress and debugging information to STDERR, keeping STDOUT clean for the translated text.

## Prerequisites

- **Ruby** (version `3.4.2` is recommended)
- **Bundler** for managing gems

## Installation 🚀

1.  **Clone the repository:**
    ```bash
    git clone <repository_url>
    cd book2ru
    ```

2.  **Install dependencies:**
    ```bash
    bundle install
    ```

3.  **Make the script executable:**
    ```bash
    chmod +x book2ru.rb
    ```

## Configuration ⚙️

1.  **API Key**: Create a `.env` file in the project root and add your OpenRouter API key:

    ```env
    # .env
    OPENROUTER_KEY='sk-or-v1-...'
    ```

2.  **Custom Settings (Optional)**: Create a `.book2ru.yml` file to override default settings like the model or translation prompt.

    ```yaml
    # .book2ru.yml
    model: "anthropic/claude-3-5-sonnet-20240620"
    prompt: "Translate the following English text to Russian. Return only the translated text, preserving all original formatting and line breaks:"
    ```

## Usage 📖

The tool reads from STDIN and writes the translated output to STDOUT.

### Basic Translation

```bash
./book2ru.rb < my_book.txt > my_book_ru.txt
```

### Using a Different Model

You can specify a different model via the command line:

```bash
./book2ru.rb --model "mistralai/mistral-large-latest" < input.txt > output.txt
```

### Command-Line Options

| Option               | Alias | Description                               |
| -------------------- | ----- | ----------------------------------------- |
| `--help`             | `-h`  | Show the help message.                    |
| `--version`          |       | Display version and model information.    |
| `--model <MODEL>`      | `-m`  | Specify the LLM model to use.             |
| `--openrouter_key <KEY>` | `-o`  | Provide the OpenRouter API key directly.  |

## How It Works

1.  The script reads the entire text content from STDIN.
2.  It splits the content into batches of approximately 10KB, ensuring that lines are not broken in the middle.
3.  Each batch is sent to the OpenRouter API for translation.
4.  The translated text is immediately written to STDOUT, preserving the original structure.
5.  All progress, status, and API responses are logged to STDERR, so they don't interfere with the output file.
