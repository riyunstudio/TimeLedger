---
name: auth-adapter-guard
description: Manages the abstraction layer between Login and real LINE Login (LIFF), ensuring seamless switching without breaking business logic.
---

# Auth Adapter & Security Guard Skill

This skill enforces the authentication layer while maintaining a path to the secure LINE Native authentication.

## Core Directives

1. **Interface Enforcement**:
   - All auth operations must go through `AuthService` interface.
   - NEVER call `liff.*` directly in generic UI components; use `useAuth()` composable.
   - NEVER query the `users` table directly for identity; use a standard `Claims` object from the middleware.

2. **Authentication**:
   - Admin login uses email/password authentication via `AuthService.AdminLogin()`.
   - Teacher login uses LINE Login (LIFF) authentication via `AuthService.TeacherLineLogin()`.
   - The `AuthService` validates credentials against the database and generates JWT tokens.

3. **Privacy Enforcement**:
   - Check all `GET /admin/centers/{id}/*` endpoints to ensure they enforce the **Security Wall** defined in `功能業務邏輯.md`.
   - Prevent cross-center data leakage (especially for teacher internal notes).

## Verification Checklist
- [ ] `AuthService` interface is implemented.
- [ ] `useAuth()` composable handles LINE authentication.
- [ ] API middleware correctly scopes center data based on the JWT `center_id` claim.
