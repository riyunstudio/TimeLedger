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

## Phase 7: Architectural Refinement
- [ ] **Rich Domain Models**: Move business rules (validation, status checks) from Services to Model structs <!-- id: 13 -->
- [ ] **DTO Implementation**: Introduce Request/Resource structs to decouple API from DB Schema <!-- id: 14 -->
- [ ] **Unified Error Handling**: Implement `AppError` and central Middleware for JSON error responses <!-- id: 15 -->
- [ ] **Interface-based DI**: Refactor Controllers to depend on Service Interfaces instead of concrete types <!-- id: 16 -->

## Phase 8: Performance & Scalability
- [ ] **N+1 Query Resolution**: Refactor `ExpandRules` to use batch fetching and in-memory mapping <!-- id: 17 -->
- [ ] **Redis Schedule Cache**: Implement caching for computed teacher schedules with smart invalidation <!-- id: 18 -->
- [ ] **Event-Driven Bus**: Implement internal event dispatcher for decoupling notifications/audits <!-- id: 19 -->
- [ ] **Fat Controller Splitting**: Decompose `TeacherController` into domain-specific sub-controllers <!-- id: 20 -->
