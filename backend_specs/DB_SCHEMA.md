# DB_SCHEMA.md — TimeLedger（MySQL 8.x）
日期：2026-01-18
charset：utf8mb4
engine：InnoDB

> 所有表皆含：id、created_at、updated_at（DATETIME(3)）
> 所有欄位必須 COMMENT
> center-scoped 表需有 center_id 並建立索引

---

## 1. 必要表
- teachers
- teacher_personal_sessions
- centers
- admin_users
- teacher_center_memberships
- offerings
- schedule_rules
- schedule_exceptions
- center_plans
- audit_logs
- event_logs

---

## 2. 索引與唯一性（重點）
- memberships：unique(center_id, teacher_id)
- rules：idx(center_id, teacher_id, weekday) + idx(center_id, offering_id)
- personal_sessions：idx(teacher_id, start_at, end_at)
- exceptions：需要「同日有效例外唯一性」（見 SQL 的備註）
- audit_logs/event_logs：idx(center_id, created_at) + idx(actor/teacher, created_at)

---

## 3. 完整建表 SQL
請見同目錄：`DB_SCHEMA_MYSQL.sql`
