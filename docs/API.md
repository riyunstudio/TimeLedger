# TimeLedger API Documentation

> Version: 1.0.0
> Base URL: `http://your-server:8080/api/v1`

---

## Table of Contents

- [Authentication](#authentication)
- [Auth Endpoints](#auth-endpoints)
- [Teacher Endpoints](#teacher-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Scheduling Endpoints](#scheduling-endpoints)
- [Smart Matching Endpoints](#smart-matching-endpoints)
- [Notification Endpoints](#notification-endpoints)
- [Export Endpoints](#export-endpoints)

---

## Authentication

### JWT Authentication

Most endpoints require a valid JWT token in the `Authorization` header:

```http
Authorization: Bearer <your-jwt-token>
```

### Roles

- **Admin**: Full access to all centers and resources
- **Teacher**: Access to own data and assigned centers
- **Center Admin**: Access to specific center resources

---

## Auth Endpoints

### Admin Login

Login with email and password (Admin only).

**Endpoint:** `POST /auth/admin/login`

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "your-password"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "name": "Admin User",
      "role": "ADMIN"
    }
  }
}
```

---

### Teacher Login (LINE)

Login with LINE OAuth (Teacher only).

**Endpoint:** `POST /auth/teacher/login`

**Request:**
```json
{
  "id_token": "LINE_ID_TOKEN",
  "access_token": "LINE_ACCESS_TOKEN"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "teacher": {
      "id": 1,
      "name": "Teacher Name",
      "line_user_id": "LINE_USER_ID"
    }
  }
}
```

---

### Refresh Token

Refresh access token using refresh token.

**Endpoint:** `POST /auth/refresh`

**Request:**
```json
{
  "refresh_token": "your-refresh-token"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Token refreshed",
  "data": {
    "token": "new-jwt-token",
    "refresh_token": "new-refresh-token"
  }
}
```

---

## Teacher Endpoints

### Get Profile

Get current teacher's profile.

**Endpoint:** `GET /teacher/me/profile`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "name": "Teacher Name",
    "email": "teacher@example.com",
    "avatar_url": "https://...",
    "bio": "Experienced piano teacher",
    "is_open_to_hiring": true,
    "city": "Taipei",
    "district": "Da'an",
    "public_contact_info": "phone: 0912345678",
    "skills": [
      {
        "skill_name": "Piano",
        "level": "Advanced"
      }
    ],
    "certificates": [],
    "personal_hashtags": []
  }
}
```

---

### Update Profile

Update teacher profile.

**Endpoint:** `PUT /teacher/me/profile`

**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "name": "Updated Name",
  "bio": "Updated bio",
  "is_open_to_hiring": false,
  "city": "Taipei",
  "district": "Xinyi",
  "public_contact_info": "email: updated@example.com"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Profile updated",
  "data": {
    "id": 1,
    "name": "Updated Name"
  }
}
```

---

### Get My Centers

Get list of centers where the teacher is a member.

**Endpoint:** `GET /teacher/me/centers`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Music Center",
      "status": "ACTIVE"
    }
  ]
}
```

---

### Upload Certificate

Upload a teaching certificate.

**Endpoint:** `POST /teacher/me/certificates`

**Headers:** `Authorization: Bearer <token>`

**Request:** `multipart/form-data`

- `certificate_name`: string (required)
- `issued_by`: string
- `issued_date`: string (YYYY-MM-DD)
- `file`: file (required)

**Response:**
```json
{
  "code": 200,
  "message": "Certificate uploaded",
  "data": {
    "id": 1,
    "certificate_name": "ABRSM Grade 8",
    "issued_by": "ABRSM",
    "file_url": "https://s3.example.com/certificates/..."
  }
}
```

---

## Admin Endpoints

### Centers

#### Create Center
**Endpoint:** `POST /admin/centers`
**Headers:** `Authorization: Bearer <token>`

#### Get Centers
**Endpoint:** `GET /admin/centers`

#### Get Center
**Endpoint:** `GET /admin/centers/:id`

#### Update Center
**Endpoint:** `PATCH /admin/centers/:id`

#### Delete Center
**Endpoint:** `DELETE /admin/centers/:id`

---

### Rooms

#### Create Room
**Endpoint:** `POST /admin/rooms`
**Headers:** `Authorization: Bearer <token>`

#### Get Rooms
**Endpoint:** `GET /admin/centers/:centerId/rooms`

#### Update Room
**Endpoint:** `PATCH /admin/rooms/:id`

#### Delete Room
**Endpoint:** `DELETE /admin/rooms/:id`

---

### Courses

#### Create Course
**Endpoint:** `POST /admin/courses`
**Headers:** `Authorization: Bearer <token>`

#### Get Courses
**Endpoint:** `GET /admin/centers/:centerId/courses`

---

### Offerings

#### Create Offering
**Endpoint:** `POST /admin/offerings`
**Headers:** `Authorization: Bearer <token>`

#### Get Offerings
**Endpoint:** `GET /admin/centers/:centerId/offerings`

---

## Scheduling Endpoints

### Check Overlap

Check for scheduling conflicts.

**Endpoint:** `POST /admin/scheduling/check-overlap`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "center_id": 1,
  "teacher_id": 1,
  "room_id": 1,
  "start_time": "2026-01-01T14:00:00Z",
  "end_time": "2026-01-01T15:00:00Z"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "valid": true,
    "conflicts": []
  }
}
```

---

### Create Schedule Rule

Create a weekly schedule rule.

**Endpoint:** `POST /admin/scheduling/rules`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "center_id": 1,
  "offering_id": 1,
  "teacher_id": 1,
  "room_id": 1,
  "weekday": 1,
  "start_time": "14:00:00",
  "end_time": "15:00:00",
  "effective_range_start_date": "2026-01-01",
  "effective_range_end_date": "2026-12-31"
}
```

