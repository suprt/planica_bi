/**
 * Reports Service
 * 
 * Сервис для работы с отчетами проектов.
 * Получает аналитические данные по Яндекс.Метрике, Директу, SEO.
 */

import { api } from './apiClient';

/**
 * Типы данных для отчетов
 */

// Динамика по сравнению с предыдущим периодом
export interface Dynamics {
    visits: number;      // % изменения визитов
    users: number;       // % изменения пользователей
    bounce: number;      // % изменения отказов
    avgSec: number;      // % изменения времени на сайте
    conv?: number;       // % изменения конверсии (опционально)
}

// Итоговые метрики по Метрике (за 1 месяц)
export interface MetricsSummary {
    month: string;           // Период (например "2024-10")
    visits: number;          // Визиты
    users: number;           // Пользователи
    bounce: number;          // Отказы (%)
    avgSec: number;         // Среднее время на сайте (секунды)
    conv?: number;           // Конверсия (%, опционально)
    dynamics?: Dynamics;      // Динамика по сравнению с предыдущим периодом
}

// Метрики по возрастным группам
export interface AgeMetrics {
    month: string;           // Период
    age: string;             // Возрастная группа (например "18-24")
    visits: number;
    users: number;
    bounce: number;
    avgSec: number;
}

// Данные кампании Яндекс.Директ
export interface DirectCampaign {
    campaignId: number;      // ID кампании
    name: string;            // Название кампании
    rows: DirectCampaignRow[]; // Данные по месяцам
}

// Строка с данными кампании за месяц
export interface DirectCampaignRow {
    month: string;           // Период
    impressions: number;     // Показы
    clicks: number;          // Клики
    ctr: number;             // CTR (%)
    cpc: number;             // CPC (стоимость за клик)
    conv?: number;           // Конверсия (опционально)
    cpa?: number;            // CPA (стоимость конверсии, опционально)
    cost: number;            // Общие расходы
}

// Итоги Яндекс.Директ (суммарно по всем кампаниям)
export interface DirectTotals {
    month: string;
    impressions: number;
    clicks: number;
    ctr: number;
    cpc: number;
    conv?: number;
    cpa?: number;
    cost: number;
}

// Данные по SEO
export interface SeoSummary {
    month: string;
    visitors: number;        // Посетители из поиска
    conv: number;            // Конверсия (%)
}

// Поисковые запросы
export interface SeoQuery {
    month: string;
    query: string;           // Текст запроса
    position: number;        // Позиция в выдаче
    url?: string;            // URL страницы (опционально)
}

// AI-инсайты (результат анализа Ollama)
export interface AiInsights {
    summary: string;         // Краткая сводка
    recommendations: string[]; // Рекомендации
}

// Полный отчет по проекту
export interface Report {
    projectId: number;       // ID проекта
    periods: string[];       // Периоды отчета (например ["2024-10", "2024-09", "2024-08"])
    metrica: {
        summary: MetricsSummary[];
        age: AgeMetrics[];
    };
    direct: {
        totals: DirectTotals[];
        campaigns: DirectCampaign[];
    };
    seo: {
        summary: SeoSummary[];
        queries: SeoQuery[];
    };
    ai_insights?: AiInsights; // Опционально (если есть AI-анализ)
}

// Статус генерации отчета (когда отчет еще не готов)
export interface ReportStatus {
    status: 'pending' | 'processing' | 'completed' | 'failed';
    message: string;
    note?: string;
    project_id: number;
    task_id?: string;
    queue?: string;
}

// Метрики по каналам (упрощенный формат для таблиц)
export interface ChannelMetrics {
    channel: string;         // Канал (например "Yandex.Metrica", "Yandex.Direct")
    data: {
        [key: string]: any;  // Данные по месяцам
    };
}

/**
 * Reports Service
 */
export const reportsService = {
    /**
     * Получить полный отчет по проекту
     * 
     * Backend endpoint: GET /api/report/:id
     * Response: { projectId, periods, metrica, direct, seo, ai_insights? }
     * 
     * @param projectId - ID проекта
     * @returns Promise с полным отчетом
     */
    async getReport(projectId: number): Promise<Report | ReportStatus> {
        try {
            const response = await api.get<Report | ReportStatus>(`/report/${projectId}`);
            console.log('[ReportsService] Fetched report for project:', projectId);
            console.log('[ReportsService] Response type:', (response.data as any).status ? 'status' : 'report');
            
            // Проверяем, это статус задачи или сам отчет
            if ((response.data as any).status) {
                // Это статус задачи, а не отчет
                return response.data as ReportStatus;
            }
            
            return response.data as Report;
        } catch (error: any) {
            console.error('[ReportsService] Failed to fetch report:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Получить метрики по каналам
     * 
     * Backend endpoint: GET /api/channel-metrics/:id?periods=2024-10,2024-09,2024-08
     * Response: массив данных по каналам
     * 
     * @param projectId - ID проекта
     * @param periods - Массив периодов (например ["2024-10", "2024-09"])
     * @returns Promise с метриками по каналам
     */
    async getChannelMetrics(projectId: number, periods: string[]): Promise<any> {
        try {
            const periodsParam = periods.join(',');
            const response = await api.get(`/channel-metrics/${projectId}`, {
                periods: periodsParam
            });
            console.log('[ReportsService] Fetched channel metrics for project:', projectId);
            return response.data;
        } catch (error: any) {
            console.error('[ReportsService] Failed to fetch channel metrics:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Проверить статус генерации отчета (для асинхронной генерации)
     * 
     * Backend endpoint: GET /api/report-status/:taskId
     * Response: { status: "pending" | "processing" | "completed" | "failed", result?: Report }
     * 
     * @param taskId - ID задачи генерации отчета
     * @returns Promise со статусом
     */
    async getReportStatus(taskId: string): Promise<{
        status: 'pending' | 'processing' | 'completed' | 'failed';
        result?: Report;
        error?: string;
    }> {
        try {
            const response = await api.get(`/report-status/${taskId}`);
            return response.data;
        } catch (error: any) {
            console.error('[ReportsService] Failed to fetch report status:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Экспорт отчета в PDF/Excel (если будет реализовано)
     * 
     * @param projectId - ID проекта
     * @param format - Формат экспорта ("pdf" | "excel")
     */
    async exportReport(projectId: number, format: 'pdf' | 'excel'): Promise<Blob> {
        try {
            const response = await api.get(`/report/${projectId}/export`, {
                format
            });
            console.log('[ReportsService] Exported report:', format);
            return response.data;
        } catch (error: any) {
            console.error('[ReportsService] Failed to export report:', error.response?.data || error.message);
            throw error;
        }
    },
};

