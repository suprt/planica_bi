USE `reports`;

-- Установка кодировки UTF-8 для корректного отображения русских символов
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Вставка тестовых проектов (соответствует GORM модели)
INSERT INTO `projects` (`name`, `slug`, `timezone`, `currency`, `is_active`) VALUES 
('Интернет-магазин электроники', 'electronics-store', 'Europe/Moscow', 'RUB', 1),
('Строительная компания "СтройГарант"', 'construction-company', 'Europe/Moscow', 'RUB', 1),
('Туристическое агентство "Вокруг света"', 'travel-agency', 'Europe/Moscow', 'RUB', 1);

-- Яндекс.Метрика счетчики (соответствует структуре GORM модели)
INSERT INTO `yandex_counters` (`project_id`, `counter_id`, `name`, `is_primary`) VALUES 
(1, 12345678, 'Основной счетчик', 1),
(2, 87654321, 'Счетчик стройки', 1),
(3, 11223344, 'Счетчик туров', 1);

-- Цели Яндекс.Метрики (соответствует структуре GORM модели)
INSERT INTO `goals` (`counter_id`, `goal_id`, `name`, `is_conversion`) VALUES 
(1, 111, 'Оформление заказа', 1),
(1, 112, 'Добавление в корзину', 0),
(2, 211, 'Отправка заявки', 1),
(3, 311, 'Бронирование тура', 1);

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

-- SEO данные (только поля, которые есть в GORM модели)
INSERT INTO `seo_queries_monthly` (`project_id`, `year`, `month`, `query`, `position`, `url`) VALUES 
(1, 2024, 11, 'купить смартфон', 5, '/catalog/smartphones'),
(1, 2024, 11, 'ноутбук недорого', 9, '/catalog/laptops'),
(1, 2024, 11, 'телевизор samsung', 3, '/catalog/tv'),
(1, 2024, 11, 'наушники беспроводные', 13, '/catalog/headphones'),
(2, 2024, 11, 'строительство домов', 2, '/services/house-building'),
(2, 2024, 11, 'ремонт квартир', 8, '/services/apartment-renovation'),
(3, 2024, 11, 'туры в турцию', 5, '/tours/turkey'),
(3, 2024, 11, 'отдых в сочи', 9, '/tours/sochi');

-- Запись в аудит о загрузке данных
INSERT INTO `schema_audit` (`change_type`, `object_type`, `object_name`, `sql_statement`) 
VALUES ('OTHER', 'TABLE', 'all', 'Test data loaded successfully');