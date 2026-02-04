-- Migration: 建立管理員登入紀錄表
-- Date: 2026-02-05

-- 如果資料表不存在才建立
CREATE TABLE IF NOT EXISTS `admin_login_histories` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `admin_id` BIGINT UNSIGNED NOT NULL COMMENT '管理員 ID',
    `email` VARCHAR(255) NOT NULL COMMENT '管理員 Email',
    `status` VARCHAR(20) NOT NULL COMMENT '登入狀態: SUCCESS, FAILED',
    `ip_address` VARCHAR(45) COMMENT 'IP 位址',
    `user_agent` VARCHAR(500) COMMENT '使用者代理程式',
    `reason` VARCHAR(255) COMMENT '失敗原因',
    `created_at` DATETIME(3) NOT NULL COMMENT '建立時間',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '刪除時間',

    PRIMARY KEY (`id`),
    INDEX `idx_admin_id` (`admin_id`),
    INDEX `idx_email` (`email`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理員登入紀錄表';
