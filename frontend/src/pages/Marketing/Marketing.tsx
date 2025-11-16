import React, { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import './Marketing.css';
import { MarketingData } from '../../types/marketing';
import { marketingService } from '../../services/api/marketingService';
import { projectStorage } from '../../utils/projectStorage';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorMessage from '../../components/ErrorMessage';

const Marketing: React.FC = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate();
    
    const [data, setData] = useState<MarketingData | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');

    // Загрузка данных через API
    useEffect(() => {
        const loadMarketingData = async () => {
            let currentProjectId = searchParams.get('project');
            
            if (!currentProjectId) {
                const lastProject = projectStorage.getLastProject();
                if (lastProject) {
                    currentProjectId = lastProject.toString();
                    setSearchParams({ project: currentProjectId }, { replace: true });
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
                
                const apiData = await marketingService.getMarketing(parsedId);
                setData(apiData);
            } catch (err: any) {
                const errorMessage = err.response?.data?.error || 
                                   err.response?.data?.message || 
                                   'Не удалось загрузить данные маркетинга';
                setError(errorMessage);
                console.error('Error loading marketing data:', err);
            } finally {
                setLoading(false);
            }
        };

        loadMarketingData();
    }, [searchParams, navigate, setSearchParams]);

    const formatValue = (value: number | string): string => {
        if (typeof value === 'number') {
            // Если значение уже содержит %, возвращаем как есть
            if (String(value).includes('%')) {
                return String(value);
            }
            // Форматируем числа с разделителями тысяч
            return value.toLocaleString('ru-RU', { 
                minimumFractionDigits: value % 1 === 0 ? 0 : 2, 
                maximumFractionDigits: 2 
            });
        }
        return value;
    };

    const formatEfficiency = (value: number): string => {
        if (value === 0) return '0.00%';
        const sign = value > 0 ? '+' : '';
        return `${sign}${value.toFixed(2)}%`;
    };

    // Определяем, является ли метрика "хорошей" (увеличение = хорошо)
    const isPositiveMetric = (indicator: string): boolean => {
        const positiveMetrics = [
            'Клики',
            'Клики MC',
            'CTR',
            'Конверсии',
            'Конверсии MC',
            'Конверсии RSA',
        ];
        return positiveMetrics.some(metric => indicator.includes(metric));
    };

    // Для CPA увеличение = плохо
    const isCostMetric = (indicator: string): boolean => {
        return indicator.includes('CPA');
    };

    const getEfficiencyClass = (value: number, indicator: string): string => {
        if (value === 0) return 'efficiency-neutral';
        
        if (isCostMetric(indicator)) {
            // Для CPA увеличение = плохо (красный), уменьшение = хорошо (зеленый)
            return value > 0 ? 'efficiency-negative' : 'efficiency-positive';
        }
        
        const isPositive = isPositiveMetric(indicator);
        if (isPositive) {
            return value > 0 ? 'efficiency-positive' : 'efficiency-negative';
        } else {
            return value > 0 ? 'efficiency-negative' : 'efficiency-positive';
        }
    };

    const getSummaryArrow = (change: number): string => {
        if (change === 0) return '';
        return change > 0 ? '↑' : '↓';
    };

    const getSummaryArrowClass = (change: number, isPositive: boolean): string => {
        if (change === 0) return 'summary-arrow-neutral';
        // Для CPA и других "плохих" метрик увеличение = красный
        if (!isPositive && change > 0) return 'summary-arrow-negative';
        if (isPositive && change > 0) return 'summary-arrow-positive';
        return 'summary-arrow-negative';
    };

    if (loading) {
        return <LoadingSpinner message="Загрузка данных маркетинга..." size="large" />;
    }

    if (error) {
        return (
            <ErrorMessage 
                title="Ошибка загрузки данных маркетинга"
                message={error}
                onRetry={() => window.location.reload()}
                type="error"
                fullPage
            />
        );
    }

    if (!data) {
        return (
            <ErrorMessage 
                title="Нет данных"
                message="Данные маркетинга для этого проекта не найдены или еще не синхронизированы."
                type="info"
                fullPage
            />
        );
    }

    return (
        <div className="marketing-page">
            <div className="marketing-header">
                <h1 className="page-title">Маркетинг</h1>
            </div>

            {/* Секция Клики */}
            <div className="marketing-section">
                <div className="section-header">
                    <h2 className="section-title">Клики</h2>
                </div>
                <div className="section-content">
                    <div className="summary-list">
                        {data.clicks.summary.map((item, index) => (
                            <div key={index} className="summary-item">
                                <span className="summary-label">{item.label}</span>
                                <div className={`summary-value ${getSummaryArrowClass(item.change, item.isPositive)}`}>
                                    <span className={`summary-arrow ${getSummaryArrowClass(item.change, item.isPositive)}`}>
                                        {getSummaryArrow(item.change)}
                                    </span>
                                    <span className="summary-text">{item.value}</span>
                                </div>
                            </div>
                        ))}
                    </div>
                    <div className="metrics-table-container">
                        <table className="metrics-table">
                            <thead>
                                <tr>
                                    <th>№ п/п</th>
                                    <th>Показатель</th>
                                    <th>Октябрь</th>
                                    <th>Сентябрь</th>
                                    <th>Август</th>
                                    <th>Эффективность, %</th>
                                </tr>
                            </thead>
                            <tbody>
                                {data.clicks.metrics.map((metric) => (
                                    <tr key={metric.id}>
                                        <td className="row-number">{metric.id}</td>
                                        <td className="indicator-cell">{metric.indicator}</td>
                                        <td>{formatValue(metric.october)}</td>
                                        <td>{formatValue(metric.september)}</td>
                                        <td>{formatValue(metric.august)}</td>
                                        <td className={`efficiency ${getEfficiencyClass(metric.efficiency, metric.indicator)}`}>
                                            {formatEfficiency(metric.efficiency)}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            {/* Секция Конверсии */}
            <div className="marketing-section">
                <div className="section-header">
                    <h2 className="section-title">Конверсии</h2>
                </div>
                <div className="section-content">
                    <div className="summary-list">
                        {data.conversions.summary.map((item, index) => (
                            <div key={index} className="summary-item">
                                <span className="summary-label">{item.label}</span>
                                <div className={`summary-value ${getSummaryArrowClass(item.change, item.isPositive)}`}>
                                    <span className={`summary-arrow ${getSummaryArrowClass(item.change, item.isPositive)}`}>
                                        {getSummaryArrow(item.change)}
                                    </span>
                                    <span className="summary-text">{item.value}</span>
                                </div>
                            </div>
                        ))}
                    </div>
                    <div className="metrics-table-container">
                        <table className="metrics-table">
                            <thead>
                                <tr>
                                    <th>№ п/п</th>
                                    <th>Показатель</th>
                                    <th>Октябрь</th>
                                    <th>Сентябрь</th>
                                    <th>Август</th>
                                    <th>Эффективность, %</th>
                                </tr>
                            </thead>
                            <tbody>
                                {data.conversions.metrics.map((metric) => (
                                    <tr key={metric.id}>
                                        <td className="row-number">{metric.id}</td>
                                        <td className="indicator-cell">{metric.indicator}</td>
                                        <td>{formatValue(metric.october)}</td>
                                        <td>{formatValue(metric.september)}</td>
                                        <td>{formatValue(metric.august)}</td>
                                        <td className={`efficiency ${getEfficiencyClass(metric.efficiency, metric.indicator)}`}>
                                            {formatEfficiency(metric.efficiency)}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Marketing;