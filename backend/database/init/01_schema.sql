-- Создание базы данных, если не существует
-- Имя базы берется из переменной окружения MYSQL_DATABASE (по умолчанию reports)
CREATE DATABASE IF NOT EXISTS `reports` 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- Создание пользователя, если не существует
-- Имя пользователя и пароль берутся из переменных окружения MYSQL_USER и MYSQL_PASSWORD
CREATE USER IF NOT EXISTS 'planica_user'@'%' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON `reports`.* TO 'planica_user'@'%';
FLUSH PRIVILEGES;

-- Использование базы данных
USE `reports`;

-- Установка кодировки UTF-8
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;