USE `planica_bi`;

-- Вставка тестовых проектов
INSERT INTO `projects` (`name`, `slug`, `description`, `timezone`, `currency`, `status`, `is_public`) VALUES 
('Интернет-магазин электроники', 'electronics-store', 'Продажа электроники и бытовой техники', 'Europe/Moscow', 'RUB', 'active', 1),
('Строительная компания "СтройГарант"', 'construction-company', 'Строительство жилых и коммерческих объектов', 'Europe/Moscow', 'RUB', 'active', 1),
('Туристическое агентство "Вокруг света"', 'travel-agency', 'Организация туров по всему миру', 'Europe/Moscow', 'USD', 'active', 0);