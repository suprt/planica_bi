/**
 * ErrorMessage Component
 * 
 * Универсальный компонент для отображения ошибок.
 * Показывает сообщение об ошибке и опциональную кнопку повтора.
 */

import React from 'react';
import './ErrorMessage.css';

interface ErrorMessageProps {
    message: string;         // Текст ошибки
    title?: string;          // Заголовок (опционально)
    onRetry?: () => void;    // Callback для повторной попытки
    type?: 'error' | 'warning' | 'info'; // Тип сообщения
    fullPage?: boolean;      // Показывать на всю страницу
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ 
    message, 
    title,
    onRetry, 
    type = 'error',
    fullPage = false 
}) => {
    const iconMap = {
        error: '❌',
        warning: '⚠️',
        info: 'ℹ️',
    };

    const containerClass = fullPage ? 'error-container-fullpage' : 'error-container';

    return (
        <div className={`${containerClass} error-${type}`}>
            <div className="error-content">
                <div className="error-icon" aria-hidden="true">
                    {iconMap[type]}
                </div>
                <div className="error-text">
                    {title && <h3 className="error-title">{title}</h3>}
                    <p className="error-message">{message}</p>
                </div>
            </div>
            {onRetry && (
                <button 
                    className="error-retry-button" 
                    onClick={onRetry}
                    aria-label="Повторить попытку"
                >
                    Повторить попытку
                </button>
            )}
        </div>
    );
};

export default ErrorMessage;

