import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { useNavigate, useLocation, Outlet } from 'react-router-dom';
import { useTheme } from '../../contexts/ThemeContext';
import { useCurrentTime, useClickOutside } from '../../hooks';
import { navigationItems } from '../../utils/navigation';
import { NavItem as NavItemType } from '../../types';
import '../../App.css';

// Мемоизированный компонент уведомления
const NotificationItem = React.memo<{ id: number; text: string; time: string }>(
    ({ text, time }) => (
        <div className="notification-item">
            <div className="notification-text">{text}</div>
            <div className="notification-time">{time}</div>
        </div>
    )
);
NotificationItem.displayName = 'NotificationItem';

// Мемоизированный компонент пункта меню
const NavItem = React.memo<{ 
    item: NavItemType; 
    isActive: boolean; 
    onClick: () => void;
}>(
    ({ item, isActive, onClick }) => (
        <li 
            className={`nav-item ${item.isSettings ? 'settings-item' : ''} ${isActive ? 'active' : ''}`}
            onClick={onClick}
        >
            <span className="nav-icon">{item.icon}</span>
            <span className="nav-label">{item.label}</span>
        </li>
    )
);
NavItem.displayName = 'NavItem';

const Dashboard: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { theme, toggleTheme } = useTheme();
    const currentTime = useCurrentTime();
    const [sidebarOpen, setSidebarOpen] = useState<boolean>(true);
    const [notificationsOpen, setNotificationsOpen] = useState<boolean>(false);
    const [userMenuOpen, setUserMenuOpen] = useState<boolean>(false);
    const [activeNavItem, setActiveNavItem] = useState<string>(() => {
        // Определяем активный элемент на основе текущего пути
        if (location.pathname.includes('/statistics')) return 'statistics';
        return '';
    });
    const notificationsRef = useRef<HTMLDivElement>(null);
    const userMenuRef = useRef<HTMLDivElement>(null);

    // Пример данных пользователя
    const currentUser = useMemo(() => ({
        name: 'Иван Иванов',
        position: 'Менеджер по продажам',
        avatar: 'ИИ'
    }), []);

    // Пример данных уведомлений
    const notifications = useMemo(() => [
        { id: 1, text: 'Новая задача от Ивана', time: '10:30' },
        { id: 2, text: 'Заказ №2456 выполнен', time: '09:45' },
        { id: 3, text: 'Поступил новый отзыв', time: '09:15' },
        { id: 4, text: 'Обновление системы', time: 'Вчера' },
    ], []);

    // Обновление активного элемента при изменении пути
    useEffect(() => {
        if (location.pathname.includes('/statistics')) {
            setActiveNavItem('statistics');
        } else {
            setActiveNavItem('');
        }
    }, [location.pathname]);

    // Закрытие выпадающих списков при клике вне их области
    useClickOutside(notificationsRef, () => setNotificationsOpen(false));
    useClickOutside(userMenuRef, () => setUserMenuOpen(false));

    const toggleSidebar = useCallback(() => {
        setSidebarOpen(prev => !prev);
    }, []);

    const toggleNotifications = useCallback(() => {
        setNotificationsOpen(prev => {
            if (!prev) {
                setUserMenuOpen(false);
            }
            return !prev;
        });
    }, []);

    const toggleUserMenu = useCallback(() => {
        setUserMenuOpen(prev => {
            if (!prev) {
                setNotificationsOpen(false);
            }
            return !prev;
        });
    }, []);

    const handleLogout = useCallback(() => {
        console.log('Выход из системы');
        setUserMenuOpen(false);
        // TODO: Реализовать очистку токенов/сессии
        navigate('/login');
    }, [navigate]);

    const handleProfile = useCallback(() => {
        console.log('Переход в профиль');
        setUserMenuOpen(false);
    }, []);

    const handleThemeToggle = useCallback(() => {
        toggleTheme();
    }, [toggleTheme]);

    const handleNavItemClick = useCallback((itemId: string) => {
        setActiveNavItem(itemId);
        // Навигация по разделам
        if (itemId === 'statistics') {
            navigate('/dashboard/statistics');
        } else {
            navigate('/dashboard');
        }
    }, [navigate]);

    return (
        <div className="app">
            <header className="header">
                <div className="header-left">
                    <button 
                        className="menu-toggle" 
                        onClick={toggleSidebar}
                        aria-label="Переключить меню"
                    >
                        <div className={`hamburger ${sidebarOpen ? 'hamburger-open' : ''}`}>
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                    </button>
                    <h1 className="title">Planica</h1>
                </div>

                <div className="search-bar">
                    <span className="search-icon" aria-hidden="true">🔍</span>
                    <input
                        type="text"
                        placeholder="Искать клиента, сотрудника, документ"
                        className="search-input"
                        aria-label="Поиск"
                    />
                </div>

                <div className="header-right">
                    <div className="time" aria-label={`Текущее время: ${currentTime}`}>
                        {currentTime}
                    </div>

                    <div className="notifications" ref={notificationsRef}>
                        <button 
                            className="notification-bell" 
                            onClick={toggleNotifications}
                            aria-label="Уведомления"
                            aria-expanded={notificationsOpen}
                        >
                            🔔
                            {notifications.length > 0 && (
                                <span className="notification-badge" aria-label={`${notifications.length} непрочитанных уведомлений`}>
                                    {notifications.length}
                                </span>
                            )}
                        </button>

                        {notificationsOpen && (
                            <div className="notifications-dropdown" role="menu">
                                <div className="notifications-header">
                                    <h3>Уведомления</h3>
                                    <span className="notifications-count">{notifications.length}</span>
                                </div>

                                <div className="notifications-list">
                                    {notifications.map(notification => (
                                        <NotificationItem
                                            key={notification.id}
                                            id={notification.id}
                                            text={notification.text}
                                            time={notification.time}
                                        />
                                    ))}
                                </div>

                                <div className="notifications-footer">
                                    <button className="view-all-btn">Показать все</button>
                                </div>
                            </div>
                        )}
                    </div>

                    <button 
                        className="theme-toggle-button" 
                        onClick={handleThemeToggle}
                        aria-label={theme === 'dark' ? 'Переключить на светлую тему' : 'Переключить на тёмную тему'}
                        title={theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}
                    >
                        {theme === 'dark' ? '☀️' : '🌙'}
                    </button>

                    <div className="user-menu" ref={userMenuRef}>
                        <button 
                            className="user-button" 
                            onClick={toggleUserMenu}
                            aria-label="Меню пользователя"
                            aria-expanded={userMenuOpen}
                        >
                            <div className="user-avatar">
                                {currentUser.avatar}
                            </div>
                            <span className="user-name">{currentUser.name}</span>
                        </button>

                        {userMenuOpen && (
                            <div className="user-dropdown" role="menu">
                                <div className="user-info">
                                    <div className="user-avatar-large">
                                        {currentUser.avatar}
                                    </div>
                                    <div className="user-details">
                                        <div className="user-name-large">{currentUser.name}</div>
                                        <div className="user-position">{currentUser.position}</div>
                                    </div>
                                </div>

                                <div className="user-menu-items">
                                    <button className="user-menu-item" onClick={handleProfile}>
                                        <span className="user-menu-icon">👤</span>
                                        Мой профиль
                                    </button>
                                    <button className="user-menu-item">
                                        <span className="user-menu-icon">⚙️</span>
                                        Настройки аккаунта
                                    </button>
                                    <div className="user-menu-divider"></div>
                                    <button className="user-menu-item logout" onClick={handleLogout}>
                                        <span className="user-menu-icon">🚪</span>
                                        Выйти
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </header>

            <div className="main-container">
                <aside className={`sidebar ${sidebarOpen ? 'sidebar-open' : 'sidebar-closed'}`}>
                    <nav className="navigation" aria-label="Основная навигация">
                        <ul className="nav-list">
                            {navigationItems.map((item) => (
                                <NavItem
                                    key={item.id}
                                    item={item}
                                    isActive={activeNavItem === item.id}
                                    onClick={() => handleNavItemClick(item.id)}
                                />
                            ))}
                        </ul>
                    </nav>
                </aside>

                <main className="main-content">
                    <Outlet />
                </main>
            </div>
        </div>
    );
};

export default Dashboard;

