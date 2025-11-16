import React from 'react';
import { useLocation } from 'react-router-dom';
import './Placeholder.css';

const Placeholder: React.FC = () => {
    const location = useLocation();
    
    // Извлекаем название раздела из пути
    const pathParts = location.pathname.split('/').filter(Boolean);
    const section = pathParts[pathParts.length - 1]; // Последняя часть пути
    
    // Маппинг ID разделов на читаемые названия
    const sectionNames: Record<string, string> = {
        sources: 'Источники',
        purchases: 'Закупки',
        tasks: 'Задачи и проекты',
        resources: 'Ресурсы',
        finance: 'Финансы',
        logistics: 'Логистика',
        innovation: 'Инноватика',
        production: 'Производство',
        company: 'Компания',
        documents: 'Документы',
        settings: 'Настройки',
    };
    
    const sectionName = sectionNames[section || ''] || 'Раздел';
    
    return (
        <div className="placeholder-page">
            <div className="placeholder-content">
                <div className="placeholder-icon">🚧</div>
                <h1 className="placeholder-title">{sectionName}</h1>
                <p className="placeholder-message">
                    Раздел "{sectionName}" находится в разработке и будет доступен в ближайшее время.
                </p>
            </div>
        </div>
    );
};

export default Placeholder;

