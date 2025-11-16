import React from 'react';
import { AppBar as RAAppBar, TitlePortal } from 'react-admin';

const AppBar = () => {
    return (
        <RAAppBar>
            <TitlePortal />
            <span style={{ flexGrow: 1, marginLeft: '1rem' }}>
                Planica BI - Админ панель
            </span>
        </RAAppBar>
    );
};

export default AppBar;

