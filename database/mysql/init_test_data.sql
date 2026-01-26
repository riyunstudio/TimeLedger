-- =====================================================
-- TimeLedger 測試環境初始化腳本
-- 執行方式: 登入 MySQL 後，選擇資料庫，執行此檔案
-- =====================================================

USE timeledger;

-- -------------------------------------------------
-- 清除現有資料（可選，如果需要乾淨的環境）
-- -------------------------------------------------
-- 注意：請依照順序執行，先刪除子表再刪除主表
DELETE FROM schedule_exceptions;
DELETE FROM schedule_rules;
DELETE FROM center_holidays;
DELETE FROM center_invitations;
DELETE FROM offerings;
DELETE FROM courses;
DELETE FROM rooms;
DELETE FROM center_memberships;
DELETE FROM admin_users;
DELETE FROM teachers;
DELETE FROM centers;

-- -------------------------------------------------
-- 1. 城市資料（精簡版）
-- -------------------------------------------------
INSERT INTO geo_cities (id, name) VALUES
(1, '臺北市'),
(2, '新北市'),
(3, '桃園市'),
(4, '臺中市'),
(5, '臺南市'),
(6, '高雄市');

INSERT INTO geo_districts (city_id, name) VALUES
-- 臺北市
(1, '中正區'), (1, '大同區'), (1, '中山區'), (1, '松山區'),
(1, '大安區'), (1, '萬華區'), (1, '信義區'), (1, '士林區'),
(1, '北投區'), (1, '內湖區'), (1, '南港區'), (1, '文山區'),
-- 新北市
(2, '板橋區'), (2, '中和區'), (2, '永和區'), (2, '新莊區'),
(2, '新店區'), (2, '淡水區'), (2, '汐止區'), (2, '土城區'),
-- 桃園市
(3, '桃園區'), (3, '中壢區'), (3, '平鎮區'), (3, '八德區'),
(3, '楊梅區'), (3, '蘆竹區'), (3, '龜山區'), (3, '龍潭區'),
-- 臺中市
(4, '中區'), (4, '東區'), (4, '南區'), (4, '西區'),
(4, '北區'), (4, '西屯區'), (4, '南屯區'), (4, '北屯區'),
-- 臺南市
(5, '中西區'), (5, '東區'), (5, '南區'), (5, '北區'),
(5, '安平區'), (5, '安南區'),
-- 高雄市
(6, '新興區'), (6, '前金區'), (6, '苓雅區'), (6, '鹽埕區'),
(6, '鼓山區'), (6, '左營區'), (6, '楠梓區'), (6, '三民區');

-- -------------------------------------------------
-- 2. 建立中心 (1個)
-- -------------------------------------------------
INSERT INTO centers (id, name, plan_level, settings, created_at) VALUES
(1, 'TimeLedger 旗艦館', 'STARTER', '{"allow_public_register": true, "default_language": "zh-TW", "exception_lead_days": 14}', NOW());

-- -------------------------------------------------
-- 3. 建立管理員 (2個)
-- 密碼: admin123
-- bcrypt hash: $2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi (admin123)
-- -------------------------------------------------
INSERT INTO admin_users (id, center_id, email, password_hash, name, role, status, created_at, updated_at) VALUES
(1, 1, 'admin@timeledger.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '超級管理員', 'OWNER', 'ACTIVE', NOW(), NOW()),
(2, 1, 'staff@timeledger.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '工作人員', 'STAFF', 'ACTIVE', NOW(), NOW());

-- -------------------------------------------------
-- 4. 建立老師 (3個)
-- -------------------------------------------------
INSERT INTO teachers (id, line_user_id, name, email, bio, city, district, is_open_to_hiring, created_at, updated_at) VALUES
(1, 'U000000000000000000000001', '李小美', 'sarah@example.com', '資深瑜伽老師，專長哈達瑜伽和熱瑜伽', '臺北市', '大安區', true, NOW(), NOW()),
(2, 'U000000000000000000000002', '陳大文', 'david@example.com', '健身教練，專長肌力訓練和功能性訓練', '新北市', '板橋區', true, NOW(), NOW()),
(3, 'U000000000000000000000003', '王小花', 'flower@example.com', '舞蹈老師，專長芭蕾和現代舞', '桃園市', '桃園區', false, NOW(), NOW());

-- -------------------------------------------------
-- 5. 老師加入中心 (3筆)
-- -------------------------------------------------
INSERT INTO center_memberships (id, teacher_id, center_id, status, joined_at, created_at, updated_at) VALUES
(1, 1, 1, 'ACTIVE', NOW(), NOW(), NOW()),
(2, 2, 1, 'ACTIVE', NOW(), NOW(), NOW()),
(3, 3, 1, 'ACTIVE', NOW(), NOW(), NOW());

-- -------------------------------------------------
-- 6. 建立設施 (教室) - 方便測試
-- -------------------------------------------------
INSERT INTO rooms (id, center_id, name, capacity, is_active, created_at, updated_at) VALUES
(1, 1, '瑜伽教室A', 20, true, NOW(), NOW()),
(2, 1, '瑜伽教室B', 15, true, NOW(), NOW()),
(3, 1, '健身教室', 25, true, NOW(), NOW()),
(4, 1, '舞蹈教室', 30, true, NOW(), NOW());

-- -------------------------------------------------
-- 7. 建立課程 - 方便測試
-- -------------------------------------------------
INSERT INTO courses (id, center_id, name, category, description, duration_minutes, room_buffer_min, teacher_buffer_min, is_active, created_at, updated_at) VALUES
(1, 1, '哈達瑜伽', '瑜伽', '傳統瑜伽課程，適合初學者', 60, 10, 10, true, NOW(), NOW()),
(2, 1, '熱瑜伽', '瑜伽', '高溫瑜伽課程，流汗排毒', 60, 15, 15, true, NOW(), NOW()),
(3, 1, '肌力訓練', '健身', '重量訓練課程，增強肌力', 45, 5, 10, true, NOW(), NOW()),
(4, 1, '芭蕾基礎', '舞蹈', '芭蕾舞基礎課程', 90, 10, 15, true, NOW(), NOW());

-- -------------------------------------------------
-- 8. 建立班別 (Offerings) - 方便測試
-- -------------------------------------------------
INSERT INTO offerings (id, center_id, course_id, name, description, is_active, created_at, updated_at) VALUES
(1, 1, 1, '週一早班哈達瑜伽', '每週一 09:00-10:00', true, NOW(), NOW()),
(2, 1, 2, '週三熱瑜伽', '每週三 18:00-19:00', true, NOW(), NOW()),
(3, 1, 3, '週五晚間肌力訓練', '每週五 19:00-19:45', true, NOW(), NOW()),
(4, 1, 4, '週六芭蕾課程', '每週六 14:00-15:30', true, NOW(), NOW());

-- =====================================================
-- 測試帳號資訊
-- =====================================================
-- 管理員帳號:
--   Email: admin@timeledger.com / staff@timeledger.com
--   密碼: admin123
--
-- 老師帳號 (LINE User ID):
--   李小美: U000000000000000000000001
--   陳大文: U000000000000000000000002
--   王小花: U000000000000000000000003
--
-- 快速登入測試:
--   請使用本機測試資料庫，並在登入頁面點擊快速登入按鈕
-- =====================================================
