-- Создание базы данных, если не существует
CREATE DATABASE IF NOT EXISTS `planica_bi` 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- Создание пользователя, если не существует
CREATE USER IF NOT EXISTS 'planica_user'@'%' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON `planica_bi`.* TO 'planica_user'@'%';
FLUSH PRIVILEGES;

-- Использование базы данных
USE `planica_bi`;