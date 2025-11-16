import React from 'react';
import { Admin, Resource } from 'react-admin';
import { dataProvider } from './dataProvider';
import { authProvider } from './authProvider';
import { UserList, UserEdit, UserCreate } from './resources/Users';
import { ProjectList, ProjectEdit, ProjectCreate } from './resources/Projects';
import Layout from './Layout';
import LoginPage from './LoginPage';

const AdminApp: React.FC = () => {
    return (
        <Admin
            dataProvider={dataProvider}
            authProvider={authProvider}
            layout={Layout}
            loginPage={LoginPage}
            requireAuth
            basename="/admin"
        >
            <Resource
                name="users"
                list={UserList}
                edit={UserEdit}
                create={UserCreate}
            />
            <Resource
                name="projects"
                list={ProjectList}
                edit={ProjectEdit}
                create={ProjectCreate}
            />
        </Admin>
    );
};

export default AdminApp;
