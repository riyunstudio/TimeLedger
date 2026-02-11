-- =====================================================
-- TimeLedger 資料庫遷移腳本
-- 新增 offerings 表的 default_start_time 和 default_end_time 欄位
-- =====================================================

USE timeledger;

-- =====================================================
-- 步驟 1: 檢查當前欄位狀態
-- =====================================================
SELECT '=== 當前欄位狀態 ===' AS status;
SHOW COLUMNS FROM offerings LIKE 'default_start_time';
SHOW COLUMNS FROM offerings LIKE 'default_end_time';

-- =====================================================
-- 步驟 2: 新增 default_start_time 欄位（如果不存在）
-- =====================================================
SELECT '=== 新增 default_start_time 欄位 ===' AS status;

ALTER TABLE offerings
ADD COLUMN IF NOT EXISTS default_start_time VARCHAR(5) NOT NULL DEFAULT '09:00' COMMENT '預設開始時間 (HH:MM)'
AFTER default_teacher_id;

SELECT CONCAT('新增 default_start_time 欄位結果: ', (CASE WHEN (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'timeledger' AND TABLE_NAME = 'offerings' AND COLUMN_NAME = 'default_start_time') > 0 THEN '成功' ELSE '欄位已存在或新增失敗' END)) AS result;

-- =====================================================
-- 步驟 3: 新增 default_end_time 欄位（如果不存在）
-- =====================================================
SELECT '=== 新增 default_end_time 欄位 ===' AS status;

ALTER TABLE offerings
ADD COLUMN IF NOT EXISTS default_end_time VARCHAR(5) NOT NULL DEFAULT '10:00' COMMENT '預設結束時間 (HH:MM)'
AFTER default_start_time;

SELECT CONCAT('新增 default_end_time 欄位結果: ', (CASE WHEN (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'timeledger' AND TABLE_NAME = 'offerings' AND COLUMN_NAME = 'default_end_time') > 0 THEN '成功' ELSE '欄位已存在或新增失敗' END)) AS result;

-- =====================================================
-- 步驟 4: 驗證欄位已新增
-- =====================================================
SELECT '=== 驗證新增後的欄位 ===' AS status;
DESCRIBE offerings;

-- =====================================================
-- 步驟 5: 更新現有資料的預設時間（可選）
-- =====================================================
SELECT '=== 更新現有資料（可選）===' AS status;

-- 如果需要為現有的 offerings 設置不同的預設時間，可以取消註解以下行：
-- UPDATE offerings SET default_start_time = '09:00', default_end_time = '10:00' WHERE default_start_time IS NULL OR default_end_time IS NULL;

SELECT '現有資料保持不變，如需批量更新請手動執行 UPDATE 語句' AS note;
