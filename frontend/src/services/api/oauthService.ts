import { api } from './apiClient';

export interface OAuthStatus {
    authorized: boolean;
    has_token: boolean;
}

export const oauthService = {
    async getStatus(): Promise<OAuthStatus> {
        try {
            const response = await api.get<OAuthStatus>('/oauth/status');
            console.log('[OAuthService] OAuth status:', response.data);
            return response.data;
        } catch (error: any) {
            console.error('[OAuthService] Failed to get OAuth status:', error.response?.data || error.message);
            throw error;
        }
    },
    
    initiateYandexAuth(): void {
        // Redirect to backend OAuth endpoint
        const backendURL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
        window.location.href = `${backendURL}/oauth/yandex`;
    },
};

