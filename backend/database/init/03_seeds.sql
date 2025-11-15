USE `planica_bi`;

-- Вставка тестовых проектов
INSERT INTO `projects` (`name`, `slug`, `description`, `timezone`, `currency`, `status`, `is_public`) VALUES 
('Интернет-магазин электроники', 'electronics-store', 'Продажа электроники и бытовой техники', 'Europe/Moscow', 'RUB', 'active', 1),
('Строительная компания "СтройГарант"', 'construction-company', 'Строительство жилых и коммерческих объектов', 'Europe/Moscow', 'RUB', 'active', 1),
('Туристическое агентство "Вокруг света"', 'travel-agency', 'Организация туров по всему миру', 'Europe/Moscow', 'USD', 'active', 0);

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
INSERT INTO `metrics_monthly` (`project_id`, `year`, `month`, `visits`, `users`, `bounce_rate`, `avg_session_duration_sec`, `conversions`) VALUES 
(1, 2024, 11, 15000, 12000, 35.2, 145, 450),
(2, 2024, 11, 8000, 6500, 28.7, 210, 230),
(3, 2024, 11, 12000, 9800, 32.1, 178, 340);

-- Метрики по возрастам
INSERT INTO `metrics_age_monthly` (`project_id`, `year`, `month`, `age_group`, `visits`, `users`, `bounce_rate`, `avg_session_duration_sec`) VALUES 
(1, 2024, 11, '18-24', 3000, 2400, 40.1, 120),
(1, 2024, 11, '25-34', 6000, 4800, 32.5, 155),
(1, 2024, 11, '35-44', 4000, 3200, 30.2, 165),
(1, 2024, 11, '45-54', 1500, 1200, 35.7, 140),
(1, 2024, 11, '55+', 500, 400, 42.3, 110);

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

-- Запись в аудит о загрузке данных
INSERT INTO `schema_audit` (`change_type`, `object_type`, `object_name`, `sql_statement`) 
VALUES ('OTHER', 'TABLE', 'all', 'Test data loaded successfully');