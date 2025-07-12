# Naming Conventions

This document outlines the naming conventions used in the `book2ru` project. Adhering to these conventions ensures code consistency, readability, and maintainability. The project follows standard Ruby style guidelines.

## 1. Files

File names should be in `snake_case`.

- **Example**: `book2ru.rb`
- **Example**: `spec_helper.rb`

## 2. Directories

Directory names should be in `snake_case`.

- **Example**: `ai-docs`
- **Example**: `_src`

## 3. Classes and Modules

Class and module names must be in `CamelCase`.

- **Example (Class)**: `Book2ru`
- **Example (Module)**: `OptionParser`

## 4. Methods

Method names must be in `snake_case`.

- **Example**: `parse_options`
- **Example**: `translate_from_stdin`

## 5. Variables

### Local Variables

Local variables must be in `snake_case`.

- **Example**: `exit_code`
- **Example**: `current_batch`

### Instance Variables

Instance variables must be in `snake_case` and prefixed with an `@` symbol.

- **Example**: `@config`
- **Example**: `@option_parser`

## 6. Constants

Constants must be in `UPPER_SNAKE_CASE`. This applies to values that are expected to remain constant throughout the execution of the program.

- **Example**: `VERSION`
- **Example**: `BATCH_SIZE`

## 7. Configuration Keys

Configuration keys in YAML files (`.book2ru.yml`) and Hashes should be in `snake_case`. When loaded, they are converted to symbols.

- **Example (YAML)**: `retry_attempts`
- **Example (Symbol)**: `:retry_attempts` 