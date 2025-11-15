-- Откат инициализации базы данных
-- ВАЖНО: Нельзя вставлять в schema_audit после DROP DATABASE, так как база уже удалена
DROP DATABASE IF EXISTS `planica_bi`;