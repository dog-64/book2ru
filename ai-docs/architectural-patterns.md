# Architectural Patterns

This document describes the key architectural patterns used in the `book2ru` project. The architecture is intentionally simple and optimized for its purpose as a command-line text processing utility.

## 1. Core Pattern: Single-Class CLI Application

The entire application logic is encapsulated within a single class, `Book2ru`, in the `book2ru.rb` file.

-   **Description**: This pattern avoids the overhead of a complex file structure (like MVC or multi-module designs) and packages the entire functionality into one cohesive unit. The class is instantiated and its `run` method is invoked when the script is executed.
-   **When to Use**: This pattern is ideal for small to medium-sized CLI tools where the logic is focused on a single, primary task (in this case, translation). It simplifies deployment and maintenance.

## 2. Data Handling Patterns

### Stream-Based I/O (Unix Pipeline)

-   **Description**: The application uses `STDIN` to read input and `STDOUT` to write output. This is a classic Unix philosophy pattern that makes the tool highly versatile. It allows `book2ru` to be chained with other command-line utilities. For example, you could `cat` a file into it, or pipe the output to `grep` or `less`.
-   **When to Use**: This pattern is fundamental for CLI tools that process text or data, as it decouples the application from the specifics of file I/O and allows for flexible composition.

### Batch Processing

-   **Description**: The core business logic revolves around processing a large text in smaller, manageable chunks. The `create_batches_from_content` method implements this pattern by splitting the input text into batches of approximately 10KB without breaking lines.
-   **When to Use**: This pattern is critical when interacting with external APIs that have request size limits or to manage memory consumption when processing very large files. It ensures that the application can handle inputs of any size without crashing or violating API constraints.

## 3. External Service Integration

### API Client / Gateway

-   **Description**: All interaction with the external OpenRouter service is isolated within the `translate_text_batch` method. This method acts as a "Gateway" to the API, encapsulating the details of building the HTTP request, setting headers, and parsing the response.
-   **When to Use**: This pattern should be used whenever the application needs to communicate with an external service. It decouples the core application logic from the specifics of the API, making it easier to maintain, test, or even replace the service in the future.

### Retry with Exponential Backoff

-   **Description**: Within the API Client pattern, a robust retry mechanism is implemented to handle transient network or API errors. If a request fails, the application waits for a progressively longer period (`sleep(2 ** attempts)`) before retrying, up to a configured maximum number of attempts.
-   **When to Use**: This is an essential pattern for any application that relies on network requests to external services. It improves the reliability and resilience of the application by automatically handling temporary issues without failing the entire process.

## 4. Configuration Management

### Hierarchical Configuration

-   **Description**: The application uses a layered approach to configuration, allowing for flexible setup with clear precedence.
    1.  **Hardcoded Defaults**: The lowest priority, defined in the `initialize` method.
    2.  **YAML File**: A `.book2ru.yml` file can override the defaults.
    3.  **Environment Variables**: The `.env` file (for secrets like API keys) overrides the YAML file.
    4.  **Command-Line Arguments**: The highest priority, allowing for on-the-fly adjustments.
-   **When to Use**: This pattern provides a powerful and user-friendly way to configure an application. It allows users to set persistent global settings in a file while still providing the ability to override them for specific runs via the command line. 