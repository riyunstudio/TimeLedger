-- Initialize database
CREATE DATABASE IF NOT EXISTS timeledger CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE timeledger;

-- Create basic user if needed
-- CREATE USER IF NOT EXISTS 'timeledger'@'%' IDENTIFIED BY 'timeledger_password_2026';
-- GRANT ALL PRIVILEGES ON timeledger.* TO 'timeledger'@'%';
-- FLUSH PRIVILEGES;
