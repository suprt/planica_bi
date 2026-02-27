/**
 * API Client Configuration
 * 
 * Централизованный клиент для всех HTTP запросов к backend API.
 * Использует axios с настроенными interceptors для:
 * - Автоматического добавления JWT токена в headers
 * - Обработки ошибок авторизации (401)
 * - Логирования запросов в development режиме
 */

import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

// Базовый URL для API - берем из environment variables или fallback на localhost
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

// Создаем экземпляр axios с базовой конфигурацией
const apiClient = axios.create({
    baseURL: API_BASE_URL,
    timeout: 30000, // 30 секунд на запрос
    headers: {
        'Content-Type': 'application/json',
    },
});

/**
 * Request Interceptor
 * Вызывается перед каждым запросом
 * Добавляет JWT токен из sessionStorage в Authorization header
 */
apiClient.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // Получаем токен из sessionStorage
        const token = sessionStorage.getItem('auth_token');
        
        // Если токен есть, добавляем его в headers
        if (token && config.headers) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        // В development режиме логируем запросы
        if (process.env.NODE_ENV === 'development') {
            console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, {
                params: config.params,
                data: config.data,
            });
        }

        return config;
    },
    (error: AxiosError) => {
        // Ошибка при подготовке запроса
        console.error('[API Request Error]', error);
        return Promise.reject(error);
    }
);

/**
 * Response Interceptor
 * Вызывается после получения ответа
 * Обрабатывает общие ошибки (401, 403, 500)
 */
apiClient.interceptors.response.use(
    (response: AxiosResponse) => {
        // Успешный ответ - логируем в development
        if (process.env.NODE_ENV === 'development') {
            console.log(`[API Response] ${response.config.method?.toUpperCase()} ${response.config.url}`, {
                status: response.status,
                data: response.data,
            });
        }
        return response;
    },
    (error: AxiosError) => {
        // Обработка ошибок
        if (error.response) {
            const status = error.response.status;
            const url = error.config?.url;

            // 401 Unauthorized - токен истек или невалиден
            if (status === 401) {
                console.warn('[API] Unauthorized - redirecting to login');

                // Удаляем невалидный токен
                sessionStorage.removeItem('auth_token');
                sessionStorage.removeItem('auth_user');
                
                // Редиректим на страницу логина (только если не на странице логина)
                if (window.location.pathname !== '/login') {
                    window.location.href = '/login';
                }
            }

            // 403 Forbidden - нет доступа к ресурсу
            if (status === 403) {
                console.error('[API] Forbidden - insufficient permissions');
            }

            // 500 Internal Server Error
            if (status === 500) {
                console.error('[API] Server Error');
            }

            // Логируем ошибку в development
            if (process.env.NODE_ENV === 'development') {
                console.error(`[API Error] ${status} ${url}`, {
                    message: error.message,
                    data: error.response.data,
                });
            }
        } else if (error.request) {
            // Запрос был отправлен, но ответа не получено
            console.error('[API] No response received', error.request);
        } else {
            // Ошибка при настройке запроса
            console.error('[API] Request setup error', error.message);
        }

        return Promise.reject(error);
    }
);

/**
 * Helper функции для типизированных запросов
 */

export const api = {
    /**
     * GET запрос
     */
    get: <T = any>(url: string, params?: any): Promise<AxiosResponse<T>> => {
        return apiClient.get<T>(url, { params });
    },

    /**
     * POST запрос
     */
    post: <T = any>(url: string, data?: any): Promise<AxiosResponse<T>> => {
        return apiClient.post<T>(url, data);
    },

    /**
     * PUT запрос
     */
    put: <T = any>(url: string, data?: any): Promise<AxiosResponse<T>> => {
        return apiClient.put<T>(url, data);
    },

    /**
     * DELETE запрос
     */
    delete: <T = any>(url: string): Promise<AxiosResponse<T>> => {
        return apiClient.delete<T>(url);
    },

    /**
     * PATCH запрос
     */
    patch: <T = any>(url: string, data?: any): Promise<AxiosResponse<T>> => {
        return apiClient.patch<T>(url, data);
    },
};

// Экспортируем сам клиент для прямого использования (если нужно)
export default apiClient;

