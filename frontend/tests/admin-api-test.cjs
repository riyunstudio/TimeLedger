const API_BASE = 'http://localhost:8888/api/v1'
const ADMIN_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3R5cGUiOiJBRE1JTiIsInVzZXJfaWQiOjEsImNlbnRlcl9pZCI6NCwiZXhwIjoxNzY5MTQ1OTYyfQ==.Vt9eVCd-qQZ6nKeSIqS68_8ztIwq_e4GhKU6cnqKoc0='

const headers = {
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${ADMIN_TOKEN}`
}

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))

const tests = {
  passed: 0,
  failed: 0,
  results: []
}

async function test(name, fn) {
  process.stdout.write(`Testing: ${name}... `)
  try {
    await fn()
    console.log('✅ PASSED')
    tests.passed++
    tests.results.push({ name, status: 'PASSED' })
  } catch (error) {
    console.log(`❌ FAILED: ${error.message}`)
    tests.failed++
    tests.results.push({ name, status: 'FAILED', error: error.message })
  }
}

async function get(endpoint) {
  const res = await fetch(`${API_BASE}${endpoint}`, { headers })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`)
  return res.json()
}

async function post(endpoint, body) {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'POST',
    headers,
    body: JSON.stringify(body)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`)
  return res.json()
}

async function postCSV(endpoint, body) {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'POST',
    headers,
    body: JSON.stringify(body)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`)
  const contentType = res.headers.get('content-type')
  if (!contentType || !contentType.includes('text/csv')) {
    throw new Error(`Expected CSV content-type, got ${contentType}`)
  }
  return res.text()
}

async function put(endpoint, body) {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'PUT',
    headers,
    body: JSON.stringify(body)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`)
  return res.json()
}

async function deleteReq(endpoint) {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'DELETE',
    headers
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`)
  return res.json()
}

