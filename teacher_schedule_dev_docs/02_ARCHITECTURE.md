# 02_ARCHITECTURE.md

## Deployment
- Single server: Go API + MySQL

## Key decisions
- Schedule view computed on query (no materialization)
- Service layer enforces gating + validate
- center scoping enforced in repository/service
