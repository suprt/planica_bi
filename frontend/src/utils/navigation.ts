import { NavItem } from '../types';

export const navigationItems: NavItem[] = [
    { id: 'statistics', label: 'Статистика', icon: '📊' },
    { id: 'sources', label: 'Источники', icon: '🔄' },
    { id: 'purchases', label: 'Закупки', icon: '🛒' },
    { id: 'tasks', label: 'Задачи и проекты', icon: '📋' },
    { id: 'resources', label: 'Ресурсы', icon: '💿' },
    { id: 'finance', label: 'Финансы', icon: '💰' },
    { id: 'logistics', label: 'Логистика', icon: '🚚' },
    { id: 'innovation', label: 'Инноватика', icon: '💡' },
    { id: 'production', label: 'Производство', icon: '🏭' },
    { id: 'company', label: 'Компания', icon: '🏢' },
    { id: 'marketing', label: 'Маркетинг', icon: '📢' },
    { id: 'documents', label: 'Документы', icon: '📄' },
    { id: 'settings', label: 'Настройки', icon: '⚙️', isSettings: true },
];
