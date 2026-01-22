/**
 * Teacher APIs Comprehensive Test Suite
 * 
 * This script tests all Teacher-related APIs using the mock JWT token.
 * 
 * Usage: node tests/teacher-api-test.cjs
 */

const { chromium } = require('playwright');

const API_BASE = 'http://localhost:8888/api/v1';
const TEACHER_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3R5cGUiOiJURUFDSEVSIiwidXNlcl9pZCI6MSwiY2VudGVyX2lkIjoxLCJsaW5lX3VzZXJfaWQiOiJMSU5FX1VTRVJfMDAxIiwiZXhwIjoxNzY5MTQ1MDgyfQ==.saYRMrPMyxtccJhIbbXyagW-lBV3RcIVllTkuTWQOtM=';

async function testTeacherAPIs() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();

  const results = [];
  const fetchAPI = async (url, options = {}) => {
    const response = await fetch(url, options);
    return await response.json();
  };

  try {
    console.log('\n========================================');
    console.log('Teacher APIs Comprehensive Test Suite');
    console.log('========================================');
    console.log('\nTeacher ID: 1 (老師1)');
    console.log('Token: ' + TEACHER_TOKEN.substring(0, 50) + '...\n');

    // === Teacher Profile APIs ===
    console.log('=== 1. Teacher Profile APIs ===\n');

    console.log('--- GET /teacher/me/profile ---');
    const profileResult = await fetchAPI(`${API_BASE}/teacher/me/profile`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', profileResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    console.log('  User:', profileResult.datas?.name || 'N/A');
    results.push({ category: 'Teacher Profile', test: 'GET /teacher/me/profile', success: profileResult.code === 0 });

    console.log('\n--- PUT /teacher/me/profile ---');
    const updateProfileResult = await fetchAPI(`${API_BASE}/teacher/me/profile`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '老師一', bio: '更新後的自我介紹' })
    });
    console.log('Response:', updateProfileResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Profile', test: 'PUT /teacher/me/profile', success: updateProfileResult.code === 0 });

    console.log('\n--- GET /teacher/me/centers ---');
    const centersResult = await fetchAPI(`${API_BASE}/teacher/me/centers`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', centersResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    const centerCount = centersResult.datas?.length || 0;
    console.log(`  Centers: ${centerCount}`);
    results.push({ category: 'Teacher Profile', test: 'GET /teacher/me/centers', success: centersResult.code === 0 });

    // === Teacher Schedule APIs ===
    console.log('\n========================================');
    console.log('=== 2. Teacher Schedule APIs ===\n');

    console.log('--- GET /teacher/me/schedule ---');
    const scheduleResult = await fetchAPI(`${API_BASE}/teacher/me/schedule?from=2026-01-01&to=2026-12-31`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', scheduleResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Schedule', test: 'GET /teacher/me/schedule', success: scheduleResult.code === 0 });

    // === Teacher Exceptions APIs ===
    console.log('\n========================================');
    console.log('=== 3. Teacher Exceptions APIs ===\n');

    console.log('--- GET /teacher/exceptions ---');
    const exceptionsResult = await fetchAPI(`${API_BASE}/teacher/exceptions`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', exceptionsResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Exceptions', test: 'GET /teacher/exceptions', success: exceptionsResult.code === 0 });

    console.log('\n--- POST /teacher/exceptions ---');
    const createExceptionResult = await fetchAPI(`${API_BASE}/teacher/exceptions`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({
        center_id: 1,
        rule_id: 1,
        original_date: '2026-02-15T00:00:00Z',
        type: 'CANCEL',
        reason: '身體不適需要請假'
      })
    });
    console.log('Response:', createExceptionResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    const exceptionId = createExceptionResult.datas?.id;
    results.push({ category: 'Teacher Exceptions', test: 'POST /teacher/exceptions', success: createExceptionResult.code === 0 });

    if (exceptionId) {
      console.log(`\n--- POST /teacher/exceptions/${exceptionId}/revoke ---`);
      const revokeResult = await fetchAPI(`${API_BASE}/teacher/exceptions/${exceptionId}/revoke`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ reason: '已康復' })
      });
      console.log('Response:', revokeResult.code === 0 ? '✅ PASS' : '❌ FAIL');
      results.push({ category: 'Teacher Exceptions', test: `POST /teacher/exceptions/${exceptionId}/revoke`, success: revokeResult.code === 0 });
    }

    // === Teacher Session Notes APIs ===
    console.log('\n========================================');
    console.log('=== 4. Teacher Session Notes APIs ===\n');

    console.log('--- GET /teacher/sessions/note ---');
    const sessionNoteResult = await fetchAPI(`${API_BASE}/teacher/sessions/note?rule_id=1&session_date=2026-02-01`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', sessionNoteResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Session Notes', test: 'GET /teacher/sessions/note', success: sessionNoteResult.code === 0 });

    console.log('\n--- PUT /teacher/sessions/note ---');
    const upsertNoteResult = await fetchAPI(`${API_BASE}/teacher/sessions/note`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ rule_id: 1, session_date: '2026-02-01', content: '今天的課程進度順利，學生表現優異。' })
    });
    console.log('Response:', upsertNoteResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Session Notes', test: 'PUT /teacher/sessions/note', success: upsertNoteResult.code === 0 });

    // === Teacher Skills APIs ===
    console.log('\n========================================');
    console.log('=== 5. Teacher Skills APIs ===\n');

    console.log('--- GET /teacher/me/skills ---');
    const skillsResult = await fetchAPI(`${API_BASE}/teacher/me/skills`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', skillsResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Skills', test: 'GET /teacher/me/skills', success: skillsResult.code === 0 });

    console.log('\n--- POST /teacher/me/skills ---');
    const createSkillResult = await fetchAPI(`${API_BASE}/teacher/me/skills`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ category: '樂器', skill_name: '鋼琴', level: 'ADVANCED' })
    });
    console.log('Response:', createSkillResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    const skillId = createSkillResult.datas?.id;
    results.push({ category: 'Teacher Skills', test: 'POST /teacher/me/skills', success: createSkillResult.code === 0 });

    if (skillId) {
      console.log(`\n--- DELETE /teacher/me/skills/${skillId} ---`);
      const deleteSkillResult = await fetchAPI(`${API_BASE}/teacher/me/skills/${skillId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
      });
      console.log('Response:', deleteSkillResult.code === 0 ? '✅ PASS' : '❌ FAIL');
      results.push({ category: 'Teacher Skills', test: `DELETE /teacher/me/skills/${skillId}`, success: deleteSkillResult.code === 0 });
    }

    // === Teacher Certificates APIs ===
    console.log('\n========================================');
    console.log('=== 6. Teacher Certificates APIs ===\n');

    console.log('--- GET /teacher/me/certificates ---');
    const certsResult = await fetchAPI(`${API_BASE}/teacher/me/certificates`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', certsResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Certificates', test: 'GET /teacher/me/certificates', success: certsResult.code === 0 });

    console.log('\n--- POST /teacher/me/certificates ---');
    const createCertResult = await fetchAPI(`${API_BASE}/teacher/me/certificates`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: '音樂教師資格證書',
        file_url: 'https://example.com/cert.pdf',
        issued_at: '2020-05-15T00:00:00Z'
      })
    });
    console.log('Response:', createCertResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    const certId = createCertResult.datas?.id;
    results.push({ category: 'Teacher Certificates', test: 'POST /teacher/me/certificates', success: createCertResult.code === 0 });

    if (certId) {
      console.log(`\n--- DELETE /teacher/me/certificates/${certId} ---`);
      const deleteCertResult = await fetchAPI(`${API_BASE}/teacher/me/certificates/${certId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
      });
      console.log('Response:', deleteCertResult.code === 0 ? '✅ PASS' : '❌ FAIL');
      results.push({ category: 'Teacher Certificates', test: `DELETE /teacher/me/certificates/${certId}`, success: deleteCertResult.code === 0 });
    }

    // === Teacher Personal Events APIs ===
    console.log('\n========================================');
    console.log('=== 7. Teacher Personal Events APIs ===\n');

    console.log('--- GET /teacher/me/personal-events ---');
    const personalEventsResult = await fetchAPI(`${API_BASE}/teacher/me/personal-events`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', personalEventsResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Teacher Personal Events', test: 'GET /teacher/me/personal-events', success: personalEventsResult.code === 0 });

    console.log('\n--- POST /teacher/me/personal-events ---');
    const createPersonalEventResult = await fetchAPI(`${API_BASE}/teacher/me/personal-events`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: '私人預約',
        start_at: '2026-02-20T14:00:00Z',
        end_at: '2026-02-20T16:00:00Z',
        is_all_day: false
      })
    });
    console.log('Response:', createPersonalEventResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    const eventId = createPersonalEventResult.datas?.id;
    results.push({ category: 'Teacher Personal Events', test: 'POST /teacher/me/personal-events', success: createPersonalEventResult.code === 0 });

    if (eventId) {
      console.log(`\n--- PATCH /teacher/me/personal-events/${eventId} ---`);
      const updateEventResult = await fetchAPI(`${API_BASE}/teacher/me/personal-events/${eventId}`, {
        method: 'PATCH',
        headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: '更新：私人預約', update_mode: 'SINGLE' })
      });
      console.log('Response:', updateEventResult.code === 0 ? '✅ PASS' : '❌ FAIL');
      results.push({ category: 'Teacher Personal Events', test: `PATCH /teacher/me/personal-events/${eventId}`, success: updateEventResult.code === 0 });

      console.log(`\n--- DELETE /teacher/me/personal-events/${eventId} ---`);
      const deleteEventResult = await fetchAPI(`${API_BASE}/teacher/me/personal-events/${eventId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
      });
      console.log('Response:', deleteEventResult.code === 0 ? '✅ PASS' : '❌ FAIL');
      results.push({ category: 'Teacher Personal Events', test: `DELETE /teacher/me/personal-events/${eventId}`, success: deleteEventResult.code === 0 });
    }

    // === Notifications APIs ===
    console.log('\n========================================');
    console.log('=== 8. Notifications APIs ===\n');

    console.log('--- GET /notifications ---');
    const notificationsResult = await fetchAPI(`${API_BASE}/notifications`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', notificationsResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Notifications', test: 'GET /notifications', success: notificationsResult.code === 0 });

    console.log('\n--- GET /notifications/unread-count ---');
    const unreadResult = await fetchAPI(`${API_BASE}/notifications/unread-count`, {
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', unreadResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Notifications', test: 'GET /notifications/unread-count', success: unreadResult.code === 0 });

    // === Logout ===
    console.log('\n========================================');
    console.log('=== 9. Logout ===\n');

    console.log('--- POST /auth/logout ---');
    const logoutResult = await fetchAPI(`${API_BASE}/auth/logout`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${TEACHER_TOKEN}` }
    });
    console.log('Response:', logoutResult.code === 0 ? '✅ PASS' : '❌ FAIL');
    results.push({ category: 'Auth', test: 'POST /auth/logout', success: logoutResult.code === 0 });

  } catch (error) {
    console.error('\n❌ Test Error:', error.message);
    results.push({ category: 'Error', test: error.message, success: false });
  }

  // Summary
  console.log('\n========================================');
  console.log('Teacher APIs Test Results Summary');
  console.log('========================================');

  const grouped = results.reduce((acc, r) => {
    if (!acc[r.category]) acc[r.category] = [];
    acc[r.category].push(r);
    return acc;
  }, {});

  for (const [category, tests] of Object.entries(grouped)) {
    const passed = tests.filter(t => t.success).length;
    console.log(`\n${category}: ${passed}/${tests.length} passed`);
    tests.forEach(t => {
      console.log(`  ${t.success ? '✅' : '❌'} ${t.test}`);
    });
  }

  const totalPassed = results.filter(r => r.success).length;
  console.log(`\n${'='.repeat(50)}`);
  console.log(`Total: ${totalPassed}/${results.length} tests passed (${Math.round(totalPassed/results.length*100)}%)`);

  await browser.close();
  process.exit(totalPassed === results.length ? 0 : 1);
}

testTeacherAPIs();
