# TimeLedger User Manual

---

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [Teacher Guide](#teacher-guide)
- [Center Admin Guide](#center-admin-guide)
- [Common Workflows](#common-workflows)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Introduction

TimeLedger is a teacher-centric multi-center scheduling platform designed for music schools and educational institutions in Taiwan. It helps manage:

- **Teacher Schedules**: Weekly class schedules across multiple centers
- **Substitute Management**: Find and book substitute teachers
- **Notifications**: Receive alerts via LINE Notify
- **Profile Management**: Update skills, certificates, and availability
- **Export**: Download schedules and reports

---

## Getting Started

### For Teachers

1. **LINE Login**
   - Open the TimeLedger app or web interface
   - Click "Login with LINE"
   - Authorize LINE access
   - Your profile will be automatically created

2. **Complete Your Profile**
   - Add your skills and proficiency levels
   - Upload teaching certificates
   - Set your availability status
   - Enable/disable hiring status

3. **Join a Center**
   - Accept invitation from a center
   - Or apply to join a center

### For Center Admins

1. **Create Account**
   - Register as an admin user
   - Verify email address

2. **Create Center**
   - Set up your center profile
   - Add rooms and courses
   - Configure schedule templates

3. **Invite Teachers**
   - Send invitations via LINE or email
   - Approve teacher applications

---

## Teacher Guide

### Dashboard

The dashboard shows your weekly schedule across all centers you're associated with.

**Features:**
- View upcoming classes
- Navigate between weeks
- See class details (room, students, notes)

---

### View Schedule

1. Open the app
2. Your schedule for the current week is displayed
3. Use Previous/Next buttons to navigate weeks
4. Tap on a class to view details

**Schedule Legend:**
- 🟢 Green: Confirmed class
- 🟡 Yellow: Pending approval
- 🔴 Red: Cancelled or rescheduled

---

### Accept/Reject Classes

When a center assigns you to a class:

1. You'll receive a LINE notification
2. Open the app
3. Go to the notification
4. Tap "Accept" or "Decline"

---

### Submit Leave Request

To request leave for a scheduled class:

1. Go to your schedule
2. Tap on the class you want to cancel
3. Tap "Request Leave"
4. Select reason and submit
5. Wait for admin approval

**Reasons:**
- Sick leave
- Personal leave
- Emergency
- Other

---

### Add Session Notes

After completing a class:

1. Go to your schedule
2. Tap on the completed class
3. Tap "Add Notes"
4. Enter teaching notes or homework
5. Save

---

### Manage Profile

Update your teaching profile:

1. Go to Profile tab
2. Edit personal information
3. Add/remove skills
4. Upload or remove certificates
5. Update availability status

---

### Enable Hiring Status

To make yourself available for substitute teaching:

1. Go to Profile
2. Toggle "Open to Hiring" to ON
3. Set your contact preferences
4. Save

---

## Center Admin Guide

### Dashboard Overview

The admin dashboard provides:
- Quick stats (classes today, teachers, rooms)
- Upcoming schedule
- Recent notifications
- Action items (pending requests, reviews)

---

### Manage Centers

#### Create a Center

1. Go to Admin > Centers
2. Click "Create Center"
3. Fill in center details:
   - Name
   - Plan level (Starter/Growth/Pro)
   - Contact information
4. Configure settings:
   - Allow public registration
   - Default language
   - Timezone
5. Save

---

### Manage Rooms

#### Add a Room

1. Go to Admin > Resources > Rooms
2. Click "Add Room"
3. Enter room details:
   - Room name
   - Capacity
   - Equipment (piano, whiteboard, etc.)
4. Save

---

### Manage Courses

#### Create a Course

1. Go to Admin > Resources > Courses
2. Click "Add Course"
3. Enter course details:
   - Course name (e.g., Piano Basic)
   - Description
   - Required equipment
   - Buffer times (teacher/room)
4. Save

---

### Create Offerings

An offering is a specific instance of a course.

1. Go to Admin > Resources > Offerings
2. Click "Add Offering"
3. Link to a course
4. Set capacity
5. Assign default room/teacher (optional)
6. Save

---

### Schedule Classes

#### Using Weekly Schedule Grid

1. Go to Admin > Scheduling
2. Select the center
3. View the weekly grid (rooms on X-axis, times on Y-axis)
4. Drag a course to a slot or click empty slot
5. Fill in class details:
   - Course/Offering
   - Teacher
   - Room
   - Time
6. Click "Create"

**Color Coding:**
- Green: No conflicts
- Yellow: Buffer warning
- Red: Conflict detected

---

#### Using Rules

Create recurring weekly schedules:

1. Go to Admin > Scheduling > Rules
2. Click "Add Rule"
3. Configure:
   - Day of week
   - Time range
   - Course/Offering
   - Teacher
   - Room
   - Effective date range
4. Save

The rule will automatically apply to all weeks in the date range.

---

### Handle Exceptions

#### Approve Leave Request

When a teacher submits leave:

1. Go to Admin > Exceptions
2. View pending requests
3. Click on a request
4. Review details
5. Approve or Reject

If approved:
- The class is marked as cancelled
- Or find a substitute (see below)

---

#### Find Substitute

1. Go to the exception detail
2. Click "Find Substitute"
3. System shows available teachers with match scores:
   - Skill match
   - Availability
   - Rating
4. Select a teacher
5. Send invitation

The substitute will receive a LINE notification.

---

#### Reschedule Class

1. Go to the exception
2. Click "Reschedule"
3. Select new date/time
4. Choose available room/teacher
5. Save

---

### Manage Teachers

#### Invite a Teacher

1. Go to Admin > Teachers
2. Click "Invite Teacher"
3. Enter teacher's LINE ID or email
4. Select center role
5. Send invitation

---

#### View Teacher Profile

1. Go to Admin > Teachers
2. Click on a teacher
3. View:
   - Skills and proficiency
   - Certificates
   - Ratings from other centers
   - Schedule conflicts
   - Notes from your center

---

#### Add Center Note

1. Go to teacher profile
2. Click "Add Note"
3. Enter internal note (private to your center)
4. Save

---

### Smart Matching

Use AI-powered matching to find the best substitutes:

1. Go to Admin > Talent Search
2. Set filters:
   - Location (city/district)
   - Required skills
   - Keywords
3. View sorted results by match score
4. Click "Invite" to send offer

---

### Export Reports

#### Export Schedule

1. Go to Admin > Reports
2. Click "Export Schedule"
3. Select date range
4. Choose format (CSV or PDF)
5. Download

---

#### Export Teacher List

1. Go to Admin > Reports
2. Click "Export Teachers"
3. Select center
4. Download CSV

---

#### Export Exceptions

1. Go to Admin > Reports
2. Click "Export Exceptions"
3. Select date range
4. Download CSV

---

## Common Workflows

### Workflow 1: Invite New Teacher

1. Admin creates center invitation
2. Teacher receives LINE notification
3. Teacher accepts invitation
4. Teacher appears in center's teacher list
5. Admin can now schedule classes for this teacher

---

### Workflow 2: Teacher Requests Leave

1. Teacher requests leave in app
2. Admin receives notification
3. Admin reviews request
4. Admin:
   - Approves and finds substitute, OR
   - Rejects with reason
5. Substitute (if found) accepts/declines
6. Final schedule updated

---

### Workflow 3: Substitute Booking

1. Admin needs a substitute for a class
2. Admin uses Smart Matching search
3. Results show available teachers with scores
4. Admin selects best match
5. System sends LINE invitation
6. Substitute accepts
7. Schedule updated

---

### Workflow 4: Weekly Schedule Review

1. Admin views weekly schedule grid
2. Checks for conflicts (red cells)
3. Adjusts classes as needed
4. System validates changes
5. Sends notifications to affected teachers
6. Schedule finalized

---

## Troubleshooting

### Login Issues

**Problem:** Can't login with LINE

**Solutions:**
- Ensure you have LINE app installed
- Check internet connection
- Try logging out and back in
- Contact support if issue persists

---

### Schedule Not Showing

**Problem:** Scheduled classes not appearing

**Solutions:**
- Refresh the app
- Check correct center is selected
- Verify date range is correct
- Contact admin if classes are missing

---

### Notifications Not Received

**Problem:** Not receiving LINE notifications

**Solutions:**
- Check LINE Notify token is set
- Ensure LINE Notify is linked
- Check notification preferences
- Test with "Send Test Notification"

---

### Can't Accept Class

**Problem:** Unable to accept assigned class

**Solutions:**
- Check if you have a schedule conflict
- Verify class hasn't been cancelled
- Contact admin

---

### Substitute Not Found

**Problem:** No substitutes available for a class

**Solutions:**
- Expand search criteria
- Increase match score threshold
- Consider teachers from nearby districts
- Contact teachers directly

---

## FAQ

### Q: Can I work at multiple centers?

**A:** Yes! Teachers can be members of multiple centers and see all their schedules in one unified view.

---

### Q: How do I become a substitute?

**A:** Enable "Open to Hiring" in your profile. Centers using Smart Matching will be able to find and invite you.

---

### Q: Can I decline a class assignment?

**A:** Yes, but please do so promptly so the center can find another teacher.

---

### Q: What happens if I'm sick and can't teach?

**A:** Submit a leave request in the app. The admin will find a substitute or reschedule the class.

---

### Q: How do I add my teaching certificates?

**A:** Go to your profile and click "Upload Certificate". Supported formats: PDF, JPG, PNG.

---

### Q: Can parents see my contact info?

**A:** Only if you enable it in your profile settings. You control your privacy.

---

### Q: How are substitute teachers matched?

**A:** Our Smart Matching algorithm considers:
- Required skills and proficiency
- Availability at the time
- Location/district
- Ratings from other centers
- Past performance

---

### Q: What's the difference between "Starter", "Growth", and "Pro" plans?

**A:**
- **Starter**: Up to 10 teachers, basic scheduling
- **Growth**: Up to 50 teachers, smart matching
- **Pro**: Unlimited teachers, priority support, custom features

---

### Q: How do I export my schedule?

**A:** Go to Reports and click "Export Schedule". Choose CSV or PDF format.

---

### Q: Can I use TimeLedger on desktop?

**A:** Yes, TimeLedger has a web interface accessible from any browser.

---

### Q: How do I change my center name or settings?

**A:** Center admins can edit center details from Admin > Centers.

---

### Q: What if I forget to submit a leave request?

**A:** Contact the admin immediately. They can create an exception on your behalf.

---

### Q: Is there a limit to how many certificates I can upload?

**A:** No, upload as many teaching certificates as you have.

---

### Q: How do I delete my account?

**A:** Contact support to request account deletion. All your data will be permanently removed.

---

## Support

### Contact Us

- **Email**: support@timeledger.com.tw
- **LINE Official Account**: @timeledger
- **Phone**: +886-2-1234-5678

### Business Hours

- Monday - Friday: 9:00 - 18:00 (GMT+8)
- Saturday: 10:00 - 14:00 (GMT+8)
- Sunday: Closed

### Emergency Support

For urgent issues outside business hours, use LINE Official Account.

---

## Version History

| Version | Date | Changes |
|:---|:---|:---|
| 1.0.0 | 2026-01-21 | Initial release |
