import { AuthProvider } from 'react-admin';
import { authService } from '../services/api/authService';

export const authProvider: AuthProvider = {
    login: async ({ username, password }) => {
        try {
            await authService.login({ email: username, password });
            
            // authService уже сохраняет токен в sessionStorage
            // Проверяем, что пользователь - админ
            const user = authService.getUser();
            if (!user) {
                throw new Error('Не удалось получить данные пользователя');
            }
            
            // Проверяем, есть ли у пользователя роль admin
            // Для этого нужно получить полные данные пользователя с проектами
            // Используем API для получения полных данных
            let isAdmin = false;
            try {
                const fullUserResponse = await fetch(`${process.env.REACT_APP_API_URL || 'http://localhost:8080/api'}/users/${user.id}`, {
                    headers: {
                        'Authorization': `Bearer ${authService.getToken()}`,
                    },
                });
                if (fullUserResponse.ok) {
                    const fullUserData = await fullUserResponse.json();
                    isAdmin = fullUserData.data?.projects?.some((p: any) => p.role === 'admin') || false;
                } else {
                    // Если не удалось получить данные, проверяем через project_roles
                    isAdmin = user.project_roles?.some((role: any) => role.role === 'admin') || false;
                }
            } catch (err: any) {
                // Если ошибка при получении данных, проверяем через project_roles
                isAdmin = user.project_roles?.some((role: any) => role.role === 'admin') || false;
            }
            
            if (!isAdmin) {
                // Очищаем токен, если пользователь не админ
                authService.logout();
                throw new Error('Доступ запрещен. Требуется роль администратора.');
            }
            
            return Promise.resolve();
        } catch (error: any) {
            return Promise.reject(error.message || 'Ошибка входа');
        }
    },
    
    logout: () => {
        authService.logout();
        return Promise.resolve();
    },
    
    checkError: ({ status }: { status: number }) => {
        if (status === 401 || status === 403) {
            authService.logout();
            return Promise.reject();
        }
        return Promise.resolve();
    },
    
    checkAuth: async () => {
        console.log('[authProvider] checkAuth called');
        const token = authService.getToken();
        console.log('[authProvider] Token exists:', !!token);
        if (!token) {
            console.log('[authProvider] No token, rejecting');
            return Promise.reject();
        }
        if (authService.isTokenExpired(token)) {
            console.log('[authProvider] Token expired, rejecting');
            return Promise.reject();
        }
        
        // Проверяем, что пользователь - админ
        const user = authService.getUser();
        console.log('[authProvider] User:', user);
        if (!user) {
            console.log('[authProvider] No user, rejecting');
            return Promise.reject();
        }
        
        // Проверяем админ-роль через API
        try {
            const fullUserResponse = await fetch(`${process.env.REACT_APP_API_URL || 'http://localhost:8080/api'}/users/${user.id}`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                },
            });
            if (fullUserResponse.ok) {
                const fullUserData = await fullUserResponse.json();
                const isAdmin = fullUserData.data?.projects?.some((p: any) => p.role === 'admin') || false;
                console.log('[authProvider] Is admin (from API):', isAdmin);
                if (!isAdmin) {
                    console.log('[authProvider] Not admin, rejecting');
                    return Promise.reject('Доступ запрещен. Требуется роль администратора.');
                }
            } else {
                // Если не удалось получить данные, проверяем через project_roles
                const isAdmin = user.project_roles?.some((role: any) => role.role === 'admin') || false;
                console.log('[authProvider] Is admin (from project_roles):', isAdmin);
                if (!isAdmin) {
                    console.log('[authProvider] Not admin, rejecting');
                    return Promise.reject('Доступ запрещен. Требуется роль администратора.');
                }
            }
        } catch (err) {
            console.error('[authProvider] Error checking admin status:', err);
            // Если ошибка, проверяем через project_roles
            const isAdmin = user.project_roles?.some((role: any) => role.role === 'admin') || false;
            console.log('[authProvider] Is admin (fallback):', isAdmin);
            if (!isAdmin) {
                console.log('[authProvider] Not admin, rejecting');
                return Promise.reject('Доступ запрещен. Требуется роль администратора.');
            }
        }
        
        console.log('[authProvider] Auth check passed');
        return Promise.resolve();
    },
    
    getPermissions: async () => {
        const user = authService.getUser();
        if (!user) {
            return Promise.reject();
        }
        
        const token = authService.getToken();
        try {
            const fullUserResponse = await fetch(`${process.env.REACT_APP_API_URL || 'http://localhost:8080/api'}/users/${user.id}`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                },
            });
            if (fullUserResponse.ok) {
                const fullUserData = await fullUserResponse.json();
                const isAdmin = fullUserData.data?.projects?.some((p: any) => p.role === 'admin') || false;
                return Promise.resolve(isAdmin ? 'admin' : 'user');
            }
        } catch (err) {
            // Если ошибка, проверяем через project_roles
        }
        
        const isAdmin = user.project_roles?.some((role: any) => role.role === 'admin') || false;
        return Promise.resolve(isAdmin ? 'admin' : 'user');
    },
    
    getIdentity: async () => {
        const user = authService.getUser();
        if (!user) {
            return Promise.reject();
        }
        
        return Promise.resolve({
            id: user.id,
            fullName: user.name,
            avatar: undefined,
        });
    },
};

