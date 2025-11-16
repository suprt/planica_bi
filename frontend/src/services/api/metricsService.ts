import { api } from './apiClient';

/**
 * Типы данных для метрик Яндекс.Метрики
 */

// Метрики за месяц
export interface MetricsRow {
    month: string;           // Период (например "2024-10")
    visits: number;          // Визиты
    users: number;           // Пользователи
    bounce_rate: number;     // Отказы (%)
    avg_sec: number;         // Среднее время на сайте (секунды)
    conversions?: number;    // Конверсия (опционально)
}

// Метрики проекта
export interface MetricsWithData {
    projectId: number;       // ID проекта
    rows: MetricsRow[];      // Данные по месяцам
}

export const metricsService = {
    /**
     * Получить метрики проекта
     */
    async getMetrics(projectId: number): Promise<MetricsWithData> {
        try {
            const response = await api.get<MetricsWithData>(`/projects/${projectId}/metrics`);
            console.log('[MetricsService] Fetched metrics:', response.data);
            return response.data;
        } catch (error: any) {
            console.error('[MetricsService] Failed to fetch metrics:', error.response?.data || error.message);
            throw error;
        }
    },
};

