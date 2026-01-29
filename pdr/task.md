# Refactoring & Optimization Plan

## Phase 1: Immediate Fixes (Performance & Security)
- [ ] **DB Index Optimization**: Add generated columns for `schedule_rules` `start_date`/`end_date` and index them <!-- id: 1 -->
- [ ] **Remove Auth Backdoor**: Remove or flag `mock-` token logic in `auth.go` <!-- id: 2 -->

## Phase 2: Structural Refactoring
- [ ] **Service Layer Boundary**: Extract scheduling logic from `TeacherController` to `SchedulingService` <!-- id: 3 -->
- [ ] **Repository Pattern**: Enforce repository usage, remove direct `MySQL.WDB` access in controllers <!-- id: 4 -->

## Phase 3: Infrastructure
- [ ] **Structured Logging**: Replace `fmt.Println` with `Zap` or `Slog` <!-- id: 5 -->
- [ ] **Async Job Queue**: Replace custom Redis loop with `Asynq` <!-- id: 6 -->

## Phase 4: Frontend Polish
- [ ] **Virtual Scrolling**: Implement `vue-virtual-scroller` for `ScheduleGrid` <!-- id: 7 -->
- [ ] **PWA Support**: Configure `@vite-pwa/nuxt` <!-- id: 8 -->

## Phase 5: Deep Dive UX/DX
- [ ] **Middleware Context**: Create `ContextHelper` to unwrap gin context type-safely <!-- id: 9 -->
- [ ] **Enhanced Export**: Implement ICS subscription and backend image generation <!-- id: 10 -->

## Phase 6: Code Consistency (DRY)
- [ ] **Generic Repository**: Implement `UsingGenericRepo[T]` for standard CRUD <!-- id: 11 -->
- [ ] **Frontend API Refactor**: Centralize header/token logic in `useApi.ts` <!-- id: 12 -->
