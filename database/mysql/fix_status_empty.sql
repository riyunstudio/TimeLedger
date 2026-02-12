-- 檢查 schedule_rules 表中所有記錄的 status 值
SELECT id, status, COUNT(*) as cnt
FROM schedule_rules
GROUP BY id, status;

-- 檢查是否有 NULL 或空的 status
SELECT id, status FROM schedule_rules WHERE status IS NULL OR status = '';

-- 更新所有空值為 CONFIRMED
UPDATE schedule_rules
SET status = 'CONFIRMED'
WHERE status IS NULL OR status = '';

-- 再次確認更新結果
SELECT status, COUNT(*) as count
FROM schedule_rules
GROUP BY status;
