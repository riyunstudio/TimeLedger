---
name: contract-sync
description: Ensures consistent naming and structure between PDR specifications (API.md, Mysql.md) and actual code implementation (Go Structs, TypeScript Types).
---

# Contract-First Synchronization Skill

Use this skill whenever you modify `pdr/API.md` or `pdr/Mysql.md` to ensure the codebase remains in lockstep with the specifications.

## Core Directives

1. **Backend Sync (Go)**:
   - Location: `backend/internal/models/` and `backend/internal/dto/`.
   - Action: Update GORM structs and JSON DTOs to match the updated fields, types, and nullability (especially `TeacherID` being nullable).
   - Validation: Run `go vet ./...` to ensure no field mismatches.

2. **Frontend Sync (TypeScript)**:
   - Location: `frontend/types/api.d.ts`.
   - Action: Regenerate or manually update TypeScript interfaces. Ensure naming (e.g., camelCase vs snake_case) matches the API spec's `json` tags.
   - Validation: Run `npm run typecheck` in the frontend directory.

3. **Documentation Consistency**:
   - Check if the `pdr/Page_API_Map.md` needs an update after an API change.

## Implementation Pattern
When this skill is triggered, you MUST provide a diff showing both the `.md` change and the corresponding `.go`/`.ts` changes in a single operation.