---

### Get Schedule Rules

Get schedule rules for a center.

**Endpoint:** `GET /admin/scheduling/rules?center_id=:centerId`

---

### Create Exception

Create a schedule exception (leave/swap).

**Endpoint:** `POST /admin/scheduling/exceptions`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "center_id": 1,
  "rule_id": 1,
  "original_date": "2026-01-06",
  "type": "CANCEL",
  "reason": "Teacher is sick"
}
```

---

### Review Exception

Approve or reject a schedule exception.

**Endpoint:** `PATCH /admin/scheduling/exceptions/:id/review`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "status": "APPROVED",
  "review_note": "Approved by admin"
}
```

---

### Expand Rules

Expand weekly rules into actual dates.

**Endpoint:** `POST /admin/scheduling/expand`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "center_id": 1,
  "start_date": "2026-01-01",
  "end_date": "2026-01-31"
}
```

---

## Smart Matching Endpoints

### Find Matches

Find substitute teachers for a session.

**Endpoint:** `POST /admin/smart-matching/find-matches`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "center_id": 1,
  "room_id": 1,
  "start_time": "2026-01-06T14:00:00Z",
  "end_time": "2026-01-06T15:00:00Z",
  "required_skills": ["Piano"],
  "exclude_teacher_ids": [1]
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "matches": [
      {
        "teacher_id": 2,
        "teacher_name": "Substitute Teacher",
        "match_score": 85.5,
        "skill_match": 90,
        "rating": 4.5,
        "notes": "Good match"
      }
    ]
  }
}
```

---

### Search Talent

Search for teachers available for hire.

**Endpoint:** `POST /admin/talent/search`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "city": "Taipei",
  "district": "Da'an",
  "keyword": "piano",
  "skills": ["Piano"],
  "limit": 10,
  "offset": 0
}
```

---

## Notification Endpoints

### List Notifications

Get notifications for current user.

**Endpoint:** `GET /notifications`
**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `limit`: number of items (default: 20)
- `offset`: pagination offset

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "total": 10,
    "unread_count": 3,
    "notifications": [
      {
        "id": 1,
        "title": "New schedule assignment",
        "message": "You have been assigned to a new class...",
        "type": "SCHEDULE",
        "is_read": false,
        "created_at": "2026-01-21T10:00:00Z"
      }
    ]
  }
}
```

---

### Mark as Read

Mark a notification as read.

**Endpoint:** `POST /notifications/:id/read`
**Headers:** `Authorization: Bearer <token>`

---

### Mark All as Read

Mark all notifications as read.

**Endpoint:** `POST /notifications/read-all`
**Headers:** `Authorization: Bearer <token>`

---

