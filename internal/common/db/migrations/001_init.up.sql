-- 001_init.up.sql
-- Users
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(200) NOT NULL,
    role VARCHAR(20) DEFAULT 'employee',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Blog Posts
CREATE TABLE IF NOT EXISTS blog_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    status VARCHAR(20) DEFAULT 'draft',
    views INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Media Files
CREATE TABLE IF NOT EXISTS media_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_name VARCHAR(500),
    file_path VARCHAR(1000),
    file_size BIGINT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Settings
CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Traffic Events
CREATE TABLE IF NOT EXISTS traffic_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id VARCHAR(36),
    page VARCHAR(500),
    ip VARCHAR(45),
    location VARCHAR(255),
    device VARCHAR(100),
    duration INTEGER DEFAULT 0,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_traffic_session ON traffic_events(session_id);
CREATE INDEX IF NOT EXISTS idx_traffic_page ON traffic_events(page);
CREATE INDEX IF NOT EXISTS idx_traffic_timestamp ON traffic_events(timestamp);

-- Seed default admin
INSERT OR IGNORE INTO users (id, username, password, name, role) VALUES 
(1, 'admin', '$2a$10$default_hash_placeholder_will_be_replaced', 'Quản trị viên', 'admin'),
(2, 'nhanvien', '$2a$10$default_hash_placeholder_will_be_replaced', 'Nhân viên kinh doanh', 'employee'),
(3, 'marketing', '$2a$10$default_hash_placeholder_will_be_replaced', 'Nhân viên marketing', 'employee');
