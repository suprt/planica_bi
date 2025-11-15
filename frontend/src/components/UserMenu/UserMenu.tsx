import React, { useState, useRef } from 'react';
import './UserMenu.css';
import { User } from '../../types';
import { useClickOutside } from '../../hooks';

interface UserMenuProps {
    user?: User;
}

const UserMenu: React.FC<UserMenuProps> = ({ user }) => {
    const [userMenuOpen, setUserMenuOpen] = useState<boolean>(false);
    const userMenuRef = useRef<HTMLDivElement>(null);

    const defaultUser: User = {
        id: '1',
        name: 'Иван Иванов',
        position: 'Менеджер по продажам',
        avatar: 'ИИ',
    };

    const currentUser = user || defaultUser;

    useClickOutside(userMenuRef, () => setUserMenuOpen(false));

    const toggleUserMenu = () => {
        setUserMenuOpen(!userMenuOpen);
    };

    const handleLogout = () => {
        console.log('Выход из системы');
        setUserMenuOpen(false);
    };

    const handleProfile = () => {
        console.log('Переход в профиль');
        setUserMenuOpen(false);
    };

    return (
        <div className="user-menu" ref={userMenuRef}>
            <button className="user-button" onClick={toggleUserMenu}>
                <div className="user-avatar">
                    {currentUser.avatar}
                </div>
                <span className="user-name">{currentUser.name}</span>
            </button>

            {userMenuOpen && (
                <div className="user-dropdown">
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
                        <button className="user-menu-item">
                            <span className="user-menu-icon">🌙</span>
                            Темная тема
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
    );
};

export default UserMenu;

