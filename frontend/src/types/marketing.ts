/**
 * Типы для модуля Маркетинг
 * Используйте эти типы при интеграции с API
 */

// Экспорт для того, чтобы файл считался модулем
export {};

export interface ClickMetric {
    id: number;
    indicator: string;
    october: number | string;
    september: number | string;
    august: number | string;
    efficiency: number;
}

export interface ConversionMetric {
    id: number;
    indicator: string;
    october: number | string;
    september: number | string;
    august: number | string;
    efficiency: number;
}

export interface SummaryItem {
    label: string;
    value: string;
    change: number;
    isPositive: boolean;
}

export interface MarketingData {
    clicks: {
        summary: SummaryItem[];
        metrics: ClickMetric[];
    };
    conversions: {
        summary: SummaryItem[];
        metrics: ConversionMetric[];
    };
}

/**
 * Пример структуры ответа API:
 * 
 * {
 *   "clicks": {
 *     "summary": [
 *     {
 *       "label": "Клики",
 *       "value": "Упало на 5%",
 *       "change": -5,
 *       "isPositive": false
 *     }
 *   ],
 *   "metrics": [
 *     {
 *       "id": 1,
 *       "indicator": "Клики, кол-во",
 *       "october": 3063,
 *       "september": 3218,
 *       "august": 3085,
 *       "efficiency": -4.82
 *     }
 *   ]
 * },
 * "conversions": { ... }
 * }
 */