-- =====================================================
-- TimeLedger 資料庫遷移腳本
-- 修復 schedule_exceptions 表 type 欄位問題（MySQL 8.0 相容版）
-- =====================================================

USE timeledger;

-- =====================================================
-- 步驟 1: 檢查欄位狀態
-- =====================================================
SELECT '=== 當前欄位狀態 ===' AS status;
SHOW COLUMNS FROM schedule_exceptions LIKE 'type';
SHOW COLUMNS FROM schedule_exceptions LIKE 'exception_type';

-- =====================================================
-- 步驟 2: 複製資料（如果 exception_type 為預設值且 type 有資料）
-- =====================================================
SELECT '=== 遷移資料 ===' AS status;

-- 查看當前資料狀態
SELECT id, type, exception_type, status FROM schedule_exceptions WHERE type IS NOT NULL AND type != '' LIMIT 10;

-- 複製資料：將 type 欄位的值複製到 exception_type
UPDATE schedule_exceptions
SET exception_type = type
WHERE (exception_type IS NULL OR exception_type = '' OR exception_type = 'CANCEL')
  AND type IS NOT NULL
  AND type != '';

-- 顯示遷移結果
SELECT id, type, exception_type, status FROM schedule_exceptions WHERE type IS NOT NULL AND type != '' LIMIT 10;

SELECT CONCAT('遷移完成，更新 ', COALESCE(ROW_COUNT(), 0), ' 筆資料') AS result;

-- =====================================================
-- 步驟 3: 驗證資料完整性
-- =====================================================
SELECT '=== 資料驗證 ===' AS status;

-- 檢查 exception_type 的分布
SELECT exception_type, COUNT(*) as count
FROM schedule_exceptions
GROUP BY exception_type;

-- 檢查是否有 type 欄位還有資料但 exception_type 沒有對應
SELECT COUNT(*) as remaining
FROM schedule_exceptions
WHERE type IS NOT NULL AND type != ''
  AND (exception_type IS NULL OR exception_type = '' OR exception_type = 'CANCEL');

-- =====================================================
-- 步驟 4: 刪除舊的 type 欄位（確認資料已正確遷移後執行）
-- =====================================================
SELECT '=== 可選：刪除舊欄位 ===' AS status;
SELECT '執行以下指令刪除舊的 type 欄位:' AS instruction;
SELECT 'ALTER TABLE schedule_exceptions DROP COLUMN type;' AS command;

-- 如果確認資料已正確遷移，可以取消註解以下行來執行刪除：
-- ALTER TABLE schedule_exceptions DROP COLUMN type;
