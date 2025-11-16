import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import './Sidebar.css';
import { navigationItems } from '../../utils/navigation';
import { projectStorage } from '../../utils/projectStorage';

interface SidebarProps {
    isOpen: boolean;
    onToggle: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen, onToggle }) => {
    const location = useLocation();
    const navigate = useNavigate();

    /**
     * Проверка, активен ли элемент навигации
     */
    const isItemActive = (itemId: string): boolean => {
        if (itemId === 'projects') {
            return location.pathname === '/dashboard/projects';
        } else if (itemId === 'statistics') {
            return location.pathname === '/dashboard/statistics';
        } else if (itemId === 'reports') {
            return location.pathname === '/dashboard/reports';
        } else if (itemId === 'metrics') {
            return location.pathname === '/dashboard/metrics';
        } else {
            return false;
        }
    };

    /**
     * Обработчик клика по элементу навигации
     */
    const handleNavClick = (itemId: string) => {
        if (itemId === 'statistics') {
            const lastProjectId = projectStorage.getLastProject();
            if (lastProjectId) {
                navigate(`/dashboard/statistics?project=${lastProjectId}`);
            } else {
                navigate('/dashboard/projects');
            }
        } else if (itemId === 'reports') {
            const lastProjectId = projectStorage.getLastProject();
            if (lastProjectId) {
                navigate(`/dashboard/reports?project=${lastProjectId}`);
            } else {
                navigate('/dashboard/projects');
            }
        } else if (itemId === 'metrics') {
            const lastProjectId = projectStorage.getLastProject();
            if (lastProjectId) {
                navigate(`/dashboard/metrics?project=${lastProjectId}`);
            } else {
                navigate('/dashboard/projects');
            }
        } else if (itemId === 'projects') {
            navigate('/dashboard/projects');
        } else {
            navigate('/dashboard');
        }
        
        if (!isOpen) {
            onToggle();
        }
    };

    return (
        <>
            <aside className={`sidebar ${isOpen ? 'sidebar-open' : 'sidebar-closed'}`}>
                <nav className="navigation">
                    <div className="sidebar-header">
                        <h1 className="title">Planica</h1>
                        <button className="menu-toggle" onClick={onToggle}>
                            <div className={`hamburger ${isOpen ? 'hamburger-open' : ''}`}>
                                <span></span>
                                <span></span>
                                <span></span>
                            </div>
                        </button>
                    </div>
                    <ul className="nav-list">
                        {navigationItems.map(item => (
                            <li
                                key={item.id}
                                className={`nav-item ${item.isSettings ? 'settings-item' : ''} ${isItemActive(item.id) ? 'active' : ''}`}
                                onClick={() => handleNavClick(item.id)}
                            >
                                <span className="nav-icon">{item.icon}</span>
                                <span className="nav-label">{item.label}</span>
                            </li>
                        ))}
                    </ul>
                </nav>
            </aside>
            {isOpen && <div className="sidebar-overlay" onClick={onToggle} />}
        </>
    );
};

export default Sidebar;

