import React from 'react';
import {
    List,
    Datagrid,
    TextField,
    EmailField,
    Edit,
    SimpleForm,
    TextInput,
    Create,
    BooleanField,
    BooleanInput,
    DateField,
    TabbedForm,
    FormTab,
} from 'react-admin';
import { UserProjects } from './UserProjects';

export const UserList = () => (
    <List>
        <Datagrid rowClick="edit">
            <TextField source="id" />
            <TextField source="name" label="Имя" />
            <EmailField source="email" label="Email" />
            <BooleanField source="is_active" label="Активен" />
            <DateField source="created_at" label="Создан" showTime />
            <DateField source="last_login_at" label="Последний вход" showTime />
        </Datagrid>
    </List>
);

export const UserEdit = () => (
    <Edit>
        <TabbedForm>
            <FormTab label="Основная информация">
                <TextInput source="id" disabled />
                <TextInput source="name" label="Имя" required />
                <TextInput source="email" label="Email" required />
                <TextInput 
                    source="password" 
                    label="Пароль" 
                    type="password" 
                    helperText="Оставьте пустым, чтобы не менять пароль. Минимальная длина: 8 символов" 
                    validate={(value) => {
                        if (value && value.length > 0 && value.length < 8) {
                            return 'Пароль должен содержать минимум 8 символов';
                        }
                        return undefined;
                    }}
                />
                <BooleanInput source="is_active" label="Активен" />
            </FormTab>
            <FormTab label="Проекты и роли">
                <UserProjects />
            </FormTab>
        </TabbedForm>
    </Edit>
);

export const UserCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="name" label="Имя" required />
            <TextInput source="email" label="Email" required />
            <TextInput 
                source="password" 
                label="Пароль" 
                type="password" 
                required 
                helperText="Минимальная длина пароля: 8 символов"
                validate={(value) => {
                    if (!value || value.length < 8) {
                        return 'Пароль должен содержать минимум 8 символов';
                    }
                    return undefined;
                }}
            />
            <BooleanInput source="is_active" label="Активен" />
        </SimpleForm>
    </Create>
);

