/**
 * Statistics Page
 * 
 * Страница со статистикой по кампаниям Яндекс.Директ проекта.
 * Показывает таблицу с показателями кампании (показы, клики, CTR, CPC, конверсии, CPA, стоимость) и динамику.
 */

import React, { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { campaignsService } from '../../services/api/campaignsService';
import type { CampaignWithMetrics } from '../../services/api/campaignsService';
import { projectStorage } from '../../utils/projectStorage';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorMessage from '../../components/ErrorMessage';
import './Statistics.css';

interface CampaignMetricRow {
    name: string;
    period1: string | number;
    period2: string | number;
    period3: string | number;
    dynamics: number;
}

const Statistics: React.FC = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate();
    
    const [campaigns, setCampaigns] = useState<CampaignWithMetrics[]>([]);
    const [selectedCampaignId, setSelectedCampaignId] = useState<number | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');

    /**
     * Загрузка кампаний при монтировании или изменении projectId
     */
    useEffect(() => {
        const fetchCampaigns = async () => {
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
                
                const campaignsData = await campaignsService.getCampaigns(parsedId);
                setCampaigns(campaignsData);
                
                // Автоматически выбираем первую кампанию, если есть
                if (campaignsData && campaignsData.length > 0) {
                    setSelectedCampaignId(campaignsData[0].campaignId);
                }
            } catch (err: any) {
                const errorMessage = err.response?.data?.error || 
                                   err.response?.data?.message || 
                                   'Не удалось загрузить кампании';
                setError(errorMessage);
            } finally {
                setLoading(false);
            }
        };

        fetchCampaigns();
    }, [searchParams, navigate, setSearchParams]);

    /**
     * Получение выбранной кампании
     */
    const getSelectedCampaign = (): CampaignWithMetrics | null => {
        if (!campaigns || !selectedCampaignId) return null;
        
        const campaign = campaigns.find(c => c.campaignId === selectedCampaignId);
        return campaign || null;
    };

    /**
     * Расчет динамики между двумя значениями
     */
    const calculateDynamics = (current: number, previous: number): number => {
        if (!previous || previous === 0) return 0;
        return ((current - previous) / previous) * 100;
    };

    /**
     * Преобразование данных кампании в формат для таблицы
     */
    const getCampaignMetrics = (): CampaignMetricRow[] => {
        const campaign = getSelectedCampaign();
        if (!campaign || !campaign.rows || campaign.rows.length === 0) {
            return [];
        }

        // Сортируем rows по месяцам (от нового к старому)
        const sortedRows = [...campaign.rows].sort((a, b) => {
            return b.month.localeCompare(a.month);
        });

        const [r0, r1, r2] = sortedRows;

        // Рассчитываем динамику для каждого показателя
        const impressionsDynamics = r1 ? calculateDynamics(r0.impressions, r1.impressions) : 0;
        const clicksDynamics = r1 ? calculateDynamics(r0.clicks, r1.clicks) : 0;
        const ctrDynamics = r1 ? calculateDynamics(r0.ctr, r1.ctr) : 0;
        const cpcDynamics = r1 ? calculateDynamics(r0.cpc, r1.cpc) : 0;
        const costDynamics = r1 ? calculateDynamics(r0.cost, r1.cost) : 0;
        
        // Для конверсий и CPA
        const conv0 = r0.conv || 0;
        const conv1 = r1?.conv || 0;
        const convDynamics = r1 && conv1 > 0 ? calculateDynamics(conv0, conv1) : 0;
        
        const cpa0 = r0.cpa || 0;
        const cpa1 = r1?.cpa || 0;
        const cpaDynamics = r1 && cpa1 > 0 ? calculateDynamics(cpa0, cpa1) : 0;

        return [
            { 
                name: 'Показы', 
                period1: r0?.impressions || 0, 
                period2: r1?.impressions || 0, 
                period3: r2?.impressions || 0, 
                dynamics: impressionsDynamics
            },
            { 
                name: 'Клики', 
                period1: r0?.clicks || 0, 
                period2: r1?.clicks || 0, 
                period3: r2?.clicks || 0, 
                dynamics: clicksDynamics
            },
            { 
                name: 'CTR, %', 
                period1: `${(r0?.ctr || 0).toFixed(2)}%`, 
                period2: `${(r1?.ctr || 0).toFixed(2)}%`, 
                period3: `${(r2?.ctr || 0).toFixed(2)}%`, 
                dynamics: ctrDynamics
            },
            { 
                name: 'CPC, ₽', 
                period1: `${(r0?.cpc || 0).toFixed(2)}`, 
                period2: `${(r1?.cpc || 0).toFixed(2)}`, 
                period3: `${(r2?.cpc || 0).toFixed(2)}`, 
                dynamics: cpcDynamics
            },
            { 
                name: 'Конверсии', 
                period1: r0?.conv || 0, 
                period2: r1?.conv || 0, 
                period3: r2?.conv || 0, 
                dynamics: convDynamics
            },
            { 
                name: 'CPA, ₽', 
                period1: r0?.cpa ? `${r0.cpa.toFixed(2)}` : '—', 
                period2: r1?.cpa ? `${r1.cpa.toFixed(2)}` : '—', 
                period3: r2?.cpa ? `${r2.cpa.toFixed(2)}` : '—', 
                dynamics: cpaDynamics
            },
            { 
                name: 'Стоимость, ₽', 
                period1: `${(r0?.cost || 0).toFixed(2)}`, 
                period2: `${(r1?.cost || 0).toFixed(2)}`, 
                period3: `${(r2?.cost || 0).toFixed(2)}`, 
                dynamics: costDynamics
            },
        ];
    };

    const metrics = getCampaignMetrics();
    const selectedCampaign = getSelectedCampaign();

    const formatDynamics = (value: number): string => {
        if (value === 0) return '0.00%';
        const sign = value > 0 ? '+' : '';
        return `${sign}${value.toFixed(2)}%`;
    };

    const getDynamicsClass = (value: number, metricName: string): string => {
        if (value === 0) return 'dynamics-neutral';
        // Для CTR, конверсий - рост это хорошо
        // Для CPC, CPA, стоимости - рост это плохо
        const isPositiveMetric = ['CTR', 'Конверсии'].some(m => metricName.includes(m));
        if (isPositiveMetric) {
            return value > 0 ? 'dynamics-positive' : 'dynamics-negative';
        } else { // Для CPC, CPA, стоимости - снижение это хорошо
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
        return <LoadingSpinner message="Загрузка статистики..." size="large" />;
    }

    if (error) {
        if (error === 'GENERATING') {
            return (
                <ErrorMessage 
                    title="Генерация отчета"
                    message="Отчет генерируется. Пожалуйста, подождите несколько секунд и обновите страницу."
                    onRetry={handleRetry}
                    type="info"
                    fullPage
                />
            );
        }
        return (
            <ErrorMessage 
                title="Ошибка загрузки статистики"
                message={error}
                onRetry={handleRetry}
                type="error"
                fullPage
            />
        );
    }

    if (!campaigns || campaigns.length === 0) {
        return (
            <ErrorMessage 
                title="Нет данных"
                message="Кампании Яндекс.Директ для этого проекта не найдены или еще не синхронизированы."
                type="info"
                fullPage
            />
        );
    }

    if (!selectedCampaign || metrics.length === 0) {
        return (
            <ErrorMessage 
                title="Нет данных"
                message="Статистика по выбранной кампании пока не сформирована или недоступна."
                type="info"
                fullPage
            />
        );
    }

    return (
        <div className="statistics-page">
            <div className="statistics-header">
                <h1 className="page-title">Статистика кампаний</h1>
                <div className="campaign-selector">
                    <label htmlFor="campaign-select">Выберите кампанию:</label>
                    <select
                        id="campaign-select"
                        value={selectedCampaignId || ''}
                        onChange={(e) => setSelectedCampaignId(Number(e.target.value))}
                        className="campaign-dropdown"
                    >
                        {campaigns.map((campaign) => (
                            <option key={campaign.campaignId} value={campaign.campaignId}>
                                {campaign.name} (ID: {campaign.campaignId})
                            </option>
                        ))}
                    </select>
                </div>
            </div>

            <div className="statistics-table-container">
                <table className="statistics-table">
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

export default Statistics;
