-- Seeds для таблицы direct_campaign_monthly
-- Тестовые данные по метрикам кампаний Яндекс.Директ по месяцам
-- Распределяем данные из direct_totals_monthly по кампаниям из direct_campaigns

USE `reports`;

-- Установка кодировки UTF-8 для корректного отображения русских символов
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Вставка тестовых данных по метрикам кампаний (помесячно)
-- Распределяем общие метрики из direct_totals_monthly по кампаниям из direct_campaigns
-- Пропорции распределения основаны на campaign_id из direct_campaigns
-- Примечание: GORM использует множественное число для названия таблицы (direct_campaign_monthlies)
INSERT INTO `direct_campaign_monthlies` (
    `project_id`, `direct_campaign_id`, `year`, `month`, 
    `impressions`, `clicks`, `ctr_pct`, `cpc`, 
    `conversions`, `cpa`, `cost`, `created_at`
)
SELECT 
    dc_project.project_id,
    dc.id AS direct_campaign_id,
    dt.year,
    dt.month,
    -- Распределяем impressions пропорционально по кампаниям
    -- Для активных кампаний больше доля, для paused меньше
    FLOOR(dt.impressions * 
        CASE 
            WHEN dc.status = 'active' THEN 
                CASE 
                    -- Project 1: 5 кампаний (id 1-5), активные получают больше
                    WHEN dc.campaign_id IN (1001, 1002, 1004) THEN 0.25  -- 3 активные по 25%
                    WHEN dc.campaign_id = 1003 THEN 0.15  -- paused получает меньше
                    -- Project 2: 4 кампании (id 6-9), активные получают больше
                    WHEN dc.campaign_id IN (2001, 2002) THEN 0.35  -- 2 активные по 35%
                    WHEN dc.campaign_id = 2003 THEN 0.20  -- paused
                    -- Project 3: 5 кампаний (id 10-14), активные получают больше
                    WHEN dc.campaign_id IN (3001, 3002, 3005) THEN 0.28  -- 3 активные
                    WHEN dc.campaign_id = 3004 THEN 0.16  -- paused
                    ELSE 0.20
                END
            ELSE 0.05  -- ended кампании получают минимум
        END
    ) AS impressions,
    -- Распределяем clicks аналогично
    FLOOR(dt.clicks * 
        CASE 
            WHEN dc.status = 'active' THEN 
                CASE 
                    WHEN dc.campaign_id IN (1001, 1002, 1004) THEN 0.25
                    WHEN dc.campaign_id = 1003 THEN 0.15
                    WHEN dc.campaign_id IN (2001, 2002) THEN 0.35
                    WHEN dc.campaign_id = 2003 THEN 0.20
                    WHEN dc.campaign_id IN (3001, 3002, 3005) THEN 0.28
                    WHEN dc.campaign_id = 3004 THEN 0.16
                    ELSE 0.20
                END
            ELSE 0.05
        END
    ) AS clicks,
    -- CTR немного варьируем (±10%)
    dt.ctr_pct * (0.9 + ((dc.id % 5) * 0.05)) AS ctr_pct,
    -- CPC немного варьируем (±10%)
    dt.cpc * (0.9 + ((dc.id % 5) * 0.05)) AS cpc,
    -- Conversions распределяем пропорционально
    CASE 
        WHEN dt.conversions IS NOT NULL THEN 
            FLOOR(dt.conversions * 
                CASE 
                    WHEN dc.status = 'active' THEN 
                        CASE 
                            WHEN dc.campaign_id IN (1001, 1002, 1004) THEN 0.25
                            WHEN dc.campaign_id = 1003 THEN 0.15
                            WHEN dc.campaign_id IN (2001, 2002) THEN 0.35
                            WHEN dc.campaign_id = 2003 THEN 0.20
                            WHEN dc.campaign_id IN (3001, 3002, 3005) THEN 0.28
                            WHEN dc.campaign_id = 3004 THEN 0.16
                            ELSE 0.20
                        END
                    ELSE 0.05
                END
            )
        ELSE NULL
    END AS conversions,
    -- CPA немного варьируем (±10%)
    CASE 
        WHEN dt.cpa IS NOT NULL THEN dt.cpa * (0.9 + ((dc.id % 5) * 0.05))
        ELSE NULL
    END AS cpa,
    -- Cost распределяем пропорционально
    dt.cost * 
        CASE 
            WHEN dc.status = 'active' THEN 
                CASE 
                    WHEN dc.campaign_id IN (1001, 1002, 1004) THEN 0.25
                    WHEN dc.campaign_id = 1003 THEN 0.15
                    WHEN dc.campaign_id IN (2001, 2002) THEN 0.35
                    WHEN dc.campaign_id = 2003 THEN 0.20
                    WHEN dc.campaign_id IN (3001, 3002, 3005) THEN 0.28
                    WHEN dc.campaign_id = 3004 THEN 0.16
                    ELSE 0.20
                END
            ELSE 0.05
        END AS cost,
    dt.created_at
FROM 
    `direct_totals_monthly` dt
    INNER JOIN (
        -- Подзапрос для получения project_id для каждой кампании
        SELECT 
            dc.id,
            dc.campaign_id,
            dc.status,
            da.project_id
        FROM 
            `direct_campaigns` dc
            INNER JOIN `direct_accounts` da ON dc.direct_account_id = da.id
    ) AS dc_project ON dt.project_id = dc_project.project_id
    INNER JOIN `direct_campaigns` dc ON dc_project.id = dc.id
WHERE 
    dt.year = 2025
    AND dt.month IN (9, 10, 11)
ORDER BY 
    dc_project.project_id, dc.id, dt.year, dt.month;