async function runTests() {
  console.log('='.repeat(60))
  console.log('ADMIN API COMPREHENSIVE TEST')
  console.log('='.repeat(60))
  console.log()

  const CENTER_ID = 4  // Admin has access to center 4
  const TEACHER_ID = 1
  const ROOM_ID = 1
  const COURSE_ID = 1
  const OFFERING_ID = 1

  // ============================================
  // SCHEDULING VALIDATION APIs
  // ============================================
  console.log('\n--- Scheduling Validation APIs ---\n')

  await test('POST /admin/centers/:id/validate (Validate Schedule)', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/validate`, {
      center_id: CENTER_ID,
      room_id: ROOM_ID,
      course_id: COURSE_ID,
      start_time: '2026-02-15T10:00:00Z',
      end_time: '2026-02-15T11:00:00Z',
      allow_buffer_override: false
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  await test('POST /admin/centers/:id/scheduling/check-overlap', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/scheduling/check-overlap`, {
      center_id: CENTER_ID,
      room_id: ROOM_ID,
      start_time: '2026-02-15T10:00:00Z',
      end_time: '2026-02-15T11:00:00Z'
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  await test('POST /admin/centers/:id/scheduling/check-teacher-buffer', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/scheduling/check-teacher-buffer`, {
      center_id: CENTER_ID,
      teacher_id: TEACHER_ID,
      room_id: ROOM_ID,
      prev_end_time: '2026-02-15T09:00:00Z',
      next_start_time: '2026-02-15T10:30:00Z',
      course_id: COURSE_ID
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  await test('POST /admin/centers/:id/scheduling/check-room-buffer', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/scheduling/check-room-buffer`, {
      center_id: CENTER_ID,
      teacher_id: TEACHER_ID,
      room_id: ROOM_ID,
      prev_end_time: '2026-02-15T09:00:00Z',
      next_start_time: '2026-02-15T10:30:00Z',
      course_id: COURSE_ID
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  // ============================================
  // SCHEDULING RULES APIs
  // ============================================
  console.log('\n--- Scheduling Rules APIs ---\n')

  await test('GET /admin/centers/:id/scheduling/rules', async () => {
    const result = await get(`/admin/centers/${CENTER_ID}/scheduling/rules`)
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    if (!Array.isArray(result.datas)) throw new Error('Expected datas to be array')
  })

  await test('POST /admin/centers/:id/scheduling/rules', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/scheduling/rules`, {
      center_id: CENTER_ID,
      name: 'Test Schedule Rule',
      offering_id: OFFERING_ID,
      teacher_id: TEACHER_ID,
      room_id: ROOM_ID,
      start_time: '10:00',
      end_time: '11:00',
      duration: 60,
      weekdays: [1, 3, 5],
      start_date: '2026-01-01'
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  // ============================================
  // SCHEDULE EXCEPTIONS APIs
  // ============================================
  console.log('\n--- Schedule Exceptions APIs ---\n')

  await test('GET /admin/centers/:id/exceptions', async () => {
    const result = await get(`/admin/centers/${CENTER_ID}/exceptions?start_date=2026-01-01&end_date=2026-12-31`)
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    if (!Array.isArray(result.datas)) throw new Error('Expected datas to be array')
  })

  // ============================================
  // SMART MATCHING APIs
  // ============================================
  console.log('\n--- Smart Matching APIs ---\n')

  await test('POST /admin/centers/:id/matching/teachers', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/matching/teachers`, {
      center_id: CENTER_ID,
      room_id: ROOM_ID,
      start_time: '2026-02-15T10:00:00Z',
      end_time: '2026-02-15T11:00:00Z',
      required_skills: ['鋼琴']
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  await test('GET /admin/centers/:id/matching/teachers/search', async () => {
    const result = await get(`/admin/centers/${CENTER_ID}/matching/teachers/search?city=台北市&skills=鋼琴`)
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  // ============================================
  // TIMETABLE TEMPLATE APIs
  // ============================================
  console.log('\n--- Timetable Template APIs ---\n')

  let templateId = 0

  await test('GET /admin/centers/:id/templates', async () => {
    const result = await get(`/admin/centers/${CENTER_ID}/templates`)
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    if (!Array.isArray(result.datas)) throw new Error('Expected datas to be array')
  })

  await test('POST /admin/centers/:id/templates', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/templates`, {
      name: 'Test Template',
      row_type: 'ROOM'
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    templateId = result.datas?.id || 0
  })

  if (templateId > 0) {
    await test('PUT /admin/centers/:id/templates/:templateId', async () => {
      const result = await put(`/admin/centers/${CENTER_ID}/templates/${templateId}`, {
        name: 'Updated Template Name'
      })
      if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    })

    await test('GET /admin/centers/:id/templates/:templateId/cells', async () => {
      const result = await get(`/admin/centers/${CENTER_ID}/templates/${templateId}/cells`)
      if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    })

    await test('POST /admin/centers/:id/templates/:templateId/cells', async () => {
      const result = await post(`/admin/centers/${CENTER_ID}/templates/${templateId}/cells`, [
        {
          row_no: 1,
          col_no: 1,
          start_time: '09:00',
          end_time: '10:00',
          room_id: ROOM_ID
        }
      ])
      if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    })

    await test('DELETE /admin/centers/:id/templates/:templateId', async () => {
      const result = await deleteReq(`/admin/centers/${CENTER_ID}/templates/${templateId}`)
      if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    })
  }

  // ============================================
  // ADMIN USERS APIs
  // ============================================
  console.log('\n--- Admin Users APIs ---\n')

  await test('GET /admin/centers/:id/users', async () => {
    const result = await get(`/admin/centers/${CENTER_ID}/users`)
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
    if (!Array.isArray(result.datas)) throw new Error('Expected datas to be array')
  })

  // ============================================
  // SCHEDULE EXPANSION API
  // ============================================
  console.log('\n--- Schedule Expansion APIs ---\n')

  await test('POST /admin/centers/:id/expand', async () => {
    const result = await post(`/admin/centers/${CENTER_ID}/expand`, {
      rule_ids: [],
      start_date: '2026-02-01T00:00:00Z',
      end_date: '2026-02-28T00:00:00Z'
    })
    if (result.code !== 0) throw new Error(`Expected code 0, got ${result.code}`)
  })

  // ============================================
  // EXPORT APIs
  // ============================================
  console.log('\n--- Export APIs ---\n')

  await test('POST /admin/export/schedule/csv', async () => {
    const csvData = await postCSV(`/admin/export/schedule/csv`, {
      center_id: CENTER_ID,
      start_date: '2026-02-01T00:00:00Z',
      end_date: '2026-02-28T00:00:00Z'
    })
    if (!csvData || csvData.length < 10) throw new Error('Empty or too short CSV response')
    console.log('CSV export received, length:', csvData.length)
  })

  // ============================================
  // SUMMARY
  // ============================================
  console.log('\n' + '='.repeat(60))
  console.log('TEST SUMMARY')
  console.log('='.repeat(60))
  console.log(`Total: ${tests.passed + tests.failed}`)
  console.log(`✅ Passed: ${tests.passed}`)
  console.log(`❌ Failed: ${tests.failed}`)
  console.log('='.repeat(60))

  if (tests.failed > 0) {
    console.log('\nFailed tests:')
    tests.results.filter(r => r.status === 'FAILED').forEach(r => {
      console.log(`  - ${r.name}: ${r.error}`)
    })
  }

  console.log('\n')
  process.exit(tests.failed > 0 ? 1 : 0)
}

runTests().catch(console.error)
