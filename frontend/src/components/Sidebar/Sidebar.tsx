import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './Sidebar.css';
import { navigationItems } from '../../utils/navigation';

interface SidebarProps {
    isOpen: boolean;
    onToggle: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen, onToggle }) => {
    const location = useLocation();

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
                                className={`nav-item ${item.isMain ? 'main-category' : ''} ${item.isSettings ? 'settings-item' : ''} ${location.pathname === item.path ? 'active' : ''}`}
                            >
                                <Link to={item.path} onClick={() => !isOpen && onToggle()}>
                                    {item.label}
                                </Link>
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

