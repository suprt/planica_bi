import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Login.css';

const Login: React.FC = () => {
    const navigate = useNavigate();
    const [login, setLogin] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [rememberMe, setRememberMe] = useState<boolean>(false);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Вход:', { login, password, rememberMe });
        // TODO: Реализовать логику входа с проверкой credentials
        // Пока просто редиректим на dashboard
        navigate('/dashboard');
    };

    return (
        <div className="login-page">
            <div className="login-container">
                <div className="login-header">
                    <h1 className="login-title active">Вход</h1>
                    <a href="/register" className="register-link" onClick={(e) => { e.preventDefault(); }}>Регистрация</a>
                </div>

                <form className="login-form" onSubmit={handleSubmit}>
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

                    <button type="submit" className="login-button">
                        Войти
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

