# LINE Bot Autonomous Registration Flow

This plan outlines the steps to implement a public registration flow for teachers via a LINE Bot Rich Menu. This allows new teachers to join the platform without a center-specific invitation, adding themselves to the "Talent Pool".

## User Review Required

> [!IMPORTANT]
> The Rich Menu setup requires an **active LINE Channel Access Token**. I will provide the configuration JSON and a `curl` command for you to run, or I can create a helper script if preferred.

## Proposed Changes

### Backend Enhancements

#### `TeacherController` [NEW] (app/controllers/teacher.go)
- [NEW] Add `PublicRegister` method:
    - Validates `LineUserID`, `Name`, and `Email`.
    - Checks if the user already exists.
    - Calls `TeacherService.RegisterPublic`.

#### `TeacherService` (app/services/teacher.go)
- [NEW] Add `RegisterPublic` method:
    - Creates a new `Teacher` record with `IsOpenToHiring = true`.
    - Generates a JWT token so the user is logged in immediately after registration.

#### `route.go` (app/servers/route.go)
- [MODIFY] Register the new endpoint:
```go
{http.MethodPost, "/api/v1/teacher/public/register", s.action.teacher.PublicRegister, []gin.HandlerFunc{}},
```

---

### Frontend Enhancements

#### `Public Registration Page` [NEW] (frontend/pages/teacher/register.vue)
- [NEW] Create a registration form:
    - **LIFF Integration**: Fetch `LineUserID` from LIFF context.
    - **Form Fields**: Name and Email.
    - **Submission**: Calls the new public register endpoint.
    - **Redirect**: On success, saves the token and redirects to `/teacher/dashboard`.

---

### LINE Bot Configuration

#### `Rich Menu Setup`
- [NEW] Provide a JSON configuration for a standard Rich Menu:
    - **Action**: Open URI (`${liffUrl}/teacher/register`).
    - **Label**: "✨ 註冊老師個人檔案" (Register Teacher Profile).
- [NEW] Provide a helper script `scripts/setup_rich_menu.sh` to upload this configuration to the LINE Messaging API.

## Verification Plan

### Automated Tests
- Run `go build ./...` to verify backend changes.
- Add a test case in `app/services/teacher_test.go` (if exists) for public registration.

### Manual Verification
1. Run the rich menu setup script.
2. Open the LINE Bot on a mobile device as a new user.
3. Click the "Register" button in the Rich Menu.
4. Complete the registration form.
5. Verify the teacher is created in the database and logged into the dashboard.
## 7. General Web Registration (/register)

This extension provides a dedicated `/register` page on the main website that uses standard LINE Login (OAuth2) instead of LIFF, making it accessible from any browser.

### Proposed Changes

#### `General Registration Page` [NEW] (frontend/pages/register.vue)
- [NEW] A clean, branded landing page for registration.
- **LINE Login Flow**:
    - "Register with LINE" button.
    - Redirects to LINE Login URL (configured in env).
    - Handles the callback to get `LineUserID`.
- **Registration Form**:
    - Reuses the logic from `teacher/register.vue` once the LINE identity is confirmed.

#### `Logic Synchronization`
- Ensure both `/teacher/register` (LIFF) and `/register` (Web) use the same backend `PublicRegister` API.

## 8. General Center Invitation Link (Multiple Teachers)

This allows a center to have a persistent invitation link that can be shared with multiple teachers simultaneously.

### Proposed Changes

#### Backend Enhancements

##### `CenterInvitation` Model (app/models/center_invitation.go)
- [MODIFY] Add `InvitationTypeGeneral` constant.
- [MODIFY] Update logic to allow `Email` to be empty for general invitations.

##### `TeacherService` (app/services/teacher.go)
- [MODIFY] `AcceptInvitationByLink`:
    - If `InviteType == GENERAL`, do NOT mark the invitation as `ACCEPTED` or set `RespondedAt` (as the link is reusable).
    - Skip the `Email` mismatch check if it's a general link.
- [NEW] `GenerateGeneralInvitationLink`: Creates a persistent link with a specific role.

#### Frontend Enhancements

##### `Admin Invitations Page` (frontend/pages/admin/invitations.vue)
- [MODIFY] Add a section to manage the "General Center Link".
- [MODIFY] Actions to "Enable/Disable" (Toggle status between `PENDING` and `EXPIRED`) and "Regenerate" (Create new link).

## Verification Plan

### Automated Tests
- Test accepting a general invitation link with multiple different `LineUserID`s.

### Manual Verification
1. As an admin, generate a "General Invitation Link".
2. Share the same link with two different teachers.
3. Both teachers should be able to join the center using that same link.
4. Disable the link and verify teachers can no longer join.
