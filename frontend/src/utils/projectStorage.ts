/**
 * Утилита для работы с последним выбранным проектом в sessionStorage
 */

const LAST_PROJECT_KEY = 'last_selected_project_id';

export const projectStorage = {
    /**
     * Сохранить ID последнего выбранного проекта
     */
    setLastProject(projectId: number | undefined | null): void {
        try {
            if (projectId === undefined || projectId === null || isNaN(Number(projectId))) {
                console.error('[ProjectStorage] Invalid project ID:', projectId);
                return;
            }
            sessionStorage.setItem(LAST_PROJECT_KEY, projectId.toString());
            console.log('[ProjectStorage] Last project saved:', projectId);
        } catch (error) {
            console.error('[ProjectStorage] Error saving last project:', error);
        }
    },

    /**
     * Получить ID последнего выбранного проекта
     */
    getLastProject(): number | null {
        try {
            const projectId = sessionStorage.getItem(LAST_PROJECT_KEY);
            if (!projectId) {
                return null;
            }
            const parsed = parseInt(projectId, 10);
            if (isNaN(parsed) || parsed <= 0) {
                return null;
            }
            return parsed;
        } catch (error) {
            console.error('[ProjectStorage] Error getting last project:', error);
            return null;
        }
    },

    /**
     * Очистить последний выбранный проект
     */
    clearLastProject(): void {
        try {
            sessionStorage.removeItem(LAST_PROJECT_KEY);
            console.log('[ProjectStorage] Last project cleared');
        } catch (error) {
            console.error('[ProjectStorage] Error clearing last project:', error);
        }
    },
};

