import { api } from './apiClient';
import type { MarketingData } from '../../types/marketing';

export const marketingService = {
    async getMarketing(projectId: number): Promise<MarketingData> {
        try {
            const response = await api.get<MarketingData>(`/projects/${projectId}/marketing`);
            console.log('[MarketingService] Fetched marketing data for project:', projectId, response.data);
            return response.data;
        } catch (error: any) {
            console.error('[MarketingService] Failed to fetch marketing data:', error.response?.data || error.message);
            throw error;
        }
    },
};

