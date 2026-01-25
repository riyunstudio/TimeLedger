import { describe, it, expect } from 'vitest'
import type { ScheduleException } from '~/types'

describe('ScheduleException Type', () => {
  it('should accept CANCEL type', () => {
    const exception: ScheduleException = {
      id: 1,
      center_id: 1,
      rule_id: 1,
      teacher_id: 1,
      original_date: '2026-01-25',
      type: 'CANCEL',
      status: 'PENDING',
      reason: 'test',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    expect(exception.type).toBe('CANCEL')
  })

  it('should accept RESCHEDULE type', () => {
    const exception: ScheduleException = {
      id: 1,
      center_id: 1,
      rule_id: 1,
      teacher_id: 1,
      original_date: '2026-01-25',
      type: 'RESCHEDULE',
      status: 'PENDING',
      new_start_at: '2026-01-26T10:00:00Z',
      new_end_at: '2026-01-26T11:00:00Z',
      reason: 'test',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    expect(exception.type).toBe('RESCHEDULE')
  })

  it('should accept REPLACE_TEACHER type', () => {
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
      reason: 'test',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    expect(exception.type).toBe('REPLACE_TEACHER')
    expect(exception.new_teacher_id).toBe(2)
    expect(exception.new_teacher_name).toBe('代課老師')
  })

  it('should allow all status values', () => {
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
        reason: 'test',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      expect(exception.status).toBe(status)
    })
  })

  it('should have optional fields', () => {
    const exception: ScheduleException = {
      id: 1,
      center_id: 1,
      rule_id: 1,
      teacher_id: 1,
      original_date: '2026-01-25',
      type: 'CANCEL',
      status: 'PENDING',
      reason: 'test',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    expect(exception.new_start_at).toBeUndefined()
    expect(exception.new_end_at).toBeUndefined()
    expect(exception.new_teacher_id).toBeUndefined()
    expect(exception.new_teacher_name).toBeUndefined()
    expect(exception.new_room_id).toBeUndefined()
  })
})
