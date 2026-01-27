-- =====================================================
-- 資料庫狀態檢查腳本
-- 執行方式: 在 MySQL 中執行
-- =====================================================

USE timeledger;

-- 1. 檢查 admin_users 表的 center_id
SELECT
    id,
    email,
    center_id,
    name,
    role,
    status,
    CASE
        WHEN center_id = 0 THEN '❌ 異常：center_id 為 0'
        WHEN center_id IS NULL THEN '❌ 異常：center_id 為 NULL'
        WHEN center_id > 0 THEN '✅ 正常'
        ELSE '❌ 異常'
    END AS status_check
FROM admin_users;

-- 2. 檢查 centers 表是否有資料
SELECT
    id,
    name,
    plan_level,
    CASE
        WHEN COUNT(*) > 0 THEN '✅ 中心資料存在'
        ELSE '❌ 異常：無中心資料'
    END AS status_check
FROM centers
GROUP BY id;

-- 3. 檢查 admin 與 center 的關聯是否正確
SELECT
    a.id,
    a.email,
    a.center_id AS admin_center_id,
    c.id AS center_id,
    c.name AS center_name,
    CASE
        WHEN a.center_id = c.id THEN '✅ 關聯正確'
        ELSE '❌ 關聯錯誤'
    END AS relation_check
FROM admin_users a
LEFT JOIN centers c ON a.center_id = c.id;

-- 4. 快速修復：如果發現 center_id = 0 或 NULL
-- 請先確保 centers 表有資料後，執行以下更新：
-- UPDATE admin_users SET center_id = 1 WHERE center_id = 0 OR center_id IS NULL;
