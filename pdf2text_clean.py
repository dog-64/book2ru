#!/usr/bin/env python3
"""
Приведение PDF к тексту

Примеры использования:
  cat "liars poker - michael lewis.pdf" | python3 pdf2text_clean.py > output.txt
или
  python3 pdf2text_clean.py < "liars poker - michael lewis.pdf" > "liars poker - michael lewis.txt"

Требования:
  pip install PyMuPDF
"""

import sys
import fitz  # PyMuPDF
import re
import io

def normalize_linebreaks(text: str) -> str:
    """Убираем одиночные переносы строк, оставляем двойные (абзацы)"""
    text = re.sub(r'(?<!\n)\n(?!\n)', ' ', text)
    return re.sub(r' +', ' ', text).strip()

def main():
    try:
        # Чтение PDF из stdin (байтовый поток)
        pdf_bytes = sys.stdin.buffer.read()
        
        if not pdf_bytes:
            print("Ошибка: Нет данных во входном потоке", file=sys.stderr)
            sys.exit(1)
        
        # Открываем PDF из памяти
        doc = fitz.open("pdf", stream=pdf_bytes)
        
        # Извлекаем и объединяем текст со всех страниц
        raw_text = "\n\n".join(page.get_text("text") for page in doc)
        
        if not raw_text.strip():
            print("Предупреждение: PDF не содержит извлекаемого текста", file=sys.stderr)
        
        cleaned_text = normalize_linebreaks(raw_text)
        
        # Вывод в stdout
        sys.stdout.write(cleaned_text)
        
        doc.close()
        
    except Exception as e:
        print(f"Ошибка при обработке PDF: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main() 