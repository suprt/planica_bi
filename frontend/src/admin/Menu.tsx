import React from 'react';
import { Menu as RAMenu } from 'react-admin';

const Menu = () => {
    return (
        <RAMenu>
            <RAMenu.Item
                to="/users"
                primaryText="Пользователи"
            />
            <RAMenu.Item
                to="/projects"
                primaryText="Проекты"
            />
        </RAMenu>
    );
};

export default Menu;

