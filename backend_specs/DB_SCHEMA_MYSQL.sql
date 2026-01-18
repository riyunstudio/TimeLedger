-- DB_SCHEMA_MYSQL.sql — TimeLedger (MySQL 8.x)
-- Date: 2026-01-18
-- charset: utf8mb4, engine: InnoDB

SET NAMES utf8mb4;
SET time_zone = '+00:00';

CREATE TABLE IF NOT EXISTS teachers (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  display_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '顯示名稱',
  avatar_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '頭像 URL',
  line_user_id VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'LINE userId（可選，未整合 LINE 先留欄位）',
  last_login_at DATETIME(3) NULL COMMENT '最後登入時間（可選）',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  UNIQUE KEY uk_teachers_line_user_id (line_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='老師（全域）';

CREATE TABLE IF NOT EXISTS teacher_personal_sessions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  teacher_id BIGINT UNSIGNED NOT NULL COMMENT '老師ID(teachers.id)',
  start_at DATETIME(3) NOT NULL COMMENT '開始時間',
  end_at DATETIME(3) NOT NULL COMMENT '結束時間',
  title VARCHAR(128) NOT NULL COMMENT '標題',
  note VARCHAR(255) NOT NULL DEFAULT '' COMMENT '備註',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  KEY idx_personal_teacher_time (teacher_id, start_at, end_at),
  CONSTRAINT fk_personal_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='老師個人行程';

CREATE TABLE IF NOT EXISTS centers (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  name VARCHAR(128) NOT NULL COMMENT '中心名稱',
  code VARCHAR(32) NOT NULL COMMENT '中心代碼（對外/邀請用）',
  status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE' COMMENT '狀態 ACTIVE/INACTIVE',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  UNIQUE KEY uk_centers_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='中心（租戶）';

CREATE TABLE IF NOT EXISTS admin_users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NOT NULL COMMENT '中心ID(centers.id)',
  account VARCHAR(64) NOT NULL COMMENT '帳號',
  password_hash VARCHAR(255) NOT NULL COMMENT '密碼雜湊（MVP 可不用）',
  token VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'MVP stub token',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  UNIQUE KEY uk_admin_center_account (center_id, account),
  KEY idx_admin_token (token),
  CONSTRAINT fk_admin_center FOREIGN KEY (center_id) REFERENCES centers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='中心管理員';

CREATE TABLE IF NOT EXISTS teacher_center_memberships (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NOT NULL COMMENT '中心ID',
  teacher_id BIGINT UNSIGNED NOT NULL COMMENT '老師ID',
  status VARCHAR(16) NOT NULL COMMENT 'INVITED/ACTIVE',
  invited_at DATETIME(3) NULL COMMENT '邀請時間',
  activated_at DATETIME(3) NULL COMMENT '啟用時間',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  UNIQUE KEY uk_membership_center_teacher (center_id, teacher_id),
  KEY idx_membership_center_status (center_id, status),
  KEY idx_membership_teacher (teacher_id),
  CONSTRAINT fk_membership_center FOREIGN KEY (center_id) REFERENCES centers(id),
  CONSTRAINT fk_membership_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='老師加入中心關係';

CREATE TABLE IF NOT EXISTS offerings (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NOT NULL COMMENT '中心ID',
  name VARCHAR(128) NOT NULL COMMENT 'Offering 名稱',
  duration_min INT NOT NULL COMMENT '時長(分鐘)',
  status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE' COMMENT 'ACTIVE/INACTIVE',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  KEY idx_offering_center (center_id),
  CONSTRAINT fk_offering_center FOREIGN KEY (center_id) REFERENCES centers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='中心課程單位';

CREATE TABLE IF NOT EXISTS schedule_rules (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NOT NULL COMMENT '中心ID',
  teacher_id BIGINT UNSIGNED NOT NULL COMMENT '老師ID',
  offering_id BIGINT UNSIGNED NOT NULL COMMENT 'Offering ID',
  weekday TINYINT NOT NULL COMMENT '星期(1=Mon..7=Sun)',
  start_time CHAR(5) NOT NULL COMMENT '開始時間HH:mm',
  end_time CHAR(5) NOT NULL COMMENT '結束時間HH:mm',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  KEY idx_rule_center_teacher_weekday (center_id, teacher_id, weekday),
  KEY idx_rule_center_offering (center_id, offering_id),
  CONSTRAINT fk_rule_center FOREIGN KEY (center_id) REFERENCES centers(id),
  CONSTRAINT fk_rule_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id),
  CONSTRAINT fk_rule_offering FOREIGN KEY (offering_id) REFERENCES offerings(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='固定排課規則';

CREATE TABLE IF NOT EXISTS schedule_exceptions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NOT NULL COMMENT '中心ID',
  rule_id BIGINT UNSIGNED NOT NULL COMMENT '對應 rule',
  teacher_id BIGINT UNSIGNED NOT NULL COMMENT '老師ID（冗餘，便於查詢）',
  type VARCHAR(16) NOT NULL COMMENT 'CANCEL/RESCHEDULE',
  status VARCHAR(16) NOT NULL COMMENT 'PENDING/APPROVED/REJECTED',
  original_date DATE NOT NULL COMMENT '原日期',
  new_start_at DATETIME(3) NULL COMMENT '新開始時間（RESCHEDULE）',
  new_end_at DATETIME(3) NULL COMMENT '新結束時間（RESCHEDULE）',
  reason VARCHAR(255) NOT NULL DEFAULT '' COMMENT '原因',
  reviewed_by BIGINT UNSIGNED NULL COMMENT '審核者 admin_users.id',
  reviewed_at DATETIME(3) NULL COMMENT '審核時間',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  UNIQUE KEY uk_exc_active (center_id, rule_id, original_date, status),
  KEY idx_exc_center_status (center_id, status),
  KEY idx_exc_teacher_date (teacher_id, original_date),
  CONSTRAINT fk_exc_center FOREIGN KEY (center_id) REFERENCES centers(id),
  CONSTRAINT fk_exc_rule FOREIGN KEY (rule_id) REFERENCES schedule_rules(id),
  CONSTRAINT fk_exc_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='例外（停課/改期）';

-- NOTE:
-- uk_exc_active 以 (status) 做 unique 會限制 REJECTED 也只能一筆。
-- MVP 可接受；若要允許多次 rejected，需改為：generated column is_active = status IN ('PENDING','APPROVED')
-- 再 unique(center_id, rule_id, original_date, is_active)。

CREATE TABLE IF NOT EXISTS center_plans (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NOT NULL COMMENT '中心ID',
  plan VARCHAR(16) NOT NULL COMMENT 'starter/pro/team',
  status VARCHAR(16) NOT NULL COMMENT 'active/expired',
  max_teachers INT NULL COMMENT '最多老師數（team 可為 NULL 表示無限）',
  trial_start_at DATETIME(3) NOT NULL COMMENT '試用開始',
  trial_end_at DATETIME(3) NOT NULL COMMENT '試用結束',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  UNIQUE KEY uk_plan_center (center_id),
  KEY idx_plan_status_end (status, trial_end_at),
  CONSTRAINT fk_plan_center FOREIGN KEY (center_id) REFERENCES centers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='中心方案/試用（不含金流）';

CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NULL COMMENT '中心ID（teacher 個人行程可為 NULL）',
  actor_type VARCHAR(16) NOT NULL COMMENT 'admin/teacher/system',
  actor_id BIGINT UNSIGNED NULL COMMENT '行為者ID',
  action VARCHAR(64) NOT NULL COMMENT '行為（create/update/delete/approve/reject）',
  entity_type VARCHAR(64) NOT NULL COMMENT '資源類型',
  entity_id BIGINT UNSIGNED NULL COMMENT '資源ID',
  payload_json JSON NULL COMMENT 'payload（可含 before/after）',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  KEY idx_audit_center_time (center_id, created_at),
  KEY idx_audit_actor_time (actor_type, actor_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作紀錄';

CREATE TABLE IF NOT EXISTS event_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK',
  center_id BIGINT UNSIGNED NULL COMMENT '中心ID',
  teacher_id BIGINT UNSIGNED NULL COMMENT '老師ID',
  event_name VARCHAR(64) NOT NULL COMMENT '事件名稱',
  event_props_json JSON NULL COMMENT '事件屬性',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '建立時間',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新時間',
  PRIMARY KEY (id),
  KEY idx_event_center_time (center_id, created_at),
  KEY idx_event_teacher_time (teacher_id, created_at),
  KEY idx_event_name_time (event_name, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件追蹤';
