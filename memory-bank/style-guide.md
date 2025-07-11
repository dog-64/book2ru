# РУКОВОДСТВО ПО СТИЛЮ

## Стиль кода Ruby

### Общие принципы
- Следовать Ruby Style Guide (https://rubystyle.guide/)
- Использовать выразительный синтаксис Ruby
- Предпочитать итерацию и модуляризацию дублированию кода
- Короткие функции лучше длинных

### Конвенции именования
- `snake_case` для файлов, методов, переменных
- `CamelCase` для классов и модулей
- Описательные имена: `user_signed_in?`, `calculate_total`
- Файлы классов: `ExampleClass` → `example_class.rb`
- Вложенные классы: `Parent::TestClass` → `parent/test_class.rb`

### Синтаксис
```ruby
# Предпочитать одинарные кавычки
message = 'Hello world'
interpolated = "Hello #{name}"

# Использовать ||= для присваивания по умолчанию
@user ||= User.current

# Использовать unless для отрицания
redirect_to login_path unless user_signed_in?

# Предпочитать &. для safe navigation
user&.name&.upcase
```

### Структура файлов
```ruby
# Документирование функций
# Calculates the total price including tax
def calculate_total_with_tax(price, tax_rate)
  price * (1 + tax_rate)
end

# Использование приватных методов
private

def helper_method
  # Вспомогательная логика
end
```

## Rails конвенции

### Модели
```ruby
class Book < ApplicationRecord
  validates :title, presence: true
  validates :isbn, uniqueness: true
  
  scope :published, -> { where(published: true) }
  
  def display_title
    title.titleize
  end
  
  private
  
  def ensure_isbn_format
    # Валидация ISBN
  end
end
```

### Контроллеры
```ruby
class BooksController < ApplicationController
  before_action :authenticate_user!
  before_action :set_book, only: [:show, :edit, :update, :destroy]
  
  def index
    @books = Book.published.includes(:author)
  end
  
  def show
    # @book уже установлен в before_action
  end
  
  private
  
  def set_book
    @book = Book.find(params[:id])
  end
  
  def book_params
    params.require(:book).permit(:title, :content, :published)
  end
end
```

### Сервисы
```ruby
# app/services/book_parser_service.rb
class BookParserService
  def initialize(file_path)
    @file_path = file_path
  end
  
  def call
    parse_content
  end
  
  private
  
  attr_reader :file_path
  
  def parse_content
    # Логика парсинга
  end
end
```

## Тестирование

### RSpec стиль
```ruby
# spec/models/book_spec.rb
RSpec.describe Book, type: :model do
  describe 'validations' do
    it 'validates presence of title' do
      book = build(:book, title: nil)
      expect(book).not_to be_valid
    end
  end
  
  describe '#display_title' do
    it 'returns titleized title' do
      book = build(:book, title: 'test book')
      expect(book.display_title).to eq('Test Book')
    end
  end
end
```

### Фабрики
```ruby
# spec/factories/books.rb
FactoryBot.define do
  factory :book do
    title { 'Sample Book' }
    content { 'Sample content' }
    published { true }
    
    trait :unpublished do
      published { false }
    end
  end
end
```

### Doubles
```ruby
# Предпочитать verifying doubles
let(:book_service) { instance_double(BookParserService) }

before do
  allow(BookParserService).to receive(:new).and_return(book_service)
  allow(book_service).to receive(:call).and_return(parsed_content)
end
```

## Frontend (Hotwire/Stimulus)

### Stimulus контроллеры
```javascript
// app/javascript/controllers/search_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = ["input", "results"]
  
  connect() {
    console.log("Search controller connected")
  }
  
  search() {
    // Логика поиска
  }
}
```

### Tailwind CSS
```erb
<!-- Использовать utility-first подход -->
<div class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl">
  <div class="md:flex">
    <div class="p-8">
      <div class="uppercase tracking-wide text-sm text-indigo-500 font-semibold">
        Book Title
      </div>
    </div>
  </div>
</div>
```

## Качество кода

### Проверки перед коммитом
```bash
# Запуск тестов
bundle exec rspec

# Проверка стиля кода
bundle exec rubocop

# Исправление автоматически исправляемых нарушений
bundle exec rubocop -a
```

### Принципы
1. **DRY** - Don't Repeat Yourself
2. **SOLID** - Принципы объектно-ориентированного программирования
3. **RESTful** - Следовать REST конвенциям
4. **Security First** - Безопасность приоритетна
5. **Performance Conscious** - Учитывать производительность

## Документация

### Комментарии
```ruby
# Документировать сложную логику
# Parses different book formats and extracts content
# Returns Hash with structured content
def parse_book_content(file_path)
  # Реализация
end
```

### README структура
- Описание проекта
- Инструкции по установке
- Примеры использования
- Конфигурация
- Тестирование
- Deployment

## Безопасность

### Валидация
```ruby
# Всегда валидировать пользовательский ввод
validates :email, presence: true, format: { with: URI::MailTo::EMAIL_REGEXP }

# Использовать strong parameters
def book_params
  params.require(:book).permit(:title, :content, :published)
end
```

### Файлы
```ruby
# Безопасная работа с файлами
def safe_file_upload(file)
  return unless file.present?
  return unless allowed_file_types.include?(file.content_type)
  
  # Дополнительная валидация
end
``` 