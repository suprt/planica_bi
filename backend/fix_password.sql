USE `reports`;

UPDATE `users` 
SET `password` = '$2a$10$Z4HYHlGo85tbkoSfIgtBoOCVIxmv8uglW97naTlfwIc4WdkAyrrYm' 
WHERE `email` = 'admin@test.ru';

