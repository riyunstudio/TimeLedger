-- Migration: Add note column to personal_events table
-- Date: 2026-01-28

-- 新增 note 欄位用於儲存個人行程備註
ALTER TABLE `personal_events`
ADD COLUMN `note` TEXT NULL DEFAULT NULL AFTER `color_hex`;

-- 如果欄位已存在則忽略（確保 idempotent）
-- 這個遷移是 idempotent 的，因為 ADD COLUMN 如果欄位已存在會報錯
-- 實際環境中請手動執行或在程式碼中處理
