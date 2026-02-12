-- 檢查 schedule_rules 表中是否有 NULL 或空的 status 值
SELECT id, status, COUNT(*) as cnt
FROM schedule_rules
WHERE status IS NULL OR status = ''
GROUP BY id, status;

-- 修復 NULL 或空的值為 CONFIRMED
UPDATE schedule_rules
SET status = 'CONFIRMED'
WHERE status IS NULL OR status = '';

-- 驗證修復後的結果
SELECT status, COUNT(*) as count
FROM schedule_rules
GROUP BY status;
