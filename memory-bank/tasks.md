# Task: CLI-приложение book2ru для перевода текстов

## Description
Создание простого CLI-приложения на Ruby для перевода английских текстов на русский язык через OpenRouter API. Приложение читает текст из stdin, разбивает на батчи ~10KB, отправляет в LLM для перевода и выводит результат в stdout.

## Complexity
Level: 3
Type: Intermediate Feature

## Technology Stack
- Language: Ruby 3.4.2
- Framework: Простое CLI-приложение (без Rails)
- HTTP Client: Net::HTTP (встроенный)
- Config: YAML (.book2ru.yml)
- Environment: dotenv gem для .env файлов
- Testing: RSpec
- CLI Parsing: OptionParser (встроенный)

## Technology Validation Checkpoints
- [x] Ruby 3.4.2 установлен и настроен
- [x] Проект инициализирован с Bundler
- [x] Gemfile создан с минимальными зависимостями
- [x] Hello world CLI-скрипт создан и запущен
- [x] Тестовый HTTP-запрос к OpenRouter API выполнен успешно
- [x] RSpec настроен для тестирования

## Status
- [x] Initialization complete
- [x] Planning complete
- [x] Technology validation complete
- [ ] Core batch processing implementation
- [ ] OpenRouter API integration
- [ ] Configuration management
- [ ] Error handling and retry logic
- [ ] CLI argument parsing
- [ ] RSpec tests implementation
- [ ] Documentation and examples

## Implementation Plan

### Phase 1: Project Setup & Hello World ✅
1. Настройка проекта Ruby
   - ✅ Создание Gemfile с минимальными зависимостями
   - ✅ Настройка rbenv для Ruby 3.4.2
   - ✅ Создание базовой структуры проекта

2. Hello World CLI
   - ✅ Создание основного файла book2ru.rb
   - ✅ Базовый CLI-интерфейс с --help и --version
   - ✅ Тестовый запуск для проверки работы

### Phase 2: Core Batch Processing
1. Реализация чтения из stdin
   - Чтение текста по строкам
   - Определение пустых строк (только whitespace)
   - Группировка непустых строк в батчи ~10KB

2. Базовая обработка батчей
   - Алгоритм разбиения без разрыва строк
   - Сохранение порядка обработки
   - Обработка пустых строк (прямое копирование)

### Phase 3: OpenRouter API Integration
1. HTTP-клиент для OpenRouter
   - Настройка базового HTTP-клиента
   - Формирование запросов к API
   - Обработка ответов

2. Обработка ошибок и retry логика
   - Timeout handling
   - Rate limiting (sleep между запросами)
   - Exponential backoff для повторов
   - Логирование ошибок

### Phase 4: Configuration Management
1. YAML конфигурация
   - Загрузка .book2ru.yml
   - Валидация обязательных параметров
   - Значения по умолчанию

2. Environment variables
   - Загрузка OPENROUTER_KEY из env
   - Поддержка .env файлов
   - CLI override для настроек

### Phase 5: CLI Interface
1. Argument parsing
   - --help, --version
   - --model, --openrouter_key
   - --rate-limits для RPM
   - Validation входных параметров

2. Output formatting
   - Metadata footer (опционально)
   - Сохранение форматирования пустых строк
   - Proper UTF-8 encoding

### Phase 6: Testing & Polish
1. RSpec test suite
   - Unit tests для batch processing
   - Mocked API tests
   - CLI argument testing
   - Error handling tests

2. Documentation
   - README с примерами использования
   - Комментарии в коде
   - Примеры конфигурации

## Creative Phases Required
- [ ] N/A - Логика приложения четко определена

## Dependencies
- Ruby 3.4.2 ✅
- OpenRouter API access ✅
- Bundler для управления зависимостями ✅
- RSpec для тестирования ✅
- dotenv для переменных окружения ✅

## Challenges & Mitigations
- **API Rate Limiting**: Реализовать sleep между запросами и exponential backoff
- **Batch Size Management**: Точный подсчет байтов без разрыва строк
- **Error Recovery**: Robust retry logic с максимальным количеством попыток
- **UTF-8 Encoding**: Правильная обработка кириллицы в stdin/stdout
- **Memory Management**: Потоковая обработка для больших файлов

## File Structure
```
book2ru/
├── book2ru.rb               # Основной исполняемый файл ✅
├── Gemfile                  # Зависимости ✅
├── .book2ru.yml.example     # Пример конфигурации ✅
├── spec/
│   ├── book2ru_spec.rb      # Основные тесты ✅
│   └── spec_helper.rb       # Настройки RSpec ✅
├── memory-bank/             # Система управления памятью ✅
└── README.md                # Документация (планируется)
```

## API Integration Details
- Endpoint: https://openrouter.ai/api/v1/chat/completions ✅
- Model: google/gemini-2.0-flash-001 (по умолчанию) ✅
- Headers: Authorization: Bearer <OPENROUTER_KEY> ✅
- Request format: OpenAI compatible JSON ✅
- Response parsing: JSON с choices[0].message.content ✅

## CLI Examples
```bash
# Основное использование
book2ru < input.txt > output-ru.txt

# С параметрами
book2ru --model claude-3-haiku --rate-limits 10 < input.txt > output.txt

# Версия и помощь
book2ru --version
book2ru --help
``` 