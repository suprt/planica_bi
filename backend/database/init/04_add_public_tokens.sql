-- Генерация public_token для существующих проектов без токена
USE `reports`;
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Функция для генерации случайного токена (64 hex символа)
-- Обновляем проекты без токена
UPDATE `projects` 
SET `public_token` = CONCAT(
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0'),
    LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0')
)
WHERE `public_token` IS NULL OR `public_token` = '';


