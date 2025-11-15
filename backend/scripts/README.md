# Скрипт анализа метрик

## Описание

Скрипт `analyze_metrics.py` анализирует метрики рекламных каналов и генерирует инсайты с помощью Ollama.

## Установка зависимостей

```bash
pip install -r requirements.txt
```

## Использование

Скрипт читает JSON из stdin и выводит результат в stdout:

```bash
cat metrics.json | python3 analyze_metrics.py
```

## Переменные окружения

- `OLLAMA_API_KEY` - API ключ Ollama (опционально, если не установлен, будет возвращен только базовый анализ без AI)
- `OLLAMA_API_URL` - URL API Ollama (по умолчанию: https://ollama.com/api для облачного, http://localhost:11434/api для локального)
- `OLLAMA_MODEL` - Название модели Ollama (по умолчанию: glm-4.6)

## Формат входных данных

```json
{
  "project": "Project A",
  "periods": ["2025-08", "2025-09", "2025-10"],
  "metrics": {
    "simple": {
      "cpc": [50.0, 55.0, 60.0],
      "impressions": [300000, 350000, 400000],
      "clicks": [2500, 3000, 3500],
      "ctr": [0.8, 0.85, 0.87],
      "conversions": [30, 35, 40],
      "cpa": [2500.0, 2800.0, 3000.0],
      "cost": [125000.0, 165000.0, 210000.0]
    },
    "МК": {...},
    "РСЯ": {...}
  }
}
```

## Формат выходных данных

```json
{
  "analytical_facts": "CTR simple вырос на 2.4% — объявления стали привлекательнее...",
  "ai_report": "На основе анализа метрик за три месяца...",
  "error": null
}
```


