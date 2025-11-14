-- Initial database schema
-- TODO: Implement all tables according to TZ

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    timezone VARCHAR(50) DEFAULT 'Europe/Moscow',
    currency ENUM('RUB') DEFAULT 'RUB',
    is_active BOOLEAN DEFAULT TRUE,
    created_at BIGINT,
    updated_at BIGINT,
    INDEX idx_slug (slug),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Yandex counters table
CREATE TABLE IF NOT EXISTS yandex_counters (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    counter_id BIGINT NOT NULL,
    name VARCHAR(255),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at BIGINT,
    updated_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id),
    INDEX idx_counter_id (counter_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Direct accounts table
CREATE TABLE IF NOT EXISTS direct_accounts (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    client_login VARCHAR(255) NOT NULL,
    account_name VARCHAR(255),
    created_at BIGINT,
    updated_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Direct campaigns table
CREATE TABLE IF NOT EXISTS direct_campaigns (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    direct_account_id INT UNSIGNED NOT NULL,
    campaign_id BIGINT NOT NULL,
    name VARCHAR(255),
    status VARCHAR(50),
    created_at BIGINT,
    updated_at BIGINT,
    FOREIGN KEY (direct_account_id) REFERENCES direct_accounts(id) ON DELETE CASCADE,
    INDEX idx_direct_account_id (direct_account_id),
    INDEX idx_campaign_id (campaign_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Goals table
CREATE TABLE IF NOT EXISTS goals (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    counter_id INT UNSIGNED NOT NULL,
    goal_id BIGINT NOT NULL,
    name VARCHAR(255),
    is_conversion BOOLEAN DEFAULT FALSE,
    created_at BIGINT,
    updated_at BIGINT,
    FOREIGN KEY (counter_id) REFERENCES yandex_counters(id) ON DELETE CASCADE,
    INDEX idx_counter_id (counter_id),
    INDEX idx_goal_id (goal_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Metrics monthly table
CREATE TABLE IF NOT EXISTS metrics_monthly (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    year INT NOT NULL,
    month TINYINT NOT NULL,
    visits INT DEFAULT 0,
    users INT DEFAULT 0,
    bounce_rate DECIMAL(5,2),
    avg_session_duration_sec INT DEFAULT 0,
    conversions INT,
    created_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE KEY idx_project_year_month (project_id, year, month),
    INDEX idx_project_id (project_id),
    INDEX idx_year_month (year, month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Metrics age monthly table
CREATE TABLE IF NOT EXISTS metrics_age_monthly (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    year INT NOT NULL,
    month TINYINT NOT NULL,
    age_group ENUM('18-24', '25-34', '35-44', '45-54', '55+', 'unknown') NOT NULL,
    visits INT DEFAULT 0,
    users INT DEFAULT 0,
    bounce_rate DECIMAL(5,2),
    avg_session_duration_sec INT DEFAULT 0,
    created_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE KEY idx_project_year_month_age (project_id, year, month, age_group),
    INDEX idx_project_id (project_id),
    INDEX idx_year_month (year, month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Direct campaign monthly table
CREATE TABLE IF NOT EXISTS direct_campaign_monthly (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    direct_campaign_id INT UNSIGNED NOT NULL,
    year INT NOT NULL,
    month TINYINT NOT NULL,
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    ctr_pct DECIMAL(6,2),
    cpc DECIMAL(12,2),
    conversions INT,
    cpa DECIMAL(12,2),
    cost DECIMAL(14,2) DEFAULT 0,
    created_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (direct_campaign_id) REFERENCES direct_campaigns(id) ON DELETE CASCADE,
    UNIQUE KEY idx_campaign_year_month (direct_campaign_id, year, month),
    INDEX idx_project_id (project_id),
    INDEX idx_year_month (year, month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Direct totals monthly table
CREATE TABLE IF NOT EXISTS direct_totals_monthly (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    year INT NOT NULL,
    month TINYINT NOT NULL,
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    ctr_pct DECIMAL(6,2),
    cpc DECIMAL(12,2),
    conversions INT,
    cpa DECIMAL(12,2),
    cost DECIMAL(14,2) DEFAULT 0,
    created_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE KEY idx_project_year_month (project_id, year, month),
    INDEX idx_project_id (project_id),
    INDEX idx_year_month (year, month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- SEO queries monthly table
CREATE TABLE IF NOT EXISTS seo_queries_monthly (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id INT UNSIGNED NOT NULL,
    year INT NOT NULL,
    month TINYINT NOT NULL,
    query VARCHAR(500),
    position INT,
    url VARCHAR(1000),
    created_at BIGINT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id),
    INDEX idx_year_month (year, month),
    INDEX idx_query (query(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

