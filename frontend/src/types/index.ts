// Типы для пользователя
export interface User {
    id: string;
    name: string;
    position: string;
    avatar: string;
    email?: string;
}

// Типы для уведомлений
export interface Notification {
    id: number;
    text: string;
    time: string;
    read?: boolean;
    type?: 'info' | 'warning' | 'success' | 'error';
}

// Типы для навигации
export interface NavItem {
    id: string;
    label: string;
    path?: string;
    icon: string;
    isMain?: boolean;
    isSettings?: boolean;
}

// Типы для поиска
export interface SearchResult {
    id: string;
    type: 'client' | 'employee' | 'document';
    title: string;
    description?: string;
}
