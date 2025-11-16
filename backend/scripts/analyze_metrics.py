#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Анализ метрик рекламных каналов с использованием Ollama
"""

import json
import sys
import os
import requests

def analyze_channel(name, ch):
    """Анализирует метрики одного канала и возвращает список инсайтов"""
    out = []
    
    cpc = ch["cpc"]
    ctr = ch["ctr"]
    cpa = ch["cpa"]
    conv = ch["conversions"]
    
    # Проверяем, что есть минимум 2 периода для сравнения
    if len(cpc) < 2 or len(ctr) < 2 or len(cpa) < 2 or len(conv) < 2:
        return out
    
    # Динамика (сравнение текущего периода с предыдущим)
    # Периоды идут от нового к старому: [текущий, предыдущий, старый]
    # Сравниваем первый (текущий) с вторым (предыдущим)
    dcpc = (cpc[0] - cpc[1]) / cpc[1] * 100 if cpc[1] != 0 else 0
    dctr = (ctr[0] - ctr[1]) / ctr[1] * 100 if ctr[1] != 0 else 0
    dcpa = (cpa[0] - cpa[1]) / cpa[1] * 100 if cpa[1] != 0 else 0
    dconv = (conv[0] - conv[1]) / conv[1] * 100 if conv[1] != 0 else 0
    
    if dctr > 5:
        out.append(f"CTR {name} вырос на {dctr:.1f}% — объявления стали привлекательнее.")
    elif dctr < -5:
        out.append(f"CTR {name} снизился на {abs(dctr):.1f}% — стоит обновить креативы.")
    
    if dcpa > 5:
        out.append(f"CPA {name} вырос на {dcpa:.1f}% — реклама дорожает.")
    elif dcpa < -5:
        out.append(f"CPA {name} снизился на {abs(dcpa):.1f}% — улучшилась эффективность.")
    
    if dconv > 5:
        out.append(f"Конверсии {name} выросли на {dconv:.1f}%.")
    elif dconv < -5:
        out.append(f"Конверсии {name} снизились на {abs(dconv):.1f}% — требуется оптимизация.")
    
    return out


def main():
    """Основная функция"""
    # Читаем JSON из stdin
    try:
        data = json.load(sys.stdin)
    except json.JSONDecodeError as e:
        print(json.dumps({"error": f"Invalid JSON: {str(e)}"}), file=sys.stderr)
        sys.exit(1)
    
    # Собираем аналитические факты
    insights = []
    for name, ch in data.get("metrics", {}).items():
        insights.extend(analyze_channel(name, ch))
    
    base_text = "\n".join(insights) if insights else "Изменения метрик незначительны."
    
    # Получаем API ключ Ollama из переменной окружения
    ollama_api_key = os.getenv("OLLAMA_API_KEY")
    ollama_api_url = os.getenv("OLLAMA_API_URL", "https://ollama.com/api")
    ollama_model = os.getenv("OLLAMA_MODEL", "glm-4.6")
    
    result = {
        "analytical_facts": base_text,
        "ai_report": None,
        "error": None
    }
    
    # Если есть API ключ, вызываем Ollama
    if ollama_api_key:
        try:
            prompt = f"""Ты опытный маркетинг‑аналитик. На основе предоставленных аналитических фактов сделай краткие выводы и рекомендации.

Аналитические факты:

{base_text}

Требования:
- Напиши максимум один абзац кратких выводов по результатам
- НЕ перечисляй конкретные цифры и проценты (пользователь их уже видит)
- Сделай выводы о трендах, проблемах и возможностях
- Дай 3-5 конкретных рекомендаций в виде списка
- Будь лаконичным и по делу

Формат: Один абзац выводов, затем список из 3-5 рекомендаций."""
            
            # Используем Ollama API
            headers = {
                "Authorization": f"Bearer {ollama_api_key}",
                "Content-Type": "application/json"
            }
            
            payload = {
                "model": ollama_model,
                "prompt": prompt,
                "stream": False,
                "options": {
                    "temperature": 0.4
                }
            }
            
            # Пробуем разные форматы запросов и URL
            # Формат: (URL, payload, headers_override, response_extractor)
            endpoints_to_try = [
                # Локальный Ollama стандартный формат
                (f"{ollama_api_url}/generate", payload, headers, lambda r: r.json()["response"]),
                # Локальный Ollama chat формат
                (f"{ollama_api_url}/chat", {
                    "model": ollama_model,
                    "messages": [
                        {"role": "system", "content": "Ты анализируешь рекламные метрики и пишешь краткие выводы для клиента."},
                        {"role": "user", "content": prompt}
                    ],
                    "stream": False,
                    "options": {"temperature": 0.4}
                }, headers, lambda r: r.json()["message"]["content"]),
                # OpenAI-совместимый формат (для облачных сервисов)
                (f"{ollama_api_url}/v1/chat/completions", {
                    "model": ollama_model,
                    "messages": [
                        {"role": "system", "content": "Ты анализируешь рекламные метрики и пишешь краткие выводы для клиента."},
                        {"role": "user", "content": prompt}
                    ],
                    "temperature": 0.4
                }, headers, lambda r: r.json()["choices"][0]["message"]["content"]),
                # Попытка без /api префикса
                (f"{ollama_api_url.rstrip('/api')}/v1/chat/completions", {
                    "model": ollama_model,
                    "messages": [
                        {"role": "system", "content": "Ты анализируешь рекламные метрики и пишешь краткие выводы для клиента."},
                        {"role": "user", "content": prompt}
                    ],
                    "temperature": 0.4
                }, headers, lambda r: r.json()["choices"][0]["message"]["content"]),
            ]
            
            last_error = None
            success = False
            
            for endpoint_url, endpoint_payload, endpoint_headers, extract_response in endpoints_to_try:
                try:
                    response = requests.post(
                        endpoint_url,
                        headers=endpoint_headers,
                        json=endpoint_payload,
                        timeout=60
                    )
                    response.raise_for_status()
                    result["ai_report"] = extract_response(response)
                    success = True
                    break
                except Exception as e:
                    last_error = str(e)
                    continue
            
            if not success:
                result["error"] = f"Ollama API error: Tried multiple endpoints. Last error: {last_error}"
                
        except Exception as e:
            result["error"] = f"Ollama API error: {str(e)}"
    else:
        result["error"] = "OLLAMA_API_KEY not set"
    
    # Выводим результат в JSON формате
    print(json.dumps(result, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()


