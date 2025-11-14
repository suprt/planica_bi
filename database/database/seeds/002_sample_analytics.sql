USE `planica_bi`;

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

-- Аккаунты Яндекс.Директа
INSERT INTO `direct_accounts` (`project_id`, `client_login`, `account_name`, `account_id`, `status`) VALUES 
(1, 'client-electro', 'Аккаунт электроники', 555111, 'active'),
(2, 'client-construction', 'Аккаунт строительства', 555222, 'active'),
(3, 'client-travel', 'Аккаунт туризма', 555333, 'active');

-- Кампании Директа
INSERT INTO `direct_campaigns` (`direct_account_id`, `campaign_id`, `name`, `status`) VALUES 
(1, 1001, 'Кампания бренда', 'ACTIVE'),
(1, 1002, 'Кампания по акциям', 'ACTIVE'),
(2, 2001, 'Кампание по строительству', 'ACTIVE'),
(3, 3001, 'Кампания туры за границу', 'ACTIVE');