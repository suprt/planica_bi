USE `reports`;

UPDATE `users` 
SET `password` = '$2a$10$Z4HYHlGo85tbkoSfIgtBoOCVIxmv8uglW97naTlfwIc4WdkAyrrYm' 
WHERE `email` = 'admin@test.ru';

SELECT id, email, LEFT(password, 40) as password_check FROM users WHERE email = 'admin@test.ru';

