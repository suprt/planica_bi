/**
 * Reports Page
 * 
 * Страница с детальным отчетом проекта.
 * Загружает полные данные из API (метрики, директ, SEO, AI-анализ).
 */

import React, { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { reportsService, Report, ReportStatus } from '../../services/api/reportsService';
import { projectsService } from '../../services/api/projectsService';
import { projectStorage } from '../../utils/projectStorage';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorMessage from '../../components/ErrorMessage';
import './Reports.css';

interface MetricRow {
    name: string;
    october: string | number;
    september: string | number;
    august: string | number;
    efficiency: number;
}

const Reports: React.FC = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate();
    
    const [report, setReport] = useState<Report | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');
    const [publicLink, setPublicLink] = useState<string>('');
    const [copySuccess, setCopySuccess] = useState<boolean>(false);

    /**
     * Загрузка отчета при монтировании или изменении projectId
     */
    useEffect(() => {
        const fetchReport = async () => {
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
            
            console.log('📡 [Reports] fetchReport called, projectId:', currentProjectId);
            
            if (!currentProjectId) {
                setError('Project ID не указан. Выберите проект.');
                setLoading(false);
                // Если нет проекта, перенаправляем на страницу проектов
                setTimeout(() => {
                    navigate('/dashboard/projects');
                }, 2000);
                return;
            }

            // Парсим ID и проверяем валидность
            const parsedId = parseInt(currentProjectId, 10);
            if (isNaN(parsedId) || parsedId <= 0) {
                setError(`Некорректный Project ID: ${currentProjectId}`);
                setLoading(false);
                console.error('❌ [Reports] Invalid project ID:', currentProjectId, 'parsed:', parsedId);
                return;
            }

            try {
                setLoading(true);
                setError('');
                
                console.log('🚀 [Reports] Fetching report for project ID:', parsedId);
                const data = await reportsService.getReport(parsedId);
                
                // Проверяем, это статус задачи или сам отчет
                if ((data as any).status) {
                    const statusData = data as ReportStatus;
                    console.log('⏳ [Reports] Report is being generated, status:', statusData.status);
                    
                    if (statusData.status === 'pending' || statusData.status === 'processing') {
                        // Устанавливаем специальный флаг для отображения информационного сообщения
                        setError('GENERATING'); // Специальное значение для информационного сообщения
                        setReport(null);
                    } else if (statusData.status === 'failed') {
                        setError(statusData.message || 'Ошибка генерации отчета');
                        setReport(null);
                    } else {
                        // completed, но почему-то вернулся статус вместо отчета
                        setError('Отчет готов, но данные не получены. Попробуйте обновить страницу.');
                        setReport(null);
                    }
                } else {
                    // Это полноценный отчет
                    setReport(data as Report);
                    console.log('✅ [Reports] Report loaded successfully');
                    console.log('📊 [Reports] Direct totals:', (data as Report).direct?.totals?.length || 0);
                    console.log('📊 [Reports] Direct data:', (data as Report).direct);
                }
            } catch (err: any) {
                const errorMessage = err.response?.data?.error || 
                                   err.response?.data?.message || 
                                   'Не удалось загрузить отчет';
                setError(errorMessage);
                console.error('❌ [Reports] Error loading report:', err);
                console.error('❌ [Reports] Response data:', err.response?.data);
            } finally {
                setLoading(false);
            }
        };

        fetchReport();
    }, [searchParams, navigate, setSearchParams]);

    /**
     * Преобразование данных из API в формат для таблицы
     */
    const getMetrics = (): MetricRow[] => {
        if (!report || !report.metrica.summary.length) {
            return [];
        }

        const [m0, m1, m2] = report.metrica.summary;

        const formatTime = (seconds: number): string => {
            const mins = Math.floor(seconds / 60);
            const secs = seconds % 60;
            return `${mins}:${secs.toString().padStart(2, '0')}`;
        };

        // Используем динамику из backend
        return [
            { 
                name: 'Посетители, кол-во', 
                october: m0?.users || 0, 
                september: m1?.users || 0, 
                august: m2?.users || 0, 
                efficiency: m0?.dynamics?.users || 0 
            },
            { 
                name: 'Визиты, кол-во', 
                october: m0?.visits || 0, 
                september: m1?.visits || 0, 
                august: m2?.visits || 0, 
                efficiency: m0?.dynamics?.visits || 0 
            },
            { 
                name: 'Кол-во отказов, %', 
                october: `${(m0?.bounce || 0).toFixed(2)}%`, 
                september: `${(m1?.bounce || 0).toFixed(2)}%`, 
                august: `${(m2?.bounce || 0).toFixed(2)}%`, 
                efficiency: m0?.dynamics?.bounce || 0 
            },
            { 
                name: 'Время на сайте, сек', 
                october: formatTime(m0?.avgSec || 0), 
                september: formatTime(m1?.avgSec || 0), 
                august: formatTime(m2?.avgSec || 0), 
                efficiency: m0?.dynamics?.avgSec || 0 
            },
            { 
                name: 'Конверсия, %', 
                october: `${(m0?.conv || 0).toFixed(2)}%`, 
                september: `${(m1?.conv || 0).toFixed(2)}%`, 
                august: `${(m2?.conv || 0).toFixed(2)}%`, 
                efficiency: m0?.dynamics?.conv || 0 
            },
        ];
    };

    const getDirectMetrics = (): MetricRow[] => {
        if (!report) {
            console.log('⚠️ [Reports] No report data');
            return [];
        }
        
        if (!report.direct || !report.direct.totals || report.direct.totals.length === 0) {
            console.log('⚠️ [Reports] No Direct totals data', {
                hasDirect: !!report.direct,
                totalsLength: report.direct?.totals?.length || 0,
                direct: report.direct
            });
            return [];
        }

        const [d0, d1, d2] = report.direct.totals;
        console.log('📊 [Reports] Direct metrics data:', { d0, d1, d2 });

        // Вычисляем динамику для Direct метрик
        const calculateDynamics = (current: number, previous: number): number => {
            if (previous === 0) return 0;
            return ((current - previous) / previous) * 100;
        };

        return [
            { 
                name: 'Показы, кол-во', 
                october: d0?.impressions || 0, 
                september: d1?.impressions || 0, 
                august: d2?.impressions || 0, 
                efficiency: d1 ? calculateDynamics(d0?.impressions || 0, d1.impressions) : 0
            },
            { 
                name: 'Клики, кол-во', 
                october: d0?.clicks || 0, 
                september: d1?.clicks || 0, 
                august: d2?.clicks || 0, 
                efficiency: d1 ? calculateDynamics(d0?.clicks || 0, d1.clicks) : 0
            },
            { 
                name: 'CTR, %', 
                october: `${(d0?.ctr || 0).toFixed(2)}%`, 
                september: `${(d1?.ctr || 0).toFixed(2)}%`, 
                august: `${(d2?.ctr || 0).toFixed(2)}%`, 
                efficiency: d1 ? calculateDynamics(d0?.ctr || 0, d1.ctr) : 0
            },
            { 
                name: 'CPC, руб.', 
                october: `${(d0?.cpc || 0).toFixed(2)}`, 
                september: `${(d1?.cpc || 0).toFixed(2)}`, 
                august: `${(d2?.cpc || 0).toFixed(2)}`, 
                efficiency: d1 ? calculateDynamics(d0?.cpc || 0, d1.cpc) : 0
            },
            { 
                name: 'CPA, руб.', 
                october: d0?.cpa ? `${d0.cpa.toFixed(2)}` : '-', 
                september: d1?.cpa ? `${d1.cpa.toFixed(2)}` : '-', 
                august: d2?.cpa ? `${d2.cpa.toFixed(2)}` : '-', 
                efficiency: (d1?.cpa && d0?.cpa) ? calculateDynamics(d0.cpa, d1.cpa) : 0
            },
            { 
                name: 'Расходы, руб.', 
                october: `${(d0?.cost || 0).toFixed(2)}`, 
                september: `${(d1?.cost || 0).toFixed(2)}`, 
                august: `${(d2?.cost || 0).toFixed(2)}`, 
                efficiency: d1 ? calculateDynamics(d0?.cost || 0, d1.cost) : 0
            },
            { 
                name: 'Конверсии, кол-во', 
                october: d0?.conv || '-', 
                september: d1?.conv || '-', 
                august: d2?.conv || '-', 
                efficiency: (d1?.conv !== undefined && d0?.conv !== undefined) ? calculateDynamics(d0.conv, d1.conv) : 0
            },
        ];
    };

    const metrics = getMetrics();
    const directMetrics = getDirectMetrics();

    const formatEfficiency = (value: number): string => {
        if (value === 0) return '0.00%';
        const sign = value > 0 ? '+' : '';
        return `${sign}${value.toFixed(2)}%`;
    };

    const isPositiveMetric = (metricName: string): boolean => {
        const positiveMetrics = [
            'Время на сайте',
            'Конверсия',
            'Посетители',
            'Визиты',
            'Клики',
            'Показы',
            'CTR',
            'Конверсии'
        ];
        const negativeMetrics = [
            'Кол-во отказов',
            'CPC',
            'CPA',
            'Расходы'
        ];
        // Для Direct метрик: CPC, CPA, Расходы - меньше лучше (отрицательная динамика = хорошо)
        if (negativeMetrics.some(metric => metricName.includes(metric))) {
            return false;
        }
        return positiveMetrics.some(metric => metricName.includes(metric));
    };

    const getEfficiencyClass = (value: number, metricName: string): string => {
        if (value === 0) return 'efficiency-neutral';
        const isPositive = isPositiveMetric(metricName);
        if (isPositive) {
            return value > 0 ? 'efficiency-positive' : 'efficiency-negative';
        } else {
            return value > 0 ? 'efficiency-negative' : 'efficiency-positive';
        }
    };

    const getEfficiencyIcon = (value: number): string => {
        if (value === 0) return '';
        return value > 0 ? '↑' : '↓';
    };

    const getSummaryData = () => {
        if (!report || !report.metrica.summary.length) {
            return { traffic: 0, conversions: 0, bounce: 0 };
        }

        const [m0] = report.metrica.summary;
        
        // Используем динамику из backend
        return {
            traffic: m0?.dynamics?.visits || 0,
            conversions: m0?.dynamics?.conv || 0,
            bounce: m0?.dynamics?.bounce || 0,
        };
    };

    const summaryData = getSummaryData();

    const handleRetry = () => {
        setError('');
        setLoading(true);
        window.location.reload();
    };

    /**
     * Получить публичную ссылку на отчет
     */
    const handleGetPublicLink = async () => {
        const projectId = searchParams.get('project') || projectStorage.getLastProject();
        if (!projectId) {
            alert('Не выбран проект');
            return;
        }

        try {
            const parsedId = parseInt(projectId.toString(), 10);
            const linkData = await projectsService.getPublicLink(parsedId);
            setPublicLink(linkData.public_url);
        } catch (err: any) {
            console.error('Failed to get public link:', err);
            alert('Не удалось получить публичную ссылку: ' + (err.response?.data?.error || err.message));
        }
    };

    /**
     * Копировать публичную ссылку в буфер обмена
     */
    const handleCopyLink = async () => {
        if (!publicLink) {
            await handleGetPublicLink();
            return;
        }

        try {
            await navigator.clipboard.writeText(publicLink);
            setCopySuccess(true);
            setTimeout(() => setCopySuccess(false), 2000);
        } catch (err) {
            console.error('Failed to copy link:', err);
            // Fallback: select text
            const textArea = document.createElement('textarea');
            textArea.value = publicLink;
            document.body.appendChild(textArea);
            textArea.select();
            try {
                document.execCommand('copy');
                setCopySuccess(true);
                setTimeout(() => setCopySuccess(false), 2000);
            } catch (fallbackErr) {
                alert('Не удалось скопировать ссылку. Скопируйте вручную: ' + publicLink);
            }
            document.body.removeChild(textArea);
        }
    };

    // Создаем компонент заголовка с кнопкой, который показывается всегда
    const renderHeader = () => (
        <div className="reports-header">
            <div>
                <h1 className="page-title">Отчет по проекту</h1>
                {report && (
                    <p className="project-periods">
                        Периоды: {report.periods.join(', ')}
                    </p>
                )}
            </div>
            <div className="reports-header-actions">
                <button
                    onClick={handleGetPublicLink}
                    className="public-link-button"
                    style={{
                        padding: '8px 16px',
                        marginRight: '10px',
                        backgroundColor: '#1976d2',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer'
                    }}
                >
                    Получить публичную ссылку
                </button>
                {publicLink && (
                    <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        <input
                            type="text"
                            value={publicLink}
                            readOnly
                            style={{
                                padding: '8px',
                                minWidth: '400px',
                                border: '1px solid #ddd',
                                borderRadius: '4px',
                                fontSize: '14px'
                            }}
                        />
                        <button
                            onClick={handleCopyLink}
                            className="copy-link-button"
                            style={{
                                padding: '8px 16px',
                                backgroundColor: copySuccess ? '#4caf50' : '#2196f3',
                                color: 'white',
                                border: 'none',
                                borderRadius: '4px',
                                cursor: 'pointer'
                            }}
                        >
                            {copySuccess ? '✓ Скопировано!' : '📋 Копировать'}
                        </button>
                    </div>
                )}
            </div>
        </div>
    );

    if (loading) {
        return (
            <div className="reports-page">
                {renderHeader()}
                <LoadingSpinner message="Загрузка отчета..." size="large" />
            </div>
        );
    }

    if (error) {
        // Специальная обработка для статуса генерации отчета
        if (error === 'GENERATING') {
            return (
                <div className="reports-page">
                    {renderHeader()}
                    <ErrorMessage 
                        title="Генерация отчета"
                        message="Отчет генерируется. Пожалуйста, подождите несколько секунд и обновите страницу."
                        onRetry={handleRetry}
                        type="info"
                        fullPage
                    />
                </div>
            );
        }
        
        return (
            <div className="reports-page">
                {renderHeader()}
                <ErrorMessage 
                    title="Ошибка загрузки отчета"
                    message={error}
                    onRetry={handleRetry}
                    type="error"
                    fullPage
                />
            </div>
        );
    }

    if (!report || metrics.length === 0) {
        return (
            <div className="reports-page">
                {renderHeader()}
                <ErrorMessage 
                    title="Нет данных"
                    message="Отчет по этому проекту пока не сформирован"
                    type="info"
                    fullPage
                />
            </div>
        );
    }

    return (
        <div className="reports-page">
            {renderHeader()}

            <div className="reports-content">
                <div className="reports-summary">
                    <div className="summary-card">
                        <div className="summary-label">Трафик</div>
                        <div className={`summary-value ${summaryData.traffic < 0 ? 'negative' : 'positive'}`}>
                            <span className="summary-arrow">{summaryData.traffic < 0 ? '↓' : '↑'}</span>
                            <span className="summary-text">
                                {summaryData.traffic < 0 ? 'Упал' : 'Вырос'} на {Math.abs(summaryData.traffic).toFixed(2)}%
                            </span>
                        </div>
                    </div>

                    <div className="summary-card">
                        <div className="summary-label">Число конверсий</div>
                        <div className={`summary-value ${summaryData.conversions < 0 ? 'negative' : 'positive'}`}>
                            <span className="summary-arrow">{summaryData.conversions < 0 ? '↓' : '↑'}</span>
                            <span className="summary-text">
                                {summaryData.conversions < 0 ? 'Упало' : 'Выросло'} на {Math.abs(summaryData.conversions).toFixed(2)}%
                            </span>
                        </div>
                    </div>

                    <div className="summary-card">
                        <div className="summary-label">Количество отказов</div>
                        <div className={`summary-value ${summaryData.bounce > 0 ? 'negative' : 'positive'}`}>
                            <span className="summary-arrow">{summaryData.bounce > 0 ? '↑' : '↓'}</span>
                            <span className="summary-text">
                                {summaryData.bounce > 0 ? 'Выросло' : 'Упало'} на {Math.abs(summaryData.bounce).toFixed(2)}%
                            </span>
                        </div>
                    </div>
                </div>

                {/* Metrica Section */}
                <div className="reports-section">
                    <h2 className="section-title">Яндекс.Метрика</h2>
                    <div className="reports-table-container">
                        <table className="reports-table">
                            <thead>
                                <tr>
                                    <th className="metric-name-col">Показатель</th>
                                    <th>{report.periods[0] || 'Месяц 1'}</th>
                                    <th>{report.periods[1] || 'Месяц 2'}</th>
                                    <th>{report.periods[2] || 'Месяц 3'}</th>
                                    <th className="efficiency-col">Эффективность, %</th>
                                </tr>
                            </thead>
                            <tbody>
                                {metrics.map((metric, index) => (
                                    <tr key={index}>
                                        <td className="metric-name">{metric.name}</td>
                                        <td>{metric.october}</td>
                                        <td>{metric.september}</td>
                                        <td>{metric.august}</td>
                                        <td className={`efficiency ${getEfficiencyClass(metric.efficiency, metric.name)}`}>
                                            <span className="efficiency-icon">{getEfficiencyIcon(metric.efficiency)}</span>
                                            {formatEfficiency(metric.efficiency)}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Direct Section */}
                {directMetrics.length > 0 && (
                    <div className="reports-section">
                        <h2 className="section-title">Яндекс.Директ</h2>
                        <div className="reports-table-container">
                            <table className="reports-table">
                                <thead>
                                    <tr>
                                        <th className="metric-name-col">Показатель</th>
                                        <th>{report.periods[0] || 'Месяц 1'}</th>
                                        <th>{report.periods[1] || 'Месяц 2'}</th>
                                        <th>{report.periods[2] || 'Месяц 3'}</th>
                                        <th className="efficiency-col">Эффективность, %</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {directMetrics.map((metric, index) => (
                                        <tr key={index}>
                                            <td className="metric-name">{metric.name}</td>
                                            <td>{metric.october}</td>
                                            <td>{metric.september}</td>
                                            <td>{metric.august}</td>
                                            <td className={`efficiency ${getEfficiencyClass(metric.efficiency, metric.name)}`}>
                                                <span className="efficiency-icon">{getEfficiencyIcon(metric.efficiency)}</span>
                                                {formatEfficiency(metric.efficiency)}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                )}
                
                {report.ai_insights && report.ai_insights.summary && (
                    <div className="ai-insights">
                        <h2>🤖 AI Анализ</h2>
                        <div className="ai-summary">
                            <p>{report.ai_insights.summary}</p>
                        </div>
                        {report.ai_insights.recommendations && report.ai_insights.recommendations.length > 0 && (
                            <div className="ai-recommendations">
                                <h3>Рекомендации:</h3>
                                <ul>
                                    {report.ai_insights.recommendations.map((rec, idx) => (
                                        <li key={idx}>{rec}</li>
                                    ))}
                                </ul>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Reports;

