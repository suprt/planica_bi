import React from 'react';
import './Header.css';
import UserMenu from '../UserMenu/UserMenu';
import Notifications from '../Notifications/Notifications';
import SearchBar from '../SearchBar/SearchBar';
import { useCurrentTime } from '../../hooks';

const Header: React.FC = () => {
    const currentTime = useCurrentTime();

    return (
        <header className="header">
            <div className="header-left">
                <div className="time">{currentTime}</div>
            </div>

            <SearchBar />

            <div className="header-right">
                <Notifications />
                <UserMenu />
            </div>
        </header>
    );
};

export default Header;

