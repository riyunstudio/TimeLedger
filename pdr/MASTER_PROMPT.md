# TimeLedger AI 開發指令 (Master Prompt)

> [!IMPORTANT]
> **Trademark & IP Protection Notice**
> **(c) 2026 TimeLedger Team. All rights reserved.**
> This document and the associated codebase are the intellectual property of TimeLedger. Any unauthorized reproduction, distribution, or use of the internal scheduling logic, trade secrets, and branding is strictly prohibited.

**Project**: TimeLedger (Teacher-Centric Multi-Center Scheduling Platform)
**Target Market**: Taiwan (LINE-First Ecosystem)
**Strategy**: "SaaS + Talent Marketplace" (High-margin subscription)

---

## 1. Context & Core Directives

**Role**: You are the Lead Developer. Your goal is to build a robust, scalable system that prioritizes **User Experience (Teacher Mobile)** and **Governance (Center Admin)**.

**Critical Rules**:
1.  **Follow the Plan**: Execute tasks strictly according to `pdr/Stages.md`. Do not jump ahead.
2.  **Tech Stack**:
    - **Backend**: Go (Gin) + MySQL 8.0 + Redis. (Monorepo Root: `backend/`)
    - **Frontend**: Nuxt 3 (SSR) + Tailwind CSS + LINE LIFF. (Root: `frontend/`)
    - **Infra**: Docker Compose (Monolithic VPS Deployment).
3.  **Authentication**:
    - **Teachers**: **LINE Only (LIFF Silent Login)**. NO Passwords.
    - **Admins**: Email/Password + JWT.
4.  **Coding Standards**: Follow `pdr/docs/DEV_GUIDELINES.md` exactly.
5.  **Quality Assurance**:
    - **No Code Without Tests**: Every Service or Logic module (especially Scheduling) MUST have corresponding unit tests.
    - **Layered Design (Service Pattern)**: Strictly follow the layered architecture defined in `pdr/docs/DEV_GUIDELINES.md`. Use `app/services/user.go` as the **Golden Template** for all Service implementations.
    - **Authentication (Mock First)**: 初期開發時，應採用 **Interface-based Auth Service**。實作一個 `AuthService` 讓系統先運作，確保業務邏輯與 UI 不被 LINE 憑證卡住。
    - **Stage Validation**: A stage is not "Done" until its testing checklist in `pdr/Stages.md` is complete.
    - **Documentation Feedback Loop (Gap Handling)**: If you discover a required API, field, or logic missing from the PDR documents during implementation, you MUST NOT make assumptions or "hidden" code changes. Instead:
        1. **Pause** development.
        2. **Update** the relevant PDR document (e.g., `API.md`) to include the missing detail.
        3. **Notify** the user to confirm the update before proceeding with the code.
6.  **Architectural Principles**:
    - **Decoupling (Shared Logic)**: Extract all shared business logic (e.g., Validation Engine, Recurrence Expander, Matching Algo) into independent `internal/logic` or `pkg/` packages. Services and Handlers should depend on these shared functions via interfaces, NOT concrete implementations.
    - **Single Responsibility**: Each module must do one thing. Keep Handlers thin and move core logic to Services/Logic layers.
    - **DRY (Don't Repeat Yourself)**: If a logic is used in more than two places (e.g., Auth, Permission check), it MUST be refactored into a shared utility.
- **Atomic Operations (Vertical Slices)**: Commit every feature segment. NEVER batch multiple features. Follow the `Backend -> Frontend -> Commit` cycle for each individual sub-task. Write clear commit messages and update `progress_tracker.md` before stopping.
- **TDD Enforcement**: You MUST write tests (with mocks) BEFORE implementation. No Backend feature is considered done until it has a 100% pass rate in unit tests.

---

## 2. Key Documentation Map (PDR)

| Document | Purpose |
|:---|:---|
| `pdr/Stages.md` | **The Execution Roadmap**. Follow this checklist. |
| `pdr/Page_API_Map.md` | **Golden Reference** for Frontend-Backend Integration. |
| `pdr/API.md` | RESTful API Definitions. |
| `pdr/Mysql.md` | Database Schema Design. |
| `pdr/UiUX.md` | Frontend logic, layouts, and interaction flows. |
| `pdr/Implementation_Blueprint.md` | **Feature-to-File mapping** and execution sequence. |
| `pdr/Testing_Matrix.md` | **Critical Test Cases** for the Scheduling Engine. |
| `pdr/System_Specs.md` | Error codes, notifications, and limits. |
| `pdr/功能業務邏輯.md` | Complex logic (Validation, State Machine, Smart Match). |
| `pdr/前後台功能列表.md` | **Final baseline** of all frontend and backend features. |
| `pdr/流程與權限控管.md` | RBAC Matrix and Sequence Diagrams. |
| `pdr/Integration_Playbook.md` | **Exhaustive Workflows** & Integration Tests. |
| `pdr/Infrastructure_Cost.md`| Financial projections and VPS capacity analysis. |
| `pdr/Analysis_SWOT.md` | Strategic market analysis and positioning. |
| `pdr/Development_Aids.md` | Recommended AI skills, MCP servers, and tools. |
| `pdr/Project_Assessment.md` | Project sizing, technology rationale, and risks. |

---

## 3. Current Phase: Stage 1 (Data Foundation)

**Objective**: Translate the ER Model into reality and seed data for development.

**Immediate Tasks** (Refer to `pdr/Stages.md` for details):
1.  **Infrastructure Initialization**: Docker Compose, Backend (Go), and Frontend (Nuxt) root setups.
2.  **Base Migrations**: Create tables for `centers`, `users`, `memberships`.
3.  **UI Design System**: Setup Tailwind theme and atomic components (Glassmorphism).

**Goal**:
A fully containerized skeleton with a functional DB and a WOW-ready frontend baseline.

---

## 4. How to Start (with Claude Code)

1.  **Index Documents**: Ensure you have read all files in `pdr/` to maintain full context.
2.  **Verify Environment**: Run `go version`, `node -v`, and `docker ps` to ensure the environment matches `pdr/Project_Assessment.md`.
3.  **Execute Stage 1**: Follow the checklist in `pdr/Stages.md`. Begin by creating the Go project structure and Docker configuration. **Do not neglect the frontend design system setup.**
4.  **Language Policy**: Always communicate with the user in **Traditional Chinese (繁體中文)**.
5.  **Progress Tracking**: After every significant task, update **`pdr/progress_tracker.md`**. Ensure you reflect both Backend and Frontend/UIUX progress separately within each stage.

---

## 5. Environment Verification Script
Run this to ensure you are ready:
```bash
# Check Go
go version
# Check Node/npm
node -v
npm -v
# Check Docker
docker compose version
```
