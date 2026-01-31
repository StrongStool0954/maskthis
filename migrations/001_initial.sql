-- maskthis.com URL Shortener Database Schema
-- Initial migration

-- URLs table
CREATE TABLE IF NOT EXISTS urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    clicks INTEGER DEFAULT 0,
    created_by_ip TEXT
);

CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_created_at ON urls(created_at DESC);

-- Click analytics table
CREATE TABLE IF NOT EXISTS clicks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT NOT NULL,
    clicked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT,
    country_code TEXT,

    FOREIGN KEY (short_code) REFERENCES urls(short_code)
);

CREATE INDEX IF NOT EXISTS idx_clicks_short_code ON clicks(short_code, clicked_at DESC);
