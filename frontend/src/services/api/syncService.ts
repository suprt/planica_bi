import { api } from './apiClient';

export interface SyncResponse {
    message: string;
    project_id: number;
    task_id: string;
    queue: string;
}

export const syncService = {
    async syncProject(projectId: number): Promise<SyncResponse> {
        try {
            const response = await api.post<SyncResponse>(`/sync/${projectId}`);
            console.log('[SyncService] Sync task enqueued for project:', projectId, response.data);
            return response.data;
        } catch (error: any) {
            console.error('[SyncService] Failed to sync project:', error.response?.data || error.message);
            throw error;
        }
    },
};
