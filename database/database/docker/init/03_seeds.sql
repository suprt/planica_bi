USE `planica_bi`;

-- Вставка тестовых проектов
INSERT INTO `projects` (`name`, `slug`, `description`, `timezone`, `currency`, `status`, `is_public`) VALUES 
('Интернет-магазин электроники', 'electronics-store', 'Продажа электроники и бытовой техники', 'Europe/Moscow', 'RUB', 'active', 1),
('Строительная компания "СтройГарант"', 'construction-company', 'Строительство жилых и коммерческих объектов', 'Europe/Moscow', 'RUB', 'active', 1),
('Туристическое агентство "Вокруг света"', 'travel-agency', 'Организация туров по всему миру', 'Europe/Moscow', 'USD', 'active', 0);

-- Вставка тестовых пользователей
INSERT INTO `users` (`name`, `email`, `password`, `timezone`, `language`, `is_active`) VALUES 
('Администратор системы', 'admin@planica.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1),
('Менеджер проектов', 'manager@planica.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1),
('Клиент Иванов', 'client@electronics-store.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1),
('Клиент Петров', 'client@construction-company.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1);

-- Назначение ролей
INSERT INTO `user_project_roles` (`user_id`, `project_id`, `role`) VALUES 
(1, 1, 'admin'),   -- Админ на все проекты
(1, 2, 'admin'),
(1, 3, 'admin'),
(2, 1, 'manager'), -- Менеджер на проекты
(2, 2, 'manager'),
(3, 1, 'client'),  -- Клиент только на свой проект
(4, 2, 'client');  -- Клиент только на свой проект

-- Яндекс.Метрика счетчики
INSERT INTO `yandex_counters` (`project_id`, `counter_id`, `name`, `site_name`, `is_primary`, `status`) VALUES 
(1, 12345678, 'Основной счетчик', 'electronics-store.ru', 1, 'active'),
(2, 87654321, 'Счетчик стройки', 'construction-company.com', 1, 'active'),
(3, 11223344, 'Счетчик туров', 'travel-agency.ru', 1, 'active');

-- Цели Яндекс.Метрики
INSERT INTO `goals` (`counter_id`, `goal_id`, `name`, `goal_type`, `is_conversion`) VALUES 
(1, 111, 'Оформление заказа', 'url', 1),
(1, 112, 'Добавление в корзину', 'action', 0),
(2, 211, 'Отправка заявки', 'url', 1),
(3, 311, 'Бронирование тура', 'url', 1);

-- Основные метрики
INSERT INTO `metrics_monthly` (`project_id`, `year`, `month`, `visits`, `users`, `bounce_rate`, `avg_session_duration_sec`, `conversions`, `source`) VALUES 
(1, 2024, 11, 15000, 12000, 35.2, 145, 450, 'all'),
(1, 2024, 11, 2500, 2000, 32.1, 160, 75, 'organic'),
(2, 2024, 11, 8000, 6500, 28.7, 210, 230, 'all'),
(3, 2024, 11, 12000, 9800, 32.1, 178, 340, 'all');

-- Метрики по возрастам
INSERT INTO `metrics_age_monthly` (`project_id`, `year`, `month`, `age_group`, `visits`, `users`, `bounce_rate`, `avg_session_duration_sec`, `source`) VALUES 
(1, 2024, 11, '18-24', 3000, 2400, 40.1, 120, 'all'),
(1, 2024, 11, '25-34', 6000, 4800, 32.5, 155, 'all'),
(1, 2024, 11, '35-44', 4000, 3200, 30.2, 165, 'all'),
(1, 2024, 11, '45-54', 1500, 1200, 35.7, 140, 'all'),
(1, 2024, 11, '55+', 500, 400, 42.3, 110, 'all');

-- SEO данные
INSERT INTO `seo_queries_monthly` (`project_id`, `year`, `month`, `query`, `position`, `url`, `impressions`, `clicks`) VALUES 
(1, 2024, 11, 'купить смартфон', 5.2, '/catalog/smartphones', 1500, 120),
(1, 2024, 11, 'ноутбук недорого', 8.7, '/catalog/laptops', 800, 45),
(1, 2024, 11, 'телевизор samsung', 3.1, '/catalog/tv', 2500, 210),
(1, 2024, 11, 'наушники беспроводные', 12.5, '/catalog/headphones', 600, 25),
(2, 2024, 11, 'строительство домов', 2.3, '/services/house-building', 1200, 180),
(2, 2024, 11, 'ремонт квартир', 7.8, '/services/apartment-renovation', 900, 65),
(3, 2024, 11, 'туры в турцию', 4.5, '/tours/turkey', 1800, 150),
(3, 2024, 11, 'отдых в сочи', 9.2, '/tours/sochi', 700, 40);

-- SEO сводки
INSERT INTO `seo_summary_monthly` (`project_id`, `year`, `month`, `total_queries`, `avg_position`, `total_impressions`, `total_clicks`, `ctr_pct`) VALUES 
(1, 2024, 11, 45, 7.32, 5400, 400, 7.41),
(2, 2024, 11, 23, 5.15, 2100, 245, 11.67),
(3, 2024, 11, 32, 6.78, 2500, 190, 7.60);

-- Запись в аудит о загрузке данных
INSERT INTO `schema_audit` (`change_type`, `object_type`, `object_name`, `sql_statement`) 
VALUES ('OTHER', 'TABLE', 'all', 'Test data loaded successfully');