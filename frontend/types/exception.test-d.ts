/**
 * Type tests for ScheduleException
 * Run with: npx tsc --noEmit
 */

// This file tests TypeScript types for the exception feature
// It will cause compile errors if types are incorrect

import type { ScheduleException } from './index'

// Test that REPLACE_TEACHER type is valid
function testReplaceTeacherTypes() {
  const exception: ScheduleException = {
    id: 1,
    center_id: 1,
    rule_id: 1,
    teacher_id: 1,
    original_date: '2026-01-25',
    type: 'REPLACE_TEACHER',
    status: 'PENDING',
    new_teacher_id: 2,
    new_teacher_name: '代課老師',
    reason: '原因',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }

  // Verify all fields are accessible
  const _type: string = exception.type
  const _teacherId: number | undefined = exception.new_teacher_id
  const _teacherName: string | undefined = exception.new_teacher_name
}

// Test that all exception types are valid
function testAllExceptionTypes() {
  const types: ScheduleException['type'][] = ['CANCEL', 'RESCHEDULE', 'REPLACE_TEACHER']

  types.forEach(type => {
    const exception: ScheduleException = {
      id: 1,
      center_id: 1,
      rule_id: 1,
      teacher_id: 1,
      original_date: '2026-01-25',
      type: type,
      status: 'PENDING',
      reason: '原因',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
    const _type: string = exception.type
  })
}

// Test that all status values are valid
function testAllStatusValues() {
  const statuses: ScheduleException['status'][] = ['PENDING', 'APPROVED', 'REJECTED', 'REVOKED']

  statuses.forEach(status => {
    const exception: ScheduleException = {
      id: 1,
      center_id: 1,
      rule_id: 1,
      teacher_id: 1,
      original_date: '2026-01-25',
      type: 'CANCEL',
      status: status,
      reason: '原因',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
    const _status: string = exception.status
  })
}
