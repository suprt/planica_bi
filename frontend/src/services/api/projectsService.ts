/**
 * Projects Service
 * 
 * Сервис для работы с проектами.
 * Получает список проектов пользователя, детали проекта и т.д.
 */

import { api } from './apiClient';

/**
 * Типы данных
 */

// Модель проекта
export interface Project {
    id: number;
    name: string;
    slug: string;
    timezone: string;
    currency: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

// Список проектов с пагинацией (если потребуется)
export interface ProjectsResponse {
    projects: Project[];
    total?: number;
}

/**
 * Projects Service
 */
export const projectsService = {
    /**
     * Получить все проекты текущего пользователя
     * 
     * Backend endpoint: GET /api/projects
     * Response: [{ id, name, slug, timezone, currency, is_active, created_at, updated_at }, ...]
     * 
     * @returns Promise со списком проектов
     */
    async getAll(): Promise<Project[]> {
        try {
            const response = await api.get<any>('/projects');
            console.log('[ProjectsService] Fetched projects response:', response.data);
            
            // Backend может возвращать либо массив, либо { data: [...], total: N }
            let projects: Project[];
            if (Array.isArray(response.data)) {
                projects = response.data;
            } else if (response.data && Array.isArray(response.data.data)) {
                projects = response.data.data;
            } else {
                console.warn('[ProjectsService] Unexpected response format:', response.data);
                projects = [];
            }
            
            console.log('[ProjectsService] Fetched projects:', projects.length);
            // Логируем первый проект для отладки
            if (projects.length > 0) {
                console.log('[ProjectsService] First project (raw):', projects[0]);
                console.log('[ProjectsService] First project keys:', Object.keys(projects[0]));
                console.log('[ProjectsService] First project ID:', projects[0].id, 'Type:', typeof projects[0].id);
            }
            return projects;
        } catch (error: any) {
            console.error('[ProjectsService] Failed to fetch projects:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Получить конкретный проект по ID
     * 
     * Backend endpoint: GET /api/projects/:id
     * Response: { id, name, slug, timezone, currency, is_active, created_at, updated_at }
     * 
     * @param projectId - ID проекта
     * @returns Promise с данными проекта
     */
    async getById(projectId: number): Promise<Project> {
        try {
            const response = await api.get<Project>(`/projects/${projectId}`);
            console.log('[ProjectsService] Fetched project:', response.data.name);
            return response.data;
        } catch (error: any) {
            console.error('[ProjectsService] Failed to fetch project:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Создать новый проект
     * 
     * Backend endpoint: POST /api/projects
     * Request: { name, slug, timezone, currency }
     * Response: { id, name, slug, timezone, currency, is_active, created_at, updated_at }
     * 
     * @param projectData - Данные нового проекта
     * @returns Promise с созданным проектом
     */
    async create(projectData: {
        name: string;
        slug: string;
        timezone?: string;
        currency?: string;
    }): Promise<Project> {
        try {
            const response = await api.post<Project>('/projects', projectData);
            console.log('[ProjectsService] Created project:', response.data.name);
            return response.data;
        } catch (error: any) {
            console.error('[ProjectsService] Failed to create project:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Обновить проект
     * 
     * Backend endpoint: PUT /api/projects/:id
     * Request: { name?, slug?, timezone?, currency?, is_active? }
     * Response: { id, name, slug, timezone, currency, is_active, created_at, updated_at }
     * 
     * @param projectId - ID проекта
     * @param updates - Данные для обновления
     * @returns Promise с обновленным проектом
     */
    async update(projectId: number, updates: Partial<{
        name: string;
        slug: string;
        timezone: string;
        currency: string;
        is_active: boolean;
    }>): Promise<Project> {
        try {
            const response = await api.put<Project>(`/projects/${projectId}`, updates);
            console.log('[ProjectsService] Updated project:', response.data.name);
            return response.data;
        } catch (error: any) {
            console.error('[ProjectsService] Failed to update project:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Удалить проект
     * 
     * Backend endpoint: DELETE /api/projects/:id
     * Response: 204 No Content
     * 
     * @param projectId - ID проекта
     */
    async delete(projectId: number): Promise<void> {
        try {
            await api.delete(`/projects/${projectId}`);
            console.log('[ProjectsService] Deleted project:', projectId);
        } catch (error: any) {
            console.error('[ProjectsService] Failed to delete project:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Получить публичную ссылку на отчет проекта
     * 
     * Backend endpoint: GET /api/projects/:id/public-link
     * Response: { public_url: string, public_token: string, project_id: number, project_name: string }
     * 
     * @param projectId - ID проекта
     * @returns Promise с публичной ссылкой
     */
    async getPublicLink(projectId: number): Promise<{ public_url: string; public_token: string; project_id: number; project_name: string }> {
        try {
            const response = await api.get<any>(`/projects/${projectId}/public-link`);
            console.log('[ProjectsService] Fetched public link:', response.data.public_url);
            return response.data;
        } catch (error: any) {
            console.error('[ProjectsService] Failed to get public link:', error.response?.data || error.message);
            throw error;
        }
    },
};

