import { api } from './apiClient';

/**
 * Типы данных для кампаний
 */

// Метрики кампании за месяц
export interface CampaignMetricsRow {
    month: string;           // Период (например "2024-10")
    impressions: number;     // Показы
    clicks: number;         // Клики
    ctr: number;            // CTR (%)
    cpc: number;           // CPC (стоимость за клик)
    conv?: number;          // Конверсия (опционально)
    cpa?: number;           // CPA (стоимость конверсии, опционально)
    cost: number;          // Общие расходы
}

// Кампания с метриками
export interface CampaignWithMetrics {
    campaignId: number;     // ID кампании
    name: string;           // Название кампании
    status: string;         // Статус кампании
    rows: CampaignMetricsRow[]; // Данные по месяцам
}

export const campaignsService = {
    /**
     * Получить все кампании проекта с метриками
     */
    async getCampaigns(projectId: number): Promise<CampaignWithMetrics[]> {
        try {
            const response = await api.get<CampaignWithMetrics[]>(`/projects/${projectId}/campaigns`);
            console.log('[CampaignsService] Fetched campaigns:', response.data.length);
            return response.data;
        } catch (error: any) {
            console.error('[CampaignsService] Failed to fetch campaigns:', error.response?.data || error.message);
            throw error;
        }
    },
};

