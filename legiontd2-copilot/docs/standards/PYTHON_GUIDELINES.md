# PYTHON GUIDELINES

## Project Conventions

- Python 3.12+
- Используется только в офлайн ML-контуре (обучение моделей)
- Не является частью runtime-приложения
- Датасет: `cmd/dataset/` на Go, обучение: Python-скрипты

## Code Style

- PEP 8
- Type hints for all function signatures
- `ruff` for linting, `black` for formatting
- Max line length: 100
