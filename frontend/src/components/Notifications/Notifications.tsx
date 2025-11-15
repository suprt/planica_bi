import React, { useState, useRef } from 'react';
import './Notifications.css';
import { Notification } from '../../types';
import { useClickOutside } from '../../hooks';

interface NotificationsProps {
    notifications?: Notification[];
}

const Notifications: React.FC<NotificationsProps> = ({ notifications: propNotifications }) => {
    const [notificationsOpen, setNotificationsOpen] = useState<boolean>(false);
    const notificationsRef = useRef<HTMLDivElement>(null);

    const defaultNotifications: Notification[] = [
        { id: 1, text: 'Новая задача от Ивана', time: '10:30', type: 'info' },
        { id: 2, text: 'Заказ №2456 выполнен', time: '09:45', type: 'success' },
        { id: 3, text: 'Поступил новый отзыв', time: '09:15', type: 'info' },
        { id: 4, text: 'Обновление системы', time: 'Вчера', type: 'warning' },
    ];

    const notifications = propNotifications || defaultNotifications;

    useClickOutside(notificationsRef, () => setNotificationsOpen(false));

    const toggleNotifications = () => {
        setNotificationsOpen(!notificationsOpen);
    };

    return (
        <div className="notifications" ref={notificationsRef}>
            <button className="notification-bell" onClick={toggleNotifications}>
                🔔
                {notifications.length > 0 && (
                    <span className="notification-badge">{notifications.length}</span>
                )}
            </button>

            {notificationsOpen && (
                <div className="notifications-dropdown">
                    <div className="notifications-header">
                        <h3>Уведомления</h3>
                        <span className="notifications-count">{notifications.length}</span>
                    </div>

                    <div className="notifications-list">
                        {notifications.map(notification => (
                            <div key={notification.id} className="notification-item">
                                <div className="notification-text">{notification.text}</div>
                                <div className="notification-time">{notification.time}</div>
                            </div>
                        ))}
                    </div>

                    <div className="notifications-footer">
                        <button className="view-all-btn">Показать все</button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Notifications;

