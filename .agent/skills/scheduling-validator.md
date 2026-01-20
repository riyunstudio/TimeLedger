---
name: scheduling-validator
description: Enforces Test-Driven Development (TDD) for the core scheduling engine, including expansion rules and conflict detection logic.
---

# Scheduling Engine Validation Skill

The scheduling engine is the heart of TimeLedger. This skill ensures that complex logic (Buffer, Overlap, Recurrence) is 100% verified by tests.

## Core Directives

1. **Test-First Requirement**:
   - Before writing any code in `backend/internal/service/schedule/`, you MUST write unit tests in `*_test.go`.
   - Tests must cover at least:
     - **Hard Overlap**: Two sessions in the same room/teacher at the same time.
     - **Null Teacher ID**: Ensuring "Draft" sessions only check room conflicts.
     - **Buffer Boundary**: Exactly 10m buffer (success) vs 9m buffer (fail/warn).
     - **Cross-day Expansion**: Rules that repeat across midnight.

2. **Logic Alignment**:
   - Reference `pdr/功能業務邏輯.md` chapter numbers in code comments.
   - Use the `ValidationResult` DTO defined in `pdr/API.md`.

3. **Concurrency Guard**:
   - When implementation involves DB writes, ensure `FOR UPDATE` locks are tested using parallel test runs or simulated race conditions.

## Verification Checklist
- [ ] Unit tests pass for the target logic.
- [ ] Edge cases (leap year, daylight saving if applicable) are considered.
- [ ] Complexity of validation function is kept low (< 15 cyclomatic complexity).
