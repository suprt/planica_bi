/**
 * ProtectedRoute Component
 * 
 * Компонент-обертка для защиты роутов от неавторизованных пользователей.
 * Если пользователь не авторизован, происходит редирект на страницу логина.
 * 
 * Использование:
 * ```tsx
 * <Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
 * ```
 */

import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

interface ProtectedRouteProps {
    children: React.ReactElement;
}

/**
 * ProtectedRoute
 * 
 * Проверяет авторизацию пользователя перед рендером дочерних компонентов.
 * 
 * @param children - Защищенный компонент (dashboard, statistics и т.д.)
 */
const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    // Показываем loader пока идет проверка токена
    if (isLoading) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '100vh',
                backgroundColor: '#f5f5f5'
            }}>
                <div style={{
                    textAlign: 'center'
                }}>
                    <div style={{
                        width: '50px',
                        height: '50px',
                        border: '5px solid #e0e0e0',
                        borderTop: '5px solid #3498db',
                        borderRadius: '50%',
                        animation: 'spin 1s linear infinite',
                        margin: '0 auto 16px'
                    }} />
                    <p style={{
                        color: '#666',
                        fontSize: '14px'
                    }}>
                        Проверка авторизации...
                    </p>
                </div>
                
                {/* CSS анимация для спиннера */}
                <style>{`
                    @keyframes spin {
                        0% { transform: rotate(0deg); }
                        100% { transform: rotate(360deg); }
                    }
                `}</style>
            </div>
        );
    }

    // Если не авторизован - редирект на login
    if (!isAuthenticated) {
        console.log('[ProtectedRoute] User not authenticated, redirecting to /login');
        return <Navigate to="/login" replace />;
    }

    // Пользователь авторизован - рендерим защищенный компонент
    return children;
};

export default ProtectedRoute;

