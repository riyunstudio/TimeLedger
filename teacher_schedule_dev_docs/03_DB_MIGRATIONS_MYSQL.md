# 03_DB_MIGRATIONS_MYSQL.md

## MySQL settings
- engine InnoDB, charset utf8mb4
- datetime DATETIME(3)
- migrations folder with versioned .sql
- use schema_migrations table

## Required
- Apply backend_specs/DB_SCHEMA_MYSQL.sql
- Ensure indices exist for time-range queries

## Backup
- daily mysqldump or managed backup; keep 7-14 days
