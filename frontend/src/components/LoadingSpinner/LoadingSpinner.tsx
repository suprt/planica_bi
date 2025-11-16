/**
 * LoadingSpinner Component
 * 
 * Универсальный компонент спиннера загрузки.
 * Можно использовать для индикации загрузки данных.
 */

import React from 'react';
import './LoadingSpinner.css';

interface LoadingSpinnerProps {
    message?: string;        // Сообщение под спиннером
    size?: 'small' | 'medium' | 'large'; // Размер спиннера
    fullScreen?: boolean;    // Показывать на весь экран
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ 
    message = 'Загрузка...', 
    size = 'medium',
    fullScreen = false 
}) => {
    const sizeClasses = {
        small: 'spinner-small',
        medium: 'spinner-medium',
        large: 'spinner-large',
    };

    const containerClass = fullScreen ? 'loading-container-fullscreen' : 'loading-container';

    return (
        <div className={containerClass}>
            <div className={`spinner ${sizeClasses[size]}`} role="status" aria-live="polite">
                <div className="spinner-circle"></div>
            </div>
            {message && <p className="loading-message">{message}</p>}
        </div>
    );
};

export default LoadingSpinner;

