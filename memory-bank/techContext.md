# ТЕХНИЧЕСКИЙ КОНТЕКСТ

## Платформа
- **ОС**: macOS (Darwin 24.5.0)
- **Архитектура**: Не определена
- **Shell**: /bin/zsh
- **Рабочая директория**: /Users/dog/Sites/book2ru

## Основной стек технологий
- **Backend**: Ruby on Rails (версия не определена)
- **Database**: PostgreSQL
- **Frontend**: Hotwire (Turbo + Stimulus)
- **Styling**: Tailwind CSS
- **Version Control**: Git

## Инструменты разработки
- **IDE**: Cursor (.cursor/ настройки обнаружены)
- **IDE**: IntelliJ IDEA (.idea/ настройки обнаружены)
- **Ruby Version Manager**: rbenv (рекомендуется)

## Структура проекта
```
book2ru/
├── .git/           # Git репозиторий
├── .gitignore      # Игнорируемые файлы
├── .cursor/        # Настройки Cursor
├── .idea/          # Настройки IntelliJ
├── _src/           # Исходные файлы
│   ├── Chip_War.epub
│   ├── Chip_War.pdf
│   └── Chip_War.txt
├── memory-bank/    # Система управления памятью
└── README.md       # Документация
```

## Требования к окружению
- **Ruby**: Версия для Rails 8 (Ruby 3.2+)
- **PostgreSQL**: Для хранения данных
- **Node.js**: Для работы с Hotwire/Stimulus
- **Bundler**: Для управления зависимостями

## Конфигурация разработки
- **Gem management**: Через rbenv (избегать глобальных gem)
- **Code style**: RuboCop для проверки стиля
- **Testing**: RSpec с FactoryBot
- **Doubles**: Предпочтение verifying doubles

## Специфичные для проекта требования
- **Обработка текста**: Работа с форматами epub, pdf, txt
- **Русский язык**: Поддержка кириллицы
- **Файловая система**: Работа с различными форматами книг

## Потенциальные зависимости
```ruby
# Предполагаемые gem'ы (требует уточнения)
gem 'rails', '~> 8.0'
gem 'pg'
gem 'turbo-rails'
gem 'stimulus-rails'
gem 'tailwindcss-rails'

# Для работы с текстом
gem 'epub-parser'     # Для epub файлов
gem 'pdf-reader'      # Для PDF файлов
gem 'nokogiri'        # Для HTML/XML

# Для тестирования
gem 'rspec-rails'
gem 'factory_bot_rails'
gem 'rubocop'
```

## Настройки безопасности
- Strong parameters в контроллерах
- Защита от XSS, CSRF
- Валидация входных данных
- Безопасная работа с файлами

## Производительность
- Кэширование (Russian Doll caching)
- Eager loading для избежания N+1 запросов
- Индексация базы данных
- Оптимизация запросов

## Deployment (будущий)
- Настройка production окружения
- Конфигурация веб-сервера
- База данных production
- Мониторинг и логирование 