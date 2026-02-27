/**
 * Authentication Service
 * 
 * Сервис для работы с авторизацией пользователей.
 * Управляет JWT токенами, login/register, проверкой авторизации.
 */

import { api } from './apiClient';

/**
 * Типы данных
 */

// Запрос на вход
export interface LoginRequest {
    email: string;
    password: string;
}

// Запрос на регистрацию
export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
}

// Роль пользователя в проекте
export interface UserProjectRole {
    id: number;
    user_id: number;
    project_id: number;
    role: 'admin' | 'manager' | 'client';
    created_at: string;
    updated_at: string;
}

// Данные пользователя
export interface User {
    id: number;
    name: string;
    email: string;
    is_active: boolean;
    project_roles?: UserProjectRole[];
}

// Ответ от сервера при успешной авторизации
export interface AuthResponse {
    token: string;
    user: User;
}

/**
 * Ключи для хранения в sessionStorage
 */
const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

/**
 * Auth Service
 */
export const authService = {
    /**
     * Вход в систему
     * 
     * @param credentials - email и password
     * @returns Promise с токеном и данными пользователя
     * 
     * Backend endpoint: POST /api/auth/login
     * Response: { "token": "...", "user": {...} }
     */
    async login(credentials: LoginRequest): Promise<AuthResponse> {
        try {
            const response = await api.post<AuthResponse>('/auth/login', credentials);
            
            // Автоматически сохраняем токен и данные пользователя после успешного входа
            if (response.data.token) {
                this.setToken(response.data.token);
                this.setUser(response.data.user);
            }
            
            return response.data;
        } catch (error: any) {
            // Логируем ошибку входа
            console.error('[AuthService] Login failed:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Регистрация нового пользователя
     * 
     * @param userData - имя, email, password
     * @returns Promise с токеном и данными пользователя
     * 
     * Backend endpoint: POST /api/auth/register
     * Response: { "token": "...", "user": {...} }
     */
    async register(userData: RegisterRequest): Promise<AuthResponse> {
        try {
            const response = await api.post<AuthResponse>('/auth/register', userData);
            
            // Автоматически сохраняем токен и данные пользователя после успешной регистрации
            if (response.data.token) {
                this.setToken(response.data.token);
                this.setUser(response.data.user);
            }
            
            return response.data;
        } catch (error: any) {
            // Логируем ошибку регистрации
            console.error('[AuthService] Registration failed:', error.response?.data || error.message);
            throw error;
        }
    },

    /**
     * Сохранение JWT токена в sessionStorage
     *
     * @param token - JWT токен
     */
    setToken(token: string): void {
        try {
            sessionStorage.setItem(TOKEN_KEY, token);
            console.log('[AuthService] Token saved');
        } catch (error) {
            console.error('[AuthService] Failed to save token:', error);
        }
    },

    /**
     * Сохранение данных пользователя в sessionStorage
     *
     * @param user - Данные пользователя
     */
    setUser(user: User): void {
        try {
            sessionStorage.setItem(USER_KEY, JSON.stringify(user));
            console.log('[AuthService] User data saved');
        } catch (error) {
            console.error('[AuthService] Failed to save user data:', error);
        }
    },

    /**
     * Получение данных пользователя из sessionStorage
     *
     * @returns Данные пользователя или null
     */
    getUser(): User | null {
        try {
            const userData = sessionStorage.getItem(USER_KEY);
            if (!userData) {
                return null;
            }
            return JSON.parse(userData) as User;
        } catch (error) {
            console.error('[AuthService] Failed to get user data:', error);
            return null;
        }
    },

    /**
     * Получение JWT токена из sessionStorage
     *
     * @returns JWT токен или null если не авторизован
     */
    getToken(): string | null {
        try {
            return sessionStorage.getItem(TOKEN_KEY);
        } catch (error) {
            console.error('[AuthService] Failed to get token:', error);
            return null;
        }
    },

    /**
     * Удаление JWT токена из sessionStorage
     * Вызывается при выходе из системы
     */
    removeToken(): void {
        try {
            sessionStorage.removeItem(TOKEN_KEY);
            sessionStorage.removeItem(USER_KEY);
            console.log('[AuthService] Token and user data removed');
        } catch (error) {
            console.error('[AuthService] Failed to remove token:', error);
        }
    },

    /**
     * Проверка авторизации
     * 
     * @returns true если пользователь авторизован (есть токен)
     */
    isAuthenticated(): boolean {
        const token = this.getToken();
        return token !== null && token.length > 0;
    },

    /**
     * Выход из системы
     * Удаляет токен и очищает данные пользователя
     */
    logout(): void {
        this.removeToken();
        console.log('[AuthService] User logged out');
    },

    /**
     * Декодирование JWT токена (базовое)
     * 
     * ВНИМАНИЕ: Это НЕ валидация токена! Токен может быть просрочен или подделан.
     * Валидацию выполняет только backend.
     * Эта функция нужна только для получения данных из payload.
     * 
     * @param token - JWT токен
     * @returns Payload токена или null при ошибке
     */
    decodeToken(token: string): any | null {
        try {
            // JWT состоит из 3 частей: header.payload.signature
            const parts = token.split('.');
            if (parts.length !== 3) {
                return null;
            }

            // Декодируем payload (вторая часть)
            const payload = parts[1];
            const decoded = atob(payload);
            return JSON.parse(decoded);
        } catch (error) {
            console.error('[AuthService] Failed to decode token:', error);
            return null;
        }
    },

    /**
     * Проверка истечения токена
     * 
     * @param token - JWT токен
     * @returns true если токен истек
     */
    isTokenExpired(token: string): boolean {
        const decoded = this.decodeToken(token);
        if (!decoded || !decoded.exp) {
            return true;
        }

        // exp в JWT хранится в секундах, Date.now() в миллисекундах
        const expirationTime = decoded.exp * 1000;
        const currentTime = Date.now();

        return currentTime >= expirationTime;
    },

    /**
     * Получение user_id из токена
     * 
     * @returns user_id или null
     */
    getUserIdFromToken(): number | null {
        const token = this.getToken();
        if (!token) {
            return null;
        }

        const decoded = this.decodeToken(token);
        return decoded?.user_id || null;
    },
};

