# Task: Academic Terms & Resource Occupancy View

- [ ] Planning & Design
    - [x] Define `CenterTerm` database schema
    - [x] Design "Term Weekly View" logic
    - [x] Design "Batch Copy Rule" logic
    - [x] Create detailed implementation plan
    - [x] Generate Cursor AI Prompts

- [ ] Stage 1: Backend Infrastructure
    - [x] Create `CenterTerm` model (`app/models/center_term.go`)
    - [x] Add GORM Migration for `center_terms` table
    - [x] Implement `AdminTermController` (CRUD)
    - [x] Implement API: `GET /admin/occupancy/rules` (Aggregation)
    - [x] Implement API: `POST /admin/terms/copy-rules` (Cloning)

- [ ] Stage 2: Term Management (Frontend)
    - [x] Create `frontend/components/Admin/TermsTab.vue`
    - [x] Integrate Terms tab into `frontend/pages/admin/resources.vue`
    - [ ] Create `TermEditModal.vue` for CRUD operations

- [ ] Stage 3: Resource Occupancy View (Frontend)
    - [ ] Create `frontend/pages/admin/resource-occupancy.vue`
    - [ ] Build `AdminOccupancyGrid.vue` (The weekly rule-based calendar)
    - [ ] Implement Drag & Drop logic for rule resizing/moving

- [ ] Stage 4: Term Copy Wizard
    - [x] Create `TermCopyModal.vue`
    - [x] Integrate into TermsTab.vue
    - [x] Implement conflict preview and confirmation logic
    - [x] Final end-to-end verification
