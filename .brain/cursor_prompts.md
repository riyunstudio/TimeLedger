# Cursor Implementation Prompts - Term Management & Occupancy View

> [!IMPORTANT]
> **Before starting any prompt, please READ the full specifications in `.brain/implementation_plan.md` and `.brain/task.md`** to ensure alignment with the overall design and progress tracking.

---

## Stage 1: Backend Foundation

### Prompt 1.1: Create CenterTerm Model & Migration
> Please read `.brain/implementation_plan.md` and `.brain/task.md` for context.
> Then, create a new backend model `CenterTerm` in `app/models/center_term.go`. 
> - Fields: `ID` (uint), `CenterID` (uint), `Name` (string), `StartDate` (time.Time), `EndDate` (time.Time), and GORM timestamps.
> - Ensure it has the necessary GORM tags and JSON tags.
> - After creating the model, add it to the auto-migration list in `database/migration.go` (or equivalent migration entry point) to ensure the `center_terms` table is created.
> - Mark the task as done in `.brain/task.md` once complete.

### Prompt 1.2: Term CRUD Controller
> Please read `.brain/implementation_plan.md` for API details.
> Implement `AdminTermController` in `app/controllers/admin_term_controller.go` to handle CRUD for `CenterTerm`.
> - Endpoints: `GET /admin/terms`, `POST /admin/terms`, `PUT /admin/terms/:id`, `DELETE /admin/terms/:id`.
> - Responses should follow the project's standard response format.
> - Use `useCenterId` logic to ensure terms are scoped to the administrator's center.
> - Update `.brain/task.md`.

### Prompt 1.3: Occupancy Aggregation & Copy Logic
> Please refer to the "Aggregation Endpoint" and "Batch Copy Endpoint" sections in `.brain/implementation_plan.md`.
> Add two specialized endpoints to `AdminTermController`:
> 1. `GET /admin/occupancy/rules`: Aggregates `ScheduleRule` records for a given `teacher_id` or `room_id` that overlap with a date range. Return rules grouped by `DayOfWeek`.
> 2. `POST /admin/terms/copy-rules`: Accepts `source_term_id`, `target_term_id`, and `rule_ids[]`. 
>    - Logic: Clone the specified rules, setting their `StartDate` and `EndDate` to match the target term's range.
>    - Save the new rules and ensure `ScheduleExpansionService` is called to generate the actual course instances for the new term.
> - Update `.brain/task.md`.

---

## Stage 2: Terms Management (Frontend)

### Prompt 2.1: TermsTab Component
> Refer to Stage 2 in `.brain/task.md`.
> Create `frontend/components/Admin/TermsTab.vue` to manage academic terms.
> - Use a list/table format to display terms (Name, Start Date, End Date).
> - Implement Add/Edit/Delete functionality using a modal.
> - Use `useApi` for requests and `useNotification` for feedback.

### Prompt 2.2: Integrate into Resources
> Modify `frontend/pages/admin/resources.vue` to include the "學期期間 (Terms)" tab.
> - Import and register `TermsTab.vue`.
> - Add a new tab entry in the template and manage the active tab state.
> - Update `.brain/task.md`.

---

## Stage 3: Occupancy View & Copy Wizard

### Prompt 3.1: Resource Occupancy Page
> Please read the "Occupancy Visualization" section in `.brain/implementation_plan.md`.
> Create `frontend/pages/admin/resource-occupancy.vue`.
> - Add filters for `Term` (dropdown), `Type` (Teacher/Room), and a search bar for the specific resource.
> - Implement a weekly grid view (Monday to Sunday) that fetches rules from `/admin/occupancy/rules`.
> - Occupied slots should show the Course Name, Teacher/Room, and Time.
> - Add a "Copy Rules" button that opens a wizard.

### Prompt 3.2: Term Copy Wizard
> Build the "Term Copy Wizard" as specified in `.brain/implementation_plan.md`.
> Create `frontend/components/Admin/TermCopyModal.vue`.
> - Step 1: Select Source Term.
> - Step 2: Display list of rules in the source term with checkboxes (default all selected).
> - Step 3: Select Target Term.
> - Step 4: Show a summary of how many rules will be copied and any date adjustments.
> - On confirm, call `POST /admin/terms/copy-rules`.
> - Update `.brain/task.md`.

### Prompt 3.3: UX Polish (Drag & Drop)
> Refer to the "UX Optimizations" section in `.brain/implementation_plan.md`.
> Add Drag-and-Drop functionality to the weekly grid in `resource-occupancy.vue`.
> - Use a library like `vuedraggable` or native HTML5 DND.
> - Dragging a rule to a different time/day should update the `day_of_week`, `start_time`, and `end_time` of the underlying `ScheduleRule`.
> - Highlight overlaps/conflicts in real-time by checking for shared time slots in the local state.