### Get Unread Count

Get count of unread notifications.

**Endpoint:** `GET /notifications/unread-count`
**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "unread_count": 3
  }
}
```

---

## Export Endpoints

### Export Schedule CSV

Export schedule to CSV format.

**Endpoint:** `POST /admin/export/schedule/csv`
**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "center_id": 1,
  "start_date": "2026-01-01",
  "end_date": "2026-01-31"
}
```

**Response:** `text/csv` file download

---

### Export Schedule PDF

Export schedule to PDF format.

**Endpoint:** `POST /admin/export/schedule/pdf`
**Headers:** `Authorization: Bearer <token>`

**Request:** Same as CSV export

**Response:** `application/pdf` file download

---

### Export Teachers CSV

Export teachers list to CSV.

**Endpoint:** `GET /admin/export/centers/:centerId/teachers/csv`
**Headers:** `Authorization: Bearer <token>`

---

### Export Exceptions CSV

Export exceptions to CSV.

**Endpoint:** `GET /admin/export/centers/:centerId/exceptions/csv`
**Headers:** `Authorization: Bearer <token>`

---

## Admin Resource Endpoints

### Dashboard

#### Get Today Summary

Get today's schedule summary for the admin dashboard.

|**Endpoint:** `GET /admin/dashboard/today-summary`
|**Headers:** `Authorization: Bearer <token>`

|**Response:**
```json
{
  "code": "SUCCESS",
  "data": {
    "today_stats": {
      "completed": 12,
      "in_progress": 3,
      "upcoming": 8
    },
    "pending_exceptions": 5,
    "schedule_changes": [
      "明日 14:00 瑜伽課程教室異動",
      "週三 10:00 鋼琴個別課暫停"
    ]
  }
}
```

---

## Timetable Template Endpoints

### Templates CRUD

#### Get Templates

Get all templates for the center.

|**Endpoint:** `GET /admin/templates`
|**Headers:** `Authorization: Bearer <token>`

---

#### Create Template

Create a new timetable template.

|**Endpoint:** `POST /admin/templates`
|**Headers:** `Authorization: Bearer <token>`

|**Request:**
```json
{
  "name": "夏季課程模板",
  "description": "夏季密集課程排課模板"
}
```

---

#### Update Template

Update an existing template.

|**Endpoint:** `PUT /admin/templates/:templateId`
|**Headers:** `Authorization: Bearer <token>`

---

#### Delete Template

Delete a template.

|**Endpoint:** `DELETE /admin/templates/:templateId`
|**Headers:** `Authorization: Bearer <token>`

---

### Template Cells

#### Get Cells

Get all cells in a template.

|**Endpoint:** `GET /admin/templates/:templateId/cells`
|**Headers:** `Authorization: Bearer <token>`

---

#### Create Cells

Create cells in a template.

|**Endpoint:** `POST /admin/templates/:templateId/cells`
|**Headers:** `Authorization: Bearer <token>`

|**Request:**
```json
{
  "cells": [
    {
      "row": 1,
      "col": 1,
      "start_time": "09:00",
      "end_time": "10:00",
      "weekday": 1
    }
  ]
}
```

---

#### Delete Cell

Delete a cell from a template.

|**Endpoint:** `DELETE /admin/templates/cells/:cellId`
|**Headers:** `Authorization: Bearer <token>`

---

#### Apply Template

Apply a template to generate schedule rules.

|**Endpoint:** `POST /admin/templates/:templateId/apply`
|**Headers:** `Authorization: Bearer <token>`

|**Request:**
```json
{
  "offering_id": 1,
  "start_date": "2026-02-01",
  "end_date": "2026-12-31",
  "weekdays": [1, 3, 5],
  "override_buffer": false
}
```

|**Response (Success):**
```json
{
  "code": "SUCCESS",
  "message": "模板套用成功",
  "data": {
    "rules_created": 15,
    "sessions_generated": 45
  }
}
```

|**Response (Buffer Conflict):**
```json
{
  "code": 40003,
  "message": "套用模板會產生緩衝時間衝突，是否繼續？",
  "datas": {
    "conflicts": [...],
    "conflict_count": 2,
    "can_override": true
  }
}
```

---

## Admin Exception Endpoints

### Get All Exceptions

Get all exception requests with filtering support.

