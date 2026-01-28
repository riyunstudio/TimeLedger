-- =====================================================
-- TimeLedger 資料庫遷移腳本
-- 修復 schedule_exceptions 表 type 欄位問題
-- 執行方式: 登入 MySQL 後，選擇 timeledger 資料庫，執行此檔案
-- =====================================================

USE timeledger;

-- =====================================================
-- 步驟 1: 檢查是否存在舊的 type 欄位
-- =====================================================
-- 檢查 schedule_exceptions 表結構
SHOW COLUMNS FROM schedule_exceptions LIKE 'type';
SHOW COLUMNS FROM schedule_exceptions LIKE 'exception_type';

-- =====================================================
-- 步驟 2: 如果存在 type 欄位且不存在 exception_type 欄位，進行遷移
-- =====================================================

-- 2.1 新增 exception_type 欄位（如果尚不存在）
ALTER TABLE schedule_exceptions
ADD COLUMN IF NOT EXISTS exception_type VARCHAR(20) NOT NULL DEFAULT 'CANCEL' AFTER original_date;

-- 2.2 複製 type 欄位的資料到 exception_type 欄位
UPDATE schedule_exceptions
SET exception_type = type
WHERE exception_type = 'CANCEL' AND type IS NOT NULL AND type != '';

-- 2.3 確認資料遷移完成後，刪除舊的 type 欄位（如果存在）
-- 注意：此步驟需要確認資料已正確遷移才執行
-- ALTER TABLE schedule_exceptions DROP COLUMN IF EXISTS type;

-- =====================================================
-- 驗證遷移結果
-- =====================================================
-- 查詢範例資料確認欄位正確
SELECT id, exception_type, status, original_date, reason
FROM schedule_exceptions
ORDER BY id DESC
LIMIT 10;

-- =====================================================
-- 說明
-- =====================================================
-- 此遷移腳本完成以下操作：
-- 1. 新增 exception_type 欄位（使用明確的欄位名稱）
-- 2. 將舊的 type 欄位資料複製到新的 exception_type 欄位
-- 3. 保留舊的 type 欄位以便驗證資料完整性後再刪除
--
-- 驗證完成後，可以執行以下指令刪除舊欄位：
-- ALTER TABLE schedule_exceptions DROP COLUMN IF EXISTS type;
-- =====================================================
