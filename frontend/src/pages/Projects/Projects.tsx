/**
 * Projects Page
 * 
 * Страница со списком всех проектов пользователя.
 * Отображает карточки проектов с возможностью перехода к отчету.
 */

import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { projectsService, Project } from '../../services/api/projectsService';
import { syncService } from '../../services/api/syncService';
import { useAuth } from '../../contexts/AuthContext';
import { projectStorage } from '../../utils/projectStorage';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorMessage from '../../components/ErrorMessage';
import './Projects.css';

const Projects: React.FC = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');
    const [syncingProjects, setSyncingProjects] = useState<Set<number>>(new Set());
    
    // Проверяем, является ли пользователь админом
    const isAdmin = user?.project_roles?.some((role: any) => role.role === 'admin') || false;
    
    // Логируем статус админа для отладки
    useEffect(() => {
        console.log('[Projects] User:', user);
        console.log('[Projects] User project_roles:', user?.project_roles);
        console.log('[Projects] Is admin:', isAdmin);
    }, [user, isAdmin]);

    /**
     * Загрузка списка проектов при монтировании
     */
    useEffect(() => {
        const fetchProjects = async () => {
            try {
                setLoading(true);
                setError('');
                
                const data = await projectsService.getAll();
                setProjects(data || []);
                
                console.log('[Projects] Loaded projects:', data?.length || 0);
            } catch (err: any) {
                const errorMessage = err.response?.data?.error || 
                                   err.response?.data?.message || 
                                   'Не удалось загрузить проекты';
                setError(errorMessage);
                console.error('[Projects] Error loading projects:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchProjects();
    }, []);

    /**
     * Переход к отчету проекта
     */
    const handleProjectClick = (projectId: number | undefined | null) => {
        console.log('🎯 [Projects] handleProjectClick called');
        console.log('🎯 [Projects] projectId:', projectId);
        console.log('🎯 [Projects] projectId type:', typeof projectId);
        
        if (!projectId || isNaN(Number(projectId))) {
            console.error('❌ [Projects] Invalid project ID:', projectId);
            return;
        }
        
        // Сохраняем последний выбранный проект
        projectStorage.setLastProject(projectId);
        navigate(`/dashboard/reports?project=${projectId}`);
    };

    /**
     * Повторная попытка загрузки
     */
    const handleRetry = () => {
        setError('');
        setLoading(true);
        // Перезагружаем компонент (useEffect сработает автоматически)
        window.location.reload();
    };

    /**
     * Синхронизация проекта (только для админов)
     */
    const handleSyncProject = async (e: React.MouseEvent, projectId: number | undefined | null) => {
        e.stopPropagation();
        
        if (!projectId || isNaN(Number(projectId))) {
            console.error('❌ [Projects] Invalid project ID for sync:', projectId);
            return;
        }

        if (!isAdmin) {
            console.warn('⚠️ [Projects] Only admins can sync projects');
            return;
        }

        try {
            setSyncingProjects(prev => new Set(prev).add(projectId));
            await syncService.syncProject(projectId);
            console.log('✅ [Projects] Sync task enqueued for project:', projectId);
            // Можно показать уведомление об успехе
            alert(`Синхронизация проекта запущена. Задача поставлена в очередь.`);
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 
                               err.response?.data?.message || 
                               'Не удалось запустить синхронизацию';
            console.error('❌ [Projects] Failed to sync project:', err);
            alert(`Ошибка синхронизации: ${errorMessage}`);
        } finally {
            setSyncingProjects(prev => {
                const newSet = new Set(prev);
                newSet.delete(projectId);
                return newSet;
            });
        }
    };

    // Состояние загрузки
    if (loading) {
        return <LoadingSpinner message="Загрузка проектов..." size="large" />;
    }

    // Состояние ошибки
    if (error) {
        return (
            <ErrorMessage 
                title="Ошибка загрузки проектов"
                message={error}
                onRetry={handleRetry}
                type="error"
                fullPage
            />
        );
    }

    // Пустой список проектов
    if (projects.length === 0) {
        return (
            <div className="projects-empty">
                <div className="empty-icon">📊</div>
                <h2>У вас пока нет проектов</h2>
                <p>Создайте первый проект для начала работы с аналитикой</p>
                <button className="create-project-button">
                    Создать проект
                </button>
            </div>
        );
    }

    // Список проектов
    return (
        <div className="projects-page">
            <div className="projects-header">
                <h1>Мои проекты</h1>
                <p className="projects-count">
                    Всего проектов: <strong>{projects.length}</strong>
                </p>
            </div>

            <div className="projects-grid">
                {projects.map((project) => {
                    console.log('🔍 [Projects] Rendering project:', project);
                    console.log('🔍 [Projects] Project ID:', project.id, 'Type:', typeof project.id);
                    console.log('🔍 [Projects] Full project object:', JSON.stringify(project, null, 2));
                    return (
                        <div 
                            key={project.id || `project-${project.slug}`}
                            className={`project-card ${!project.is_active ? 'project-inactive' : ''}`}
                            onClick={() => {
                                console.log('👆 [Projects] Card clicked, project:', project);
                                console.log('👆 [Projects] Project ID:', project.id);
                                handleProjectClick(project.id);
                            }}
                        >
                        <div className="project-card-header">
                            <h3 className="project-name">{project.name}</h3>
                            <span className={`project-status ${project.is_active ? 'status-active' : 'status-inactive'}`}>
                                {project.is_active ? '●' : '○'}
                            </span>
                        </div>

                        <div className="project-card-body">
                            <div className="project-info">
                                <div className="info-item">
                                    <span className="info-label">Slug:</span>
                                    <span className="info-value">{project.slug}</span>
                                </div>
                                <div className="info-item">
                                    <span className="info-label">Валюта:</span>
                                    <span className="info-value">{project.currency}</span>
                                </div>
                                <div className="info-item">
                                    <span className="info-label">Часовой пояс:</span>
                                    <span className="info-value">{project.timezone}</span>
                                </div>
                            </div>
                        </div>

                        <div className="project-card-footer">
                            <span className="project-date">
                                Создан: {new Date(project.created_at).toLocaleDateString('ru-RU')}
                            </span>
                            <div className="project-card-actions">
                                {isAdmin && (
                                    <button 
                                        className="sync-button"
                                        onClick={(e) => handleSyncProject(e, project.id)}
                                        disabled={syncingProjects.has(project.id as number)}
                                        title="Синхронизировать данные проекта"
                                        style={{ display: 'flex' }}
                                    >
                                        <span style={{ fontSize: '18px' }}>
                                            {syncingProjects.has(project.id as number) ? '⏳' : '🔄'}
                                        </span>
                                    </button>
                                )}
                                {!isAdmin && (
                                    <span style={{ fontSize: '12px', color: '#999', marginRight: '8px' }}>
                                        Только для админов
                                    </span>
                                )}
                                <button 
                                    className="view-report-button"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        console.log('👆 [Projects] Button clicked, project:', project);
                                        handleProjectClick(project.id);
                                    }}
                                >
                                    Открыть отчет →
                                </button>
                            </div>
                        </div>
                    </div>
                    );
                })}
            </div>
        </div>
    );
};

export default Projects;

