USE `reports`;

-- SEO данные
INSERT INTO `seo_queries_monthly` (`project_id`, `year`, `month`, `query`, `position`, `url`, `impressions`, `clicks`) VALUES 
(1, 2024, 11, 'купить смартфон', 5.2, '/catalog/smartphones', 1500, 120),
(1, 2024, 11, 'ноутбук недорого', 8.7, '/catalog/laptops', 800, 45),
(1, 2024, 11, 'телевизор samsung', 3.1, '/catalog/tv', 2500, 210),
(1, 2024, 11, 'наушники беспроводные', 12.5, '/catalog/headphones', 600, 25),
(2, 2024, 11, 'строительство домов', 2.3, '/services/house-building', 1200, 180),
(2, 2024, 11, 'ремонт квартир', 7.8, '/services/apartment-renovation', 900, 65),
(3, 2024, 11, 'туры в турцию', 4.5, '/tours/turkey', 1800, 150),
(3, 2024, 11, 'отдых в сочи', 9.2, '/tours/sochi', 700, 40);
