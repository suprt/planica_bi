/**
 * Metrics Page
 * 
 * Страница с метриками Яндекс.Метрики проекта.
 * Показывает таблицу с показателями (визиты, пользователи, отказы, время на сайте, конверсии) и динамику.
 */

import React, { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { metricsService } from '../../services/api/metricsService';
import type { MetricsWithData } from '../../services/api/metricsService';
import { projectStorage } from '../../utils/projectStorage';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorMessage from '../../components/ErrorMessage';
import './Metrics.css';

interface MetricRow {
    name: string;
    period1: string | number;
    period2: string | number;
    period3: string | number;
    dynamics: number;
}

const Metrics: React.FC = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate();
    
    const [metricsData, setMetricsData] = useState<MetricsWithData | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');

    /**
     * Загрузка метрик при монтировании или изменении projectId
     */
    useEffect(() => {
        const fetchMetrics = async () => {
            // Получаем projectId из URL
            let currentProjectId = searchParams.get('project');
            
            // Если projectId нет в URL, берем из localStorage и обновляем URL
            if (!currentProjectId) {
                const lastProject = projectStorage.getLastProject();
                if (lastProject) {
                    currentProjectId = lastProject.toString();
                    // Обновляем URL без перезагрузки страницы
                    setSearchParams({ project: currentProjectId }, { replace: true });
                    // Ждем обновления URL перед загрузкой данных
                    return;
                }
            }
            
            if (!currentProjectId) {
                setError('Project ID не указан. Выберите проект.');
                setLoading(false);
                setTimeout(() => {
                    navigate('/dashboard/projects');
                }, 2000);
                return;
            }

            const parsedId = parseInt(currentProjectId, 10);
            if (isNaN(parsedId) || parsedId <= 0) {
                setError(`Некорректный Project ID: ${currentProjectId}`);
                setLoading(false);
                return;
            }

            try {
                setLoading(true);
                setError('');
                
                const data = await metricsService.getMetrics(parsedId);
                setMetricsData(data);
            } catch (err: any) {
                const errorMessage = err.response?.data?.error || 
                                   err.response?.data?.message || 
                                   'Не удалось загрузить метрики';
                setError(errorMessage);
            } finally {
                setLoading(false);
            }
        };

        fetchMetrics();
    }, [searchParams, navigate, setSearchParams]);

    /**
     * Расчет динамики между двумя значениями
     */
    const calculateDynamics = (current: number, previous: number): number => {
        if (!previous || previous === 0) return 0;
        return ((current - previous) / previous) * 100;
    };

    /**
     * Преобразование данных в формат для таблицы
     */
    const getMetrics = (): MetricRow[] => {
        if (!metricsData || !metricsData.rows || metricsData.rows.length === 0) {
            return [];
        }

        // Сортируем rows по месяцам (от нового к старому)
        const sortedRows = [...metricsData.rows].sort((a, b) => {
            return b.month.localeCompare(a.month);
        });

        const [r0, r1, r2] = sortedRows;

        const formatTime = (seconds: number): string => {
            const mins = Math.floor(seconds / 60);
            const secs = seconds % 60;
            return `${mins}:${secs.toString().padStart(2, '0')}`;
        };

        // Рассчитываем динамику для каждого показателя
        const visitsDynamics = r1 ? calculateDynamics(r0.visits, r1.visits) : 0;
        const usersDynamics = r1 ? calculateDynamics(r0.users, r1.users) : 0;
        const bounceDynamics = r1 ? calculateDynamics(r0.bounce_rate, r1.bounce_rate) : 0;
        const avgSecDynamics = r1 ? calculateDynamics(r0.avg_sec, r1.avg_sec) : 0;
        
        // Для конверсий
        const conv0 = r0.conversions || 0;
        const conv1 = r1?.conversions || 0;
        const convDynamics = r1 && conv1 > 0 ? calculateDynamics(conv0, conv1) : 0;

        return [
            { 
                name: 'Визиты', 
                period1: r0?.visits || 0, 
                period2: r1?.visits || 0, 
                period3: r2?.visits || 0, 
                dynamics: visitsDynamics
            },
            { 
                name: 'Пользователи', 
                period1: r0?.users || 0, 
                period2: r1?.users || 0, 
                period3: r2?.users || 0, 
                dynamics: usersDynamics
            },
            { 
                name: 'Отказы, %', 
                period1: `${(r0?.bounce_rate || 0).toFixed(2)}%`, 
                period2: `${(r1?.bounce_rate || 0).toFixed(2)}%`, 
                period3: `${(r2?.bounce_rate || 0).toFixed(2)}%`, 
                dynamics: bounceDynamics
            },
            { 
                name: 'Время на сайте', 
                period1: formatTime(r0?.avg_sec || 0), 
                period2: formatTime(r1?.avg_sec || 0), 
                period3: formatTime(r2?.avg_sec || 0), 
                dynamics: avgSecDynamics
            },
            { 
                name: 'Конверсии', 
                period1: r0?.conversions || 0, 
                period2: r1?.conversions || 0, 
                period3: r2?.conversions || 0, 
                dynamics: convDynamics
            },
        ];
    };

    const metrics = getMetrics();

    const formatDynamics = (value: number): string => {
        if (value === 0) return '0.00%';
        const sign = value > 0 ? '+' : '';
        return `${sign}${value.toFixed(2)}%`;
    };

    const getDynamicsClass = (value: number, metricName: string): string => {
        if (value === 0) return 'dynamics-neutral';
        // Для визитов, пользователей, времени на сайте, конверсий - рост это хорошо
        // Для отказов - рост это плохо
        const isPositiveMetric = ['Визиты', 'Пользователи', 'Время на сайте', 'Конверсии'].some(m => metricName.includes(m));
        if (isPositiveMetric) {
            return value > 0 ? 'dynamics-positive' : 'dynamics-negative';
        } else { // Для отказов - снижение это хорошо
            return value > 0 ? 'dynamics-negative' : 'dynamics-positive';
        }
    };

    const getDynamicsIcon = (value: number): string => {
        if (value === 0) return '';
        return value > 0 ? '↑' : '↓';
    };

    const handleRetry = () => {
        setError('');
        window.location.reload();
    };

    if (loading) {
        return <LoadingSpinner message="Загрузка метрик..." size="large" />;
    }

    if (error) {
        return (
            <ErrorMessage 
                title="Ошибка загрузки метрик"
                message={error}
                onRetry={handleRetry}
                type="error"
                fullPage
            />
        );
    }

    if (!metricsData || metrics.length === 0) {
        return (
            <ErrorMessage 
                title="Нет данных"
                message="Метрики Яндекс.Метрики для этого проекта не найдены или еще не синхронизированы."
                type="info"
                fullPage
            />
        );
    }

    return (
        <div className="metrics-page">
            <div className="metrics-header">
                <h1 className="page-title">Метрики Яндекс.Метрики</h1>
            </div>

            <div className="metrics-table-container">
                <table className="metrics-table">
                    <thead>
                        <tr>
                            <th className="metric-name-col">Показатель</th>
                            <th>Текущий период</th>
                            <th>Предыдущий период</th>
                            <th>Период M-2</th>
                            <th className="dynamics-col">Динамика, %</th>
                        </tr>
                    </thead>
                    <tbody>
                        {metrics.map((metric, index) => (
                            <tr key={index}>
                                <td className="metric-name">{metric.name}</td>
                                <td>{metric.period1}</td>
                                <td>{metric.period2}</td>
                                <td>{metric.period3}</td>
                                <td className={`dynamics ${getDynamicsClass(metric.dynamics, metric.name)}`}>
                                    <span className="dynamics-icon">{getDynamicsIcon(metric.dynamics)}</span>
                                    {formatDynamics(metric.dynamics)}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Metrics;

