-- Создание базы данных, если не существует
-- Имя базы берется из переменной окружения MYSQL_DATABASE (по умолчанию planica_bi)
CREATE DATABASE IF NOT EXISTS `planica_bi` 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- Создание пользователя, если не существует
-- Имя пользователя и пароль берутся из переменных окружения MYSQL_USER и MYSQL_PASSWORD
CREATE USER IF NOT EXISTS 'planica_user'@'%' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON `planica_bi`.* TO 'planica_user'@'%';
FLUSH PRIVILEGES;

-- Использование базы данных
USE `planica_bi`;