|**Endpoint:** `GET /admin/exceptions/all`
|**Headers:** `Authorization: Bearer <token>`

|**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|:---|:---|:---:|:---:|:---|
| `status` | string | No | - | Filter by status (PENDING, APPROVED, REJECTED, REVOKED) |
| `date_from` | string | No | - | Start date (ISO 8601) |
| `date_to` | string | No | - | End date (ISO 8601) |
| `teacher_id` | int | No | - | Filter by teacher |
| `room_id` | int | No | - | Filter by room |

|**Response:**
```json
{
  "code": 200,
  "datas": [
    {
      "id": 1,
      "type": "CANCEL",
      "status": "PENDING",
      "reason": "身體不適需要休息",
      "original_date": "2026-01-20",
      "offering_name": "瑜伽基礎",
      "rule": {
        "id": 10,
        "start_time": "09:00",
        "end_time": "10:00",
        "teacher": { "id": 1, "name": "張老師" },
        "room": { "id": 1, "name": "A教室" }
      }
    }
  ]
}
```

---

### Review Exception

Approve or reject an exception request.

|**Endpoint:** `POST /admin/scheduling/exceptions/:id/review`
|**Headers:** `Authorization: Bearer <token>`

|**Request:**
```json
{
  "action": "APPROVED",
  "reason": "已確認代課老師"
}
```

| Action Value | Description |
|:---|:---|
| `APPROVED` | Approve the exception request |
| `REJECTED` | Reject the exception request |

|**Response (Success):**
```json
{
  "code": 200,
  "message": "Exception reviewed successfully"
}
```

|**Response (Conflict Error):**
```json
{
  "code": 409,
  "message": "Cannot approve: schedule conflict detected",
  "datas": {
    "conflicts": [...]
  }
}
```

---

## Teacher Schedule Endpoints

### Get My Schedule

Get the teacher's schedule.

|**Endpoint:** `GET /teacher/me/schedule`
|**Headers:** `Authorization: Bearer <token>`

|**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|:---|:---|:---:|:---:|:---|
| `start_date` | string | No | Today | Start date (YYYY-MM-DD) |
| `end_date` | string | No | +7 days | End date (YYYY-MM-DD) |
| `view` | string | No | week | View mode (day, week, month) |

|**Response:**
```json
{
  "code": 200,
  "datas": [
    {
      "id": 1,
      "date": "2026-01-20",
      "start_time": "09:00",
      "end_time": "10:00",
      "offering_name": "瑜伽基礎",
      "center_name": "台北館",
      "room_name": "A教室",
      "is_holiday": false
    }
  ]
}
```

---

### Get My Personal Events

Get the teacher's personal events.

|**Endpoint:** `GET /teacher/me/personal-events`
|**Headers:** `Authorization: Bearer <token>`

|**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|:---|:---|:---:|:---:|:---|
| `start_date` | string | No | Today | Start date (YYYY-MM-DD) |
| `end_date` | string | No | +30 days | End date (YYYY-MM-DD) |

|**Response:**
```json
{
  "code": 200,
  "datas": [
    {
      "id": 1,
      "title": "醫院回診",
      "date": "2026-01-25",
      "start_time": "14:00",
      "end_time": "16:00",
      "is_recurring": false
    }
  ]
}
```

---

## Error Codes

| Code | Description |
|:---|:---|
| 200 | Success |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 422 | Validation Error |
| 500 | Internal Server Error |

---

## Common Errors

### Unauthorized
```json
{
  "code": 401,
  "message": "Unauthorized: Invalid or expired token"
}
```

### Validation Error
```json
{
  "code": 422,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Email is required"
    }
  ]
}
```

### Conflict
```json
{
  "code": 409,
  "message": "Schedule conflict detected",
  "conflicts": [
    {
      "type": "TEACHER_OVERLAP",
      "message": "Teacher is already assigned at this time"
    }
  ]
}
```

---

## Rate Limiting

API requests are rate limited to:
- 100 requests per minute per user
- 1000 requests per minute per IP

Rate limit headers are included in every response:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642761600
```

---

## Pagination

List endpoints support pagination:

```
GET /admin/centers?limit=10&offset=0
```

Response includes pagination info:
```json
{
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "offset": 0,
    "has_more": true
  }
}
```
