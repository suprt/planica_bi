import React, { useState, useEffect } from 'react';
import {
    useDataProvider,
    useNotify,
    useRecordContext,
    useRefresh,
} from 'react-admin';

interface UserProject {
    project_id: number;
    project_name: string;
    role: 'admin' | 'manager' | 'client';
}

interface Project {
    id: number;
    name: string;
}

export const UserProjects: React.FC = () => {
    const record = useRecordContext();
    const dataProvider = useDataProvider();
    const notify = useNotify();
    const refresh = useRefresh();

    const [projects, setProjects] = useState<UserProject[]>([]);
    const [allProjects, setAllProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [editingProjectId, setEditingProjectId] = useState<number | null>(null);
    const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
    const [selectedRole, setSelectedRole] = useState<'admin' | 'manager' | 'client'>('client');

    // Загрузка проектов пользователя
    useEffect(() => {
        if (!record?.id) return;

        const loadData = async () => {
            try {
                setLoading(true);
                // Загружаем пользователя с проектами
                const userResponse = await dataProvider.getOne('users', { id: record.id });
                const user = userResponse.data;
                
                if (user.projects) {
                    setProjects(user.projects || []);
                }

                // Загружаем все проекты для выбора
                const projectsResponse = await dataProvider.getList('projects', {
                    pagination: { page: 1, perPage: 1000 },
                    sort: { field: 'name', order: 'ASC' },
                    filter: {},
                });
                setAllProjects(projectsResponse.data || []);
            } catch (error: any) {
                console.error('Failed to load user projects:', error);
                notify('Ошибка загрузки проектов', { type: 'error' });
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [record?.id, dataProvider, notify]);

    // Добавить проект
    const handleAddProject = async () => {
        if (!selectedProjectId || !record?.id) {
            notify('Выберите проект', { type: 'warning' });
            return;
        }

        try {
            const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
            const token = sessionStorage.getItem('auth_token');

            const response = await fetch(`${apiUrl}/projects/${selectedProjectId}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    user_id: record.id,
                    role: selectedRole,
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Ошибка назначения проекта');
            }

            notify('Проект успешно добавлен', { type: 'success' });
            setSelectedProjectId(null);
            setSelectedRole('client');
            
            // Перезагружаем проекты пользователя
            const userResponse = await dataProvider.getOne('users', { id: record.id });
            if (userResponse.data.projects) {
                setProjects(userResponse.data.projects);
            }
            refresh();
        } catch (error: any) {
            console.error('Failed to add project:', error);
            notify(error.message || 'Ошибка добавления проекта', { type: 'error' });
        }
    };

    // Обновить роль в проекте
    const handleUpdateRole = async (projectId: number, newRole: 'admin' | 'manager' | 'client') => {
        if (!record?.id) return;

        try {
            const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
            const token = sessionStorage.getItem('auth_token');

            const response = await fetch(`${apiUrl}/projects/${projectId}/users/${record.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    role: newRole,
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Ошибка обновления роли');
            }

            notify('Роль успешно обновлена', { type: 'success' });
            setEditingProjectId(null);
            
            // Перезагружаем проекты пользователя
            const userResponse = await dataProvider.getOne('users', { id: record.id });
            if (userResponse.data.projects) {
                setProjects(userResponse.data.projects);
            }
            refresh();
        } catch (error: any) {
            console.error('Failed to update role:', error);
            notify(error.message || 'Ошибка обновления роли', { type: 'error' });
        }
    };

    // Удалить проект
    const handleRemoveProject = async (projectId: number) => {
        if (!record?.id) return;
        if (!window.confirm('Вы уверены, что хотите удалить доступ пользователя к этому проекту?')) {
            return;
        }

        try {
            const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
            const token = sessionStorage.getItem('auth_token');

            const response = await fetch(`${apiUrl}/projects/${projectId}/users/${record.id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Ошибка удаления проекта');
            }

            notify('Проект успешно удален', { type: 'success' });
            
            // Перезагружаем проекты пользователя
            const userResponse = await dataProvider.getOne('users', { id: record.id });
            if (userResponse.data.projects) {
                setProjects(userResponse.data.projects);
            }
            refresh();
        } catch (error: any) {
            console.error('Failed to remove project:', error);
            notify(error.message || 'Ошибка удаления проекта', { type: 'error' });
        }
    };

    // Получаем доступные проекты (те, которые еще не назначены)
    const availableProjects = allProjects.filter(
        p => !projects.some(up => up.project_id === p.id)
    );

    if (loading) {
        return <div>Загрузка...</div>;
    }

    if (!record?.id) {
        return null;
    }

    return (
        <div style={{ marginTop: '20px' }}>
            <h3>Проекты и роли</h3>

            {/* Таблица текущих проектов */}
            {projects.length > 0 && (
                <div style={{ marginTop: '15px', marginBottom: '20px', border: '1px solid #ddd', borderRadius: '4px' }}>
                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                        <thead>
                            <tr style={{ backgroundColor: '#f5f5f5' }}>
                                <th style={{ padding: '10px', textAlign: 'left', borderBottom: '1px solid #ddd' }}>Проект</th>
                                <th style={{ padding: '10px', textAlign: 'left', borderBottom: '1px solid #ddd' }}>Роль</th>
                                <th style={{ padding: '10px', textAlign: 'left', borderBottom: '1px solid #ddd' }}>Действия</th>
                            </tr>
                        </thead>
                        <tbody>
                            {projects.map((userProject) => (
                                <tr key={userProject.project_id}>
                                    <td style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
                                        {userProject.project_name || `Проект ${userProject.project_id}`}
                                    </td>
                                    <td style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
                                        {editingProjectId === userProject.project_id ? (
                                            <select
                                                value={userProject.role}
                                                onChange={(e) => {
                                                    const newRole = e.target.value as 'admin' | 'manager' | 'client';
                                                    handleUpdateRole(userProject.project_id, newRole);
                                                }}
                                                style={{ padding: '5px', minWidth: '150px' }}
                                            >
                                                <option value="admin">Админ</option>
                                                <option value="manager">Менеджер</option>
                                                <option value="client">Клиент</option>
                                            </select>
                                        ) : (
                                            <span>
                                                {userProject.role === 'admin' && 'Админ'}
                                                {userProject.role === 'manager' && 'Менеджер'}
                                                {userProject.role === 'client' && 'Клиент'}
                                            </span>
                                        )}
                                    </td>
                                    <td style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
                                        <button
                                            onClick={() => {
                                                if (editingProjectId === userProject.project_id) {
                                                    setEditingProjectId(null);
                                                } else {
                                                    setEditingProjectId(userProject.project_id);
                                                }
                                            }}
                                            style={{ marginRight: '10px', padding: '5px 10px', cursor: 'pointer' }}
                                        >
                                            {editingProjectId === userProject.project_id ? 'Отмена' : 'Изменить'}
                                        </button>
                                        <button
                                            onClick={() => handleRemoveProject(userProject.project_id)}
                                            style={{ padding: '5px 10px', cursor: 'pointer', color: 'red' }}
                                        >
                                            Удалить
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Добавить новый проект */}
            <div style={{ marginTop: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '4px' }}>
                <h4 style={{ marginTop: 0 }}>Добавить проект</h4>
                <div style={{ display: 'flex', gap: '10px', alignItems: 'flex-start', marginTop: '10px' }}>
                    <select
                        value={selectedProjectId || ''}
                        onChange={(e) => setSelectedProjectId(e.target.value ? Number(e.target.value) : null)}
                        style={{ padding: '8px', minWidth: '200px' }}
                    >
                        <option value="">Выберите проект</option>
                        {availableProjects.map(p => (
                            <option key={p.id} value={p.id}>{p.name}</option>
                        ))}
                    </select>
                    <select
                        value={selectedRole}
                        onChange={(e) => setSelectedRole(e.target.value as 'admin' | 'manager' | 'client')}
                        style={{ padding: '8px', minWidth: '150px' }}
                    >
                        <option value="admin">Админ</option>
                        <option value="manager">Менеджер</option>
                        <option value="client">Клиент</option>
                    </select>
                    <button
                        onClick={handleAddProject}
                        disabled={!selectedProjectId}
                        style={{ 
                            padding: '8px 16px', 
                            cursor: selectedProjectId ? 'pointer' : 'not-allowed',
                            backgroundColor: selectedProjectId ? '#1976d2' : '#ccc',
                            color: 'white',
                            border: 'none',
                            borderRadius: '4px'
                        }}
                    >
                        Добавить
                    </button>
                </div>
            </div>

            {projects.length === 0 && (
                <p style={{ marginTop: '15px', color: '#666' }}>
                    У пользователя пока нет назначенных проектов
                </p>
            )}
        </div>
    );
};

