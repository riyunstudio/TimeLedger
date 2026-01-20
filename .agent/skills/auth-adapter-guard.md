---
name: auth-adapter-guard
description: Manages the abstraction layer between Mock Login and real LINE Login (LIFF), ensuring seamless switching without breaking business logic.
---

# Auth Adapter & Security Guard Skill

This skill enforces the "Mock First" strategy while maintaining a path to the secure LINE Native authentication.

## Core Directives

1. **Interface Enforcement**:
   - All auth operations must go through `AuthService` interface.
   - NEVER call `liff.*` directly in generic UI components; use `useAuth()` composable.
   - NEVER query the `users` table directly for identity; use a standard `Claims` object from the middleware.

2. **Mock Mode Management**:
   - When `APP_ENV=development`, use `MockAuthService`.
   - Provide a persistent test token that identifies as "Admin (Owner)" and "Teacher" for local testing.
   - Ensure `MockAuthService` generates a valid `IdentInfo` that mimics real LINE payloads.

3. **Privacy Enforcement**:
   - Check all `GET /admin/centers/{id}/*` endpoints to ensure they enforce the **Security Wall** defined in `功能業務邏輯.md`.
   - Prevent cross-center data leakage (especially for teacher internal notes).

## Verification Checklist
- [ ] `AuthService` interface is implemented.
- [ ] `useAuth()` composable handles both LIFF and Mock states.
- [ ] API middleware correctly scopes center data based on the JWT `center_id` claim.
