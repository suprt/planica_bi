import { DataProvider, fetchUtils } from 'react-admin';

const httpClient = (url: string, options: any = {}) => {
    // Добавляем JWT токен из sessionStorage
    const token = sessionStorage.getItem('auth_token');
    
    if (!options.headers) {
        options.headers = new Headers();
    }
    
    if (token) {
        options.headers.set('Authorization', `Bearer ${token}`);
    }
    
    // Используем базовый URL из переменных окружения
    const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
    const fullUrl = url.startsWith('http') ? url : `${baseURL}${url}`;
    
    return fetchUtils.fetchJson(fullUrl, options);
};

export const dataProvider: DataProvider = {
    getList: async (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        
        // Формируем URL с параметрами
        let url = `/${resource}?page=${page}&perPage=${perPage}`;
        if (field) {
            url += `&sort=${field}&order=${order}`;
        }
        
        // Добавляем фильтры
        if (Object.keys(params.filter).length > 0) {
            Object.keys(params.filter).forEach(key => {
                url += `&${key}=${params.filter[key]}`;
            });
        }
        
        const { json } = await httpClient(url);
        
        // Если ответ уже в формате react-admin (с data и total)
        if (json && typeof json === 'object' && json.data && json.total !== undefined) {
            return {
                data: json.data,
                total: json.total,
            };
        }
        
        // Если ответ - просто массив, оборачиваем его
        if (Array.isArray(json)) {
            return {
                data: json,
                total: json.length,
            };
        }
        
        return {
            data: [],
            total: 0,
        };
    },

    getOne: async (resource, params) => {
        const { json } = await httpClient(`/${resource}/${params.id}`);
        
        // Если ответ уже в формате react-admin (с data)
        if (json.data) {
            return { data: json.data };
        }
        
        // Если ответ - просто объект
        return { data: json };
    },

    getMany: async (resource, params) => {
        const promises = params.ids.map(id => 
            httpClient(`/${resource}/${id}`).then(({ json }) => {
                if (json.data) return json.data;
                return json;
            })
        );
        const data = await Promise.all(promises);
        return { data };
    },

    getManyReference: async (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        
        // Формируем URL с параметрами
        let url = `/${resource}?page=${page}&perPage=${perPage}`;
        if (field) {
            url += `&sort=${field}&order=${order}`;
        }
        
        // Добавляем фильтр по связанному ресурсу
        url += `&${params.target}=${params.id}`;
        
        // Добавляем другие фильтры
        if (Object.keys(params.filter).length > 0) {
            Object.keys(params.filter).forEach(key => {
                url += `&${key}=${params.filter[key]}`;
            });
        }
        
        const { json } = await httpClient(url);
        
        // Если ответ уже в формате react-admin (с data и total)
        if (json.data && json.total !== undefined) {
            return {
                data: json.data,
                total: json.total,
            };
        }
        
        // Если ответ - просто массив
        if (Array.isArray(json)) {
            return {
                data: json,
                total: json.length,
            };
        }
        
        return {
            data: [],
            total: 0,
        };
    },

    create: async (resource, params) => {
        // Фильтруем данные перед отправкой - убираем поля, которые не нужны бэкенду
        const dataToSend = { ...params.data };
        
        // Для users - убираем поля, которые не нужны при создании
        if (resource === 'users') {
            delete dataToSend.timezone;
            delete dataToSend.language;
            delete dataToSend.created_at;
            delete dataToSend.updated_at;
            delete dataToSend.last_login_at;
            delete dataToSend.id; // ID генерируется на бэкенде
        }
        
        // Для projects - убираем поля, которые не нужны при создании
        if (resource === 'projects') {
            delete dataToSend.created_at;
            delete dataToSend.updated_at;
            delete dataToSend.public_token; // Генерируется на бэкенде
            delete dataToSend.id; // ID генерируется на бэкенде
        }
        
        const { json } = await httpClient(`/${resource}`, {
            method: 'POST',
            body: JSON.stringify(dataToSend),
        });
        
        // Если ответ уже в формате react-admin (с data)
        if (json && typeof json === 'object' && json.data) {
            return { data: json.data };
        }
        
        // Если ответ - просто объект
        return { data: json };
    },

    update: async (resource, params) => {
        // Фильтруем данные перед отправкой - убираем поля, которые не нужны бэкенду
        const dataToSend = { ...params.data };
        
        // Для users - убираем поля, которые не нужны при обновлении
        if (resource === 'users') {
            delete dataToSend.timezone;
            delete dataToSend.language;
            delete dataToSend.created_at;
            delete dataToSend.updated_at;
            delete dataToSend.last_login_at;
            // Если пароль пустой, удаляем его (бэкенд не должен обновлять пароль)
            if (!dataToSend.password || dataToSend.password.trim() === '') {
                delete dataToSend.password;
            }
        }
        
        // Для projects - убираем поля, которые не нужны при обновлении
        if (resource === 'projects') {
            delete dataToSend.created_at;
            delete dataToSend.updated_at;
            delete dataToSend.public_token; // Не обновляется
            delete dataToSend.id; // ID не обновляется
        }
        
        const { json } = await httpClient(`/${resource}/${params.id}`, {
            method: 'PUT',
            body: JSON.stringify(dataToSend),
        });
        
        // Если ответ уже в формате react-admin (с data)
        if (json && typeof json === 'object' && json.data) {
            return { data: json.data };
        }
        
        // Если ответ - просто объект
        return { data: json };
    },

    updateMany: async (resource, params) => {
        const promises = params.ids.map(id =>
            httpClient(`/${resource}/${id}`, {
                method: 'PUT',
                body: JSON.stringify(params.data),
            }).then(({ json }) => {
                if (json.data) return json.data;
                return json;
            })
        );
        const data = await Promise.all(promises);
        return { data };
    },

    delete: async (resource, params) => {
        const { json } = await httpClient(`/${resource}/${params.id}`, {
            method: 'DELETE',
        });
        
        // Если ответ уже в формате react-admin (с data)
        if (json && typeof json === 'object' && json.data) {
            return { data: json.data };
        }
        
        // Если ответ - просто объект с id
        return { data: { id: params.id } };
    },

    deleteMany: async (resource, params) => {
        const promises = params.ids.map(id =>
            httpClient(`/${resource}/${id}`, {
                method: 'DELETE',
            })
        );
        await Promise.all(promises);
        return { data: params.ids };
    },
};

