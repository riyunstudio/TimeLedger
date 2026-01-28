-- =====================================================
-- 資料庫索引優化遷移腳本
-- 日期：2026年1月28日
-- 目的：提升熱門查詢的效能
-- =====================================================

USE timeledger;

-- =====================================================
-- 1. schedule_exceptions 表索引優化
-- =====================================================

-- 現有索引：idx_rule_date (rule_id + original_date)
-- 缺少：center_id + status（審核列表查詢）
-- 缺少：status（狀態篩選）

-- 新增：center_id + status 複合索引（用於取得某中心的待審核例外）
ALTER TABLE schedule_exceptions
ADD INDEX idx_exception_center_status (center_id, status);

-- 新增：status 單一索引（用於統計各狀態數量）
ALTER TABLE schedule_exceptions
ADD INDEX idx_exception_status (status);

-- 新增：center_id + status + created_at（用於列表排序）
ALTER TABLE schedule_exceptions
ADD INDEX idx_exception_list (center_id, status, created_at DESC);


-- =====================================================
-- 2. center_memberships 表索引優化
-- =====================================================

-- 現有索引：teacher_id, center_id
-- 缺少：center_id + teacher_id + status（活躍會籍查詢）

-- 新增：複合索引（用於查詢某中心某教師的活躍會籍）
ALTER TABLE center_memberships
ADD INDEX idx_membership_active (center_id, teacher_id, status);


-- =====================================================
-- 3. schedule_rules 表索引優化
-- =====================================================

-- 現有索引：
--   - idx_center_weekday_time (center_id + weekday + start_time)
--   - idx_teacher_time (teacher_id + start_time)
--   - idx_room_time (room_id)

-- 新增：effective_range 查詢索引（用於展開排課規則）
ALTER TABLE schedule_rules
ADD INDEX idx_rule_effective_start (effective_range (CAST(NULL AS JSON)));

-- 新增：center_id + is_active（用於取得某中心的所有有效規則）
ALTER TABLE schedule_rules
ADD INDEX idx_rule_active (center_id, is_active);


-- =====================================================
-- 4. personal_events 表索引優化
-- =====================================================

-- 現有索引：idx_teacher_time (teacher_id + start_at)
-- 缺少：teacher_id + start_at + end_at（日期範圍查詢）

-- 修改：將現有索引擴展為複合索引
ALTER TABLE personal_events
DROP INDEX idx_teacher_time;

ALTER TABLE personal_events
ADD INDEX idx_teacher_timerange (teacher_id, start_at, end_at);


-- =====================================================
-- 5. session_notes 表索引優化
-- =====================================================

-- 現有索引：rule_id, teacher_id
-- 缺少：teacher_id + session_date（教師某日筆記查詢）
-- 缺少：teacher_id + rule_id + session_date（教師某課程筆記查詢）

-- 新增：教師日期查詢索引
ALTER TABLE session_notes
ADD INDEX idx_note_teacher_date (teacher_id, session_date);

-- 新增：完整查詢索引
ALTER TABLE session_notes
ADD INDEX idx_note_full (teacher_id, rule_id, session_date);


-- =====================================================
-- 6. center_teacher_notes 表索引優化
-- =====================================================

-- 新增：center_id + teacher_id 索引
ALTER TABLE center_teacher_notes
ADD INDEX idx_note_center_teacher (center_id, teacher_id);


-- =====================================================
-- 7. center_invitations 表索引優化
-- =====================================================

-- 現有索引：status, expires_at
-- 缺少：center_id + status（某中心邀請列表）

-- 新增：邀請列表查詢索引
ALTER TABLE center_invitations
ADD INDEX idx_invitation_center_status (center_id, status);


-- =====================================================
-- 8. teacher_skills 表索引優化
-- =====================================================

-- 新增：category 索引（人才庫技能搜尋）
ALTER TABLE teacher_skills
ADD INDEX idx_skill_category (category);

-- 新增：level 索引
ALTER TABLE teacher_skills
ADD INDEX idx_skill_level (level);


-- =====================================================
-- 9. notifications 表索引優化
-- =====================================================

-- 新增：user_id + is_read（通知列表）
ALTER TABLE notifications
ADD INDEX idx_notification_list (user_id, is_read, created_at DESC);


-- =====================================================
-- 10. notification_queues 表索引優化
-- =====================================================

-- 新增：status + retry_count（佇列處理）
ALTER TABLE notification_queues
ADD INDEX idx_queue_pending (status, retry_count, created_at);


-- =====================================================
-- 驗證索引是否建立成功
-- =====================================================
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX,
    NON_UNIQUE
FROM INFORMATION_SCHEMA.STATISTICS
WHERE TABLE_SCHEMA = 'timeledger'
    AND TABLE_NAME IN (
        'schedule_exceptions',
        'center_memberships',
        'schedule_rules',
        'personal_events',
        'session_notes',
        'center_teacher_notes',
        'center_invitations',
        'teacher_skills',
        'notifications',
        'notification_queues'
    )
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;


-- =====================================================
-- 完成
-- =====================================================
SELECT '索引優化完成！共優化 10 個資料表。' AS Status;
