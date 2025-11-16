import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import './Login.css';

const Login: React.FC = () => {
    const navigate = useNavigate();
    const { login: authLogin, isLoading } = useAuth();
    
    // Form state
    const [login, setLogin] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [rememberMe, setRememberMe] = useState<boolean>(false);
    
    // Error state
    const [error, setError] = useState<string>('');

    /**
     * Обработка отправки формы
     * Вызывает API login через AuthContext
     */
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        // Сбрасываем предыдущую ошибку
        setError('');
        
        // Базовая валидация
        if (!login.trim() || !password.trim()) {
            setError('Заполните все поля');
            return;
        }

        try {
            // Вызываем login из AuthContext
            // email = login (в форме может быть и логин, и email)
            await authLogin(login.trim(), password);
            
            // Успешный вход - редиректим на dashboard
            console.log('[Login] Success, redirecting to dashboard');
            navigate('/dashboard');
        } catch (err: any) {
            // Обрабатываем ошибку
            console.error('[Login] Error:', err);
            
            const errorMessage = err.message || 'Ошибка входа в систему. Проверьте данные и попробуйте снова.';
            setError(errorMessage);
        }
    };

    return (
        <div className="login-page">
            <div className="login-container">
                <div className="login-header">
                    <h1 className="login-title active">Вход</h1>
                    <a href="/register" className="register-link" onClick={(e) => { e.preventDefault(); }}>Регистрация</a>
                </div>

                <form className="login-form" onSubmit={handleSubmit}>
                    {/* Сообщение об ошибке */}
                    {error && (
                        <div className="error-message" style={{
                            padding: '12px',
                            marginBottom: '16px',
                            backgroundColor: '#fee',
                            border: '1px solid #fcc',
                            borderRadius: '4px',
                            color: '#c33',
                            fontSize: '14px'
                        }}>
                            {error}
                        </div>
                    )}
                    
                    <div className="form-group">
                        <label htmlFor="login" className="form-label">
                            Логин или e-mail
                        </label>
                        <input
                            type="text"
                            id="login"
                            className="form-input"
                            value={login}
                            onChange={(e) => setLogin(e.target.value)}
                            placeholder=""
                            required
                            disabled={isLoading}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password" className="form-label">
                            Пароль
                        </label>
                        <input
                            type="password"
                            id="password"
                            className="form-input"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder=""
                            required
                            disabled={isLoading}
                        />
                    </div>

                    <div className="form-options">
                        <label className="checkbox-label">
                            <input
                                type="checkbox"
                                checked={rememberMe}
                                onChange={(e) => setRememberMe(e.target.checked)}
                                className="checkbox-input"
                            />
                            <span className="checkbox-text">Запомнить меня</span>
                        </label>
                    </div>

                    <button 
                        type="submit" 
                        className="login-button"
                        disabled={isLoading}
                        style={{
                            opacity: isLoading ? 0.6 : 1,
                            cursor: isLoading ? 'not-allowed' : 'pointer'
                        }}
                    >
                        {isLoading ? 'Вход...' : 'Войти'}
                    </button>

                    <div className="login-links">
                        <a href="/forgot-password" className="link" onClick={(e) => { e.preventDefault(); }}>Напомнить пароль</a>
                        <a href="/resend-email" className="link" onClick={(e) => { e.preventDefault(); }}>Не пришло письмо?</a>
                    </div>
                </form>

                <div className="social-login">
                    <div className="social-login-header">
                        <span className="social-login-text">Войти с помощью</span>
                    </div>
                    <button 
                        className="yandex-id-button" 
                        onClick={() => {
                            console.log('Вход через Яндекс ID');
                            // TODO: Реализовать OAuth авторизацию через Яндекс ID
                            navigate('/dashboard');
                        }}
                        aria-label="Войти через Яндекс ID"
                    >
                        <span className="yandex-icon">Я</span>
                        <span className="yandex-text">ID</span>
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Login;

