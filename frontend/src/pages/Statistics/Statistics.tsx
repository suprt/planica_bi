import React from 'react';
import './Statistics.css';

interface MetricRow {
    name: string;
    october: string | number;
    september: string | number;
    august: string | number;
    efficiency: number;
}

const Statistics: React.FC = () => {
    const metrics: MetricRow[] = [
        { name: 'Посетители, кол-во', october: 3246, september: 3971, august: 3416, efficiency: -18.26 },
        { name: 'Новые посетители, кол-во', october: 2958, september: 3642, august: 3134, efficiency: -18.78 },
        { name: 'Визиты, кол-во', october: 4268, september: 4926, august: 4356, efficiency: -13.36 },
        { name: 'Кол-во отказов, %', october: '28.66%', september: '22.78%', august: '26.86%', efficiency: 25.81 },
        { name: 'Время на сайте, сек', october: '1:42', september: '1:24', august: '1:18', efficiency: 21.43 },
        { name: 'Всего заявок, кол-во', october: 83, september: 98, august: 81, efficiency: -15.31 },
        { name: 'Клик на номер, кол-во', october: 26, september: 308, august: 109, efficiency: -91.56 },
        { name: 'Заявка', october: 12, september: 12, august: 8, efficiency: 0.00 },
        { name: 'Спецпредложение', october: 29, september: 29, august: 24, efficiency: 0.00 },
        { name: 'Лизинг', october: 0, september: 1, august: 1, efficiency: -100.00 },
        { name: 'Тест-драйв', october: 4, september: 3, august: 6, efficiency: 33.33 },
        { name: 'Госпрограмма', october: 1, september: 5, august: 6, efficiency: -80.00 },
        { name: 'Звонок CallTouch', october: 37, september: 41, august: 32, efficiency: -9.76 },
        { name: 'Таймер', october: 0, september: 7, august: 4, efficiency: -100.00 },
    ];

    const formatEfficiency = (value: number): string => {
        if (value === 0) return '0.00%';
        const sign = value > 0 ? '+' : '';
        return `${sign}${value.toFixed(2)}%`;
    };

    // Определяем, является ли метрика "хорошей" (увеличение = хорошо)
    const isPositiveMetric = (metricName: string): boolean => {
        const positiveMetrics = [
            'Время на сайте',
            'Всего заявок',
            'Клик на номер',
            'Заявка',
            'Спецпредложение',
            'Лизинг',
            'Тест-драйв',
            'Госпрограмма',
            'Звонок CallTouch',
            'Таймер',
            'Посетители',
            'Новые посетители',
            'Визиты'
        ];
        return positiveMetrics.some(metric => metricName.includes(metric));
    };

    const getEfficiencyClass = (value: number, metricName: string): string => {
        if (value === 0) return 'efficiency-neutral';
        const isPositive = isPositiveMetric(metricName);
        // Если метрика "хорошая", то положительное изменение = хорошо (зеленый)
        // Если метрика "плохая" (например, отказы), то положительное изменение = плохо (красный)
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

    return (
        <div className="statistics-page">
            <div className="statistics-header">
                <h1 className="page-title">Статистика / Аналитика сайта</h1>
            </div>

            <div className="statistics-content">
                <div className="statistics-summary">
                    <div className="summary-card">
                        <div className="summary-label">Трафик</div>
                        <div className="summary-value negative">
                            <span className="summary-arrow">↓</span>
                            <span className="summary-text">Упал на 13%</span>
                        </div>
                    </div>

                    <div className="summary-card">
                        <div className="summary-label">Число конверсий</div>
                        <div className="summary-value negative">
                            <span className="summary-arrow">↓</span>
                            <span className="summary-text">Упало на 15%</span>
                        </div>
                    </div>

                    <div className="summary-card">
                        <div className="summary-label">Количество отказов</div>
                        <div className="summary-value negative">
                            <span className="summary-arrow">↑</span>
                            <span className="summary-text">Выросло на 26%</span>
                        </div>
                    </div>
                </div>

                <div className="statistics-table-container">
                    <table className="statistics-table">
                        <thead>
                            <tr>
                                <th className="metric-name-col">Показатель</th>
                                <th>Октябрь</th>
                                <th>Сентябрь</th>
                                <th>Август</th>
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
        </div>
    );
};

export default Statistics;

