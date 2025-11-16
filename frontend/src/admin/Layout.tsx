import React from 'react';
import { Layout as RALayout, LayoutProps } from 'react-admin';
import AppBar from './AppBar';
import Menu from './Menu';

const Layout: React.FC<LayoutProps> = (props) => {
    return <RALayout {...props} appBar={AppBar} menu={Menu} />;
};

export default Layout;

