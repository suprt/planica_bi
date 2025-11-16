/**
 * AuthContext
 * 
 * React Context для управления глобальным состоянием авторизации.
 * Предоставляет информацию о текущем пользователе и методы для входа/выхода.
 */

import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react';
import { authService, User, LoginRequest, RegisterRequest } from '../services/api/authService';

/**
 * Типы для AuthContext
 */
interface AuthContextType {
    // Состояние
    user: User | null;                          // Данные текущего пользователя
    isAuthenticated: boolean;                   // Авторизован ли пользователь
    isLoading: boolean;                         // Идет ли загрузка/проверка токена
    
    // Методы
    login: (email: string, password: string) => Promise<void>;
    register: (name: string, email: string, password: string) => Promise<void>;
    logout: () => void;
}

/**
 * Создаем Context
 */
const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * Props для AuthProvider
 */
interface AuthProviderProps {
    children: ReactNode;
}

/**
 * AuthProvider
 * 
 * Оборачивает всё приложение и предоставляет доступ к состоянию авторизации.
 * 
 * Использование в App.tsx:
 * ```tsx
 * <AuthProvider>
 *   <App />
 * </AuthProvider>
 * ```
 */
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    // Состояние пользователя
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    /**
     * Инициализация при монтировании компонента
     * Проверяем есть ли сохраненный токен и данные пользователя
     */
    useEffect(() => {
        const initializeAuth = async () => {
            try {
                const token = authService.getToken();
                
                if (!token) {
                    // Нет токена - пользователь не авторизован
                    setIsLoading(false);
                    return;
                }

                // Проверяем не истек ли токен
                if (authService.isTokenExpired(token)) {
                    console.log('[AuthContext] Token expired, removing');
                    authService.removeToken();
                    setIsLoading(false);
                    return;
                }

                // Восстанавливаем данные пользователя из localStorage
                const savedUser = authService.getUser();
                
                if (savedUser) {
                    setUser(savedUser);
                    console.log('[AuthContext] User authenticated from localStorage:', savedUser.email);
                } else {
                    // Если данных пользователя нет, удаляем токен (некорректное состояние)
                    console.log('[AuthContext] User data not found, removing token');
                    authService.removeToken();
                }
            } catch (error) {
                console.error('[AuthContext] Initialization error:', error);
                authService.removeToken();
            } finally {
                setIsLoading(false);
            }
        };

        initializeAuth();
    }, []);

    /**
     * Вход в систему
     * 
     * @param email - Email пользователя
     * @param password - Пароль
     * @throws Error если вход не удался
     */
    const login = useCallback(async (email: string, password: string): Promise<void> => {
        try {
            setIsLoading(true);
            
            const credentials: LoginRequest = { email, password };
            const response = await authService.login(credentials);
            
            // Сохраняем данные пользователя
            setUser(response.user);
            
            console.log('[AuthContext] Login successful:', response.user.email);
        } catch (error: any) {
            console.error('[AuthContext] Login failed:', error);
            
            // Пробрасываем ошибку дальше для обработки в компоненте
            throw new Error(
                error.response?.data?.error || 
                error.response?.data?.message || 
                'Ошибка входа в систему'
            );
        } finally {
            setIsLoading(false);
        }
    }, []);

    /**
     * Регистрация нового пользователя
     * 
     * @param name - Имя пользователя
     * @param email - Email
     * @param password - Пароль
     * @throws Error если регистрация не удалась
     */
    const register = useCallback(async (name: string, email: string, password: string): Promise<void> => {
        try {
            setIsLoading(true);
            
            const userData: RegisterRequest = { name, email, password };
            const response = await authService.register(userData);
            
            // Сохраняем данные пользователя
            setUser(response.user);
            
            console.log('[AuthContext] Registration successful:', response.user.email);
        } catch (error: any) {
            console.error('[AuthContext] Registration failed:', error);
            
            // Пробрасываем ошибку дальше для обработки в компоненте
            throw new Error(
                error.response?.data?.error || 
                error.response?.data?.message || 
                'Ошибка регистрации'
            );
        } finally {
            setIsLoading(false);
        }
    }, []);

    /**
     * Выход из системы
     * Удаляет токен и очищает состояние пользователя
     */
    const logout = useCallback((): void => {
        authService.logout();
        setUser(null);
        console.log('[AuthContext] User logged out');
    }, []);

    /**
     * Значение контекста
     */
    const contextValue: AuthContextType = {
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
    };

    return (
        <AuthContext.Provider value={contextValue}>
            {children}
        </AuthContext.Provider>
    );
};

/**
 * Custom Hook для использования AuthContext
 * 
 * Использование в компонентах:
 * ```tsx
 * const { user, isAuthenticated, login, logout } = useAuth();
 * ```
 * 
 * @throws Error если используется вне AuthProvider
 */
export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    
    if (context === undefined) {
        throw new Error('useAuth must be used within AuthProvider');
    }
    
    return context;
};

