import React from 'react';
import {
    List,
    Datagrid,
    TextField,
    Edit,
    SimpleForm,
    TextInput,
    Create,
    BooleanField,
    BooleanInput,
    DateField,
} from 'react-admin';

export const ProjectList = () => (
    <List>
        <Datagrid rowClick="edit">
            <TextField source="id" />
            <TextField source="name" label="Название" />
            <TextField source="slug" label="Slug" />
            <TextField source="public_token" label="Публичный токен" />
            <TextField source="timezone" label="Часовой пояс" />
            <TextField source="currency" label="Валюта" />
            <BooleanField source="is_active" label="Активен" />
            <DateField source="created_at" label="Создан" showTime />
        </Datagrid>
    </List>
);

export const ProjectEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="id" disabled />
            <TextInput source="name" label="Название" required />
            <TextInput source="slug" label="Slug" required />
            <TextInput source="public_token" label="Публичный токен" />
            <TextInput source="timezone" label="Часовой пояс" />
            <TextInput source="currency" label="Валюта" />
            <BooleanInput source="is_active" label="Активен" />
        </SimpleForm>
    </Edit>
);

export const ProjectCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="name" label="Название" required />
            <TextInput source="slug" label="Slug" required />
            <TextInput source="timezone" label="Часовой пояс" defaultValue="Europe/Moscow" />
            <TextInput source="currency" label="Валюта" defaultValue="RUB" />
            <BooleanInput source="is_active" label="Активен" />
        </SimpleForm>
    </Create>
);

