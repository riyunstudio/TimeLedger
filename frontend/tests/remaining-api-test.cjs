const { chromium } = require('playwright');

async function testAllRemainingAPIs() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();

  const results = [];
  const API_BASE = 'http://localhost:8888/api/v1';

  const fetchAPI = async (url, options = {}) => {
    const response = await fetch(url, options);
    return await response.json();
  };

  try {
    // Login first
    console.log('\n=== Login ===');
    const loginResult = await fetchAPI(`${API_BASE}/auth/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'admin@timeledger.com', password: 'admin123' })
    });
    const adminToken = loginResult.datas?.token;
    const centerId = 1;
    console.log('Login:', loginResult.code === 0 ? '✅' : '❌');

    // === Teacher Profile APIs ===
    console.log('\n=== Teacher Profile APIs ===');
    
    console.log('\n--- GET /teacher/me/profile ---');
    const profileResult = await fetchAPI(`${API_BASE}/teacher/me/profile`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Profile:', profileResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Profile', test: 'GET /teacher/me/profile', success: profileResult.code === 0 });

    console.log('\n--- GET /teacher/me/centers ---');
    const teacherCentersResult = await fetchAPI(`${API_BASE}/teacher/me/centers`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Teacher Centers:', teacherCentersResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Profile', test: 'GET /teacher/me/centers', success: teacherCentersResult.code === 0 });

    console.log('\n--- GET /teacher/me/schedule ---');
    const teacherScheduleResult = await fetchAPI(`${API_BASE}/teacher/me/schedule`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Teacher Schedule:', teacherScheduleResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Profile', test: 'GET /teacher/me/schedule', success: teacherScheduleResult.code === 0 });

    // === Teacher Skills APIs ===
    console.log('\n=== Teacher Skills APIs ===');
    
    console.log('\n--- GET /teacher/me/skills ---');
    const skillsResult = await fetchAPI(`${API_BASE}/teacher/me/skills`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Skills:', skillsResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Skills', test: 'GET /teacher/me/skills', success: skillsResult.code === 0 });

    console.log('\n--- POST /teacher/me/skills ---');
    const createSkillResult = await fetchAPI(`${API_BASE}/teacher/me/skills`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '鋼琴', level: 5 })
    });
    console.log('Create Skill:', createSkillResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Skills', test: 'POST /teacher/me/skills', success: createSkillResult.code === 0 });

    // === Teacher Certificates APIs ===
    console.log('\n=== Teacher Certificates APIs ===');
    
    console.log('\n--- GET /teacher/me/certificates ---');
    const certsResult = await fetchAPI(`${API_BASE}/teacher/me/certificates`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Certificates:', certsResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Certificates', test: 'GET /teacher/me/certificates', success: certsResult.code === 0 });

    console.log('\n--- POST /teacher/me/certificates ---');
    const createCertResult = await fetchAPI(`${API_BASE}/teacher/me/certificates`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '音樂教師證書', issuer: '音樂協會', issued_at: '2020-01-01' })
    });
    console.log('Create Certificate:', createCertResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Certificates', test: 'POST /teacher/me/certificates', success: createCertResult.code === 0 });

    // === Teacher Personal Events APIs ===
    console.log('\n=== Teacher Personal Events APIs ===');
    
    console.log('\n--- GET /teacher/me/personal-events ---');
    const personalEventsResult = await fetchAPI(`${API_BASE}/teacher/me/personal-events`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Personal Events:', personalEventsResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Personal Events', test: 'GET /teacher/me/personal-events', success: personalEventsResult.code === 0 });

    console.log('\n--- POST /teacher/me/personal-events ---');
    const createPersonalEventResult = await fetchAPI(`${API_BASE}/teacher/me/personal-events`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: '私人時間', start_at: '2026-02-01 10:00:00', end_at: '2026-02-01 12:00:00', type: 'PERSONAL' })
    });
    console.log('Create Personal Event:', createPersonalEventResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Personal Events', test: 'POST /teacher/me/personal-events', success: createPersonalEventResult.code === 0 });

    // === Teacher Exceptions APIs ===
    console.log('\n=== Teacher Exceptions APIs ===');
    
    console.log('\n--- GET /teacher/exceptions ---');
    const teacherExceptionsResult = await fetchAPI(`${API_BASE}/teacher/exceptions`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Teacher Exceptions:', teacherExceptionsResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Exceptions', test: 'GET /teacher/exceptions', success: teacherExceptionsResult.code === 0 });

    // === Admin Scheduling Validation APIs ===
    console.log('\n=== Admin Scheduling Validation APIs ===');
    
    console.log('\n--- POST /admin/centers/:id/validate ---');
    const validateResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/validate`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ start_date: '2026-02-01', end_date: '2026-02-28' })
    });
    console.log('Validate Schedule:', validateResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Scheduling Validation', test: 'POST /admin/centers/:id/validate', success: validateResult.code === 0 });

    console.log('\n--- POST /admin/centers/:id/scheduling/check-overlap ---');
    const checkOverlapResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/scheduling/check-overlap`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ center_id: centerId, course_id: 1, weekday: 1, start_time: '10:00', end_time: '11:00', exclude_offering_id: 0 })
    });
    console.log('Check Overlap:', checkOverlapResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Scheduling Validation', test: 'POST /admin/centers/:id/scheduling/check-overlap', success: checkOverlapResult.code === 0 });

    // === Admin Schedule Exceptions APIs ===
    console.log('\n=== Admin Schedule Exceptions APIs ===');
    
    console.log('\n--- POST /admin/centers/:id/exceptions ---');
    const createExceptionResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/exceptions`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ original_date: '2026-02-15', type: 'CANCEL', reason: '春節假期' })
    });
    console.log('Create Exception:', createExceptionResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Schedule Exceptions', test: 'POST /admin/centers/:id/exceptions', success: createExceptionResult.code === 0 });

    console.log('\n--- GET /admin/centers/:id/rules/:ruleId/exceptions ---');
    const ruleExceptionsResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/rules/1/exceptions`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Rule Exceptions:', ruleExceptionsResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Schedule Exceptions', test: 'GET /admin/centers/:id/rules/:ruleId/exceptions', success: ruleExceptionsResult.code === 0 });

    // === Smart Matching APIs ===
    console.log('\n=== Smart Matching APIs ===');
    
    console.log('\n--- POST /admin/centers/:id/matching/teachers ---');
    const matchingResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/matching/teachers`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ course_id: 1, weekday: 1, start_time: '10:00', end_time: '11:00' })
    });
    console.log('Find Matching Teachers:', matchingResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Smart Matching', test: 'POST /admin/centers/:id/matching/teachers', success: matchingResult.code === 0 });

    // === Offerings CRUD APIs ===
    console.log('\n=== Offerings CRUD APIs ===');
    
    console.log('\n--- PUT /admin/centers/:id/offerings/:offering_id ---');
    const updateOfferingResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/offerings/1`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '更新後的課程', allow_buffer_override: true })
    });
    console.log('Update Offering:', updateOfferingResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Offerings CRUD', test: 'PUT /admin/centers/:id/offerings/:offering_id', success: updateOfferingResult.code === 0 });

    console.log('\n--- POST /admin/centers/:id/offerings/:offering_id/copy ---');
    const copyOfferingResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/offerings/1/copy`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_weekday: 3 })
    });
    console.log('Copy Offering:', copyOfferingResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Offerings CRUD', test: 'POST /admin/centers/:id/offerings/:offering_id/copy', success: copyOfferingResult.code === 0 });

    // === Timetable Templates APIs ===
    console.log('\n=== Timetable Templates APIs ===');
    
    console.log('\n--- GET /admin/centers/:id/templates ---');
    const templatesResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/templates`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Templates:', templatesResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Timetable Templates', test: 'GET /admin/centers/:id/templates', success: templatesResult.code === 0 });

    console.log('\n--- POST /admin/centers/:id/templates ---');
    const createTemplateResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/templates`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '範本 A', description: '一般課表' })
    });
    console.log('Create Template:', createTemplateResult.code === 0 ? '✅' : '❌');
    const templateId = createTemplateResult.datas?.id;
    results.push({ category: 'Timetable Templates', test: 'POST /admin/centers/:id/templates', success: createTemplateResult.code === 0 });

    if (templateId) {
      console.log('\n--- GET /admin/centers/:id/templates/:templateId ---');
      const getTemplateResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/templates/${templateId}`, {
        headers: { 'Authorization': `Bearer ${adminToken}` }
      });
      console.log('Get Template:', getTemplateResult.code === 0 ? '✅' : '❌');
      results.push({ category: 'Timetable Templates', test: 'GET /admin/centers/:id/templates/:templateId', success: getTemplateResult.code === 0 });
    }

    // === Admin Users CRUD APIs ===
    console.log('\n=== Admin Users CRUD APIs ===');
    
    console.log('\n--- PUT /admin/centers/:id/users/:adminId ---');
    const updateAdminResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/users/16`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '更新的名稱', role: 'ADMIN' })
    });
    console.log('Update Admin User:', updateAdminResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Admin Users CRUD', test: 'PUT /admin/centers/:id/users/:adminId', success: updateAdminResult.code === 0 });

    console.log('\n--- DELETE /admin/centers/:id/users/:adminId ---');
    const deleteAdminResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/users/54`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Delete Admin User:', deleteAdminResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Admin Users CRUD', test: 'DELETE /admin/centers/:id/users/:adminId', success: deleteAdminResult.code === 0 });

    // === Notifications APIs ===
    console.log('\n=== Notifications APIs ===');
    
    console.log('\n--- POST /notifications/:id/read ---');
    const markReadResult = await fetchAPI(`${API_BASE}/notifications/1/read`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Mark As Read:', markReadResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Notifications', test: 'POST /notifications/:id/read', success: markReadResult.code === 0 });

    console.log('\n--- POST /notifications/read-all ---');
    const markAllReadResult = await fetchAPI(`${API_BASE}/notifications/read-all`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Mark All As Read:', markAllReadResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Notifications', test: 'POST /notifications/read-all', success: markAllReadResult.code === 0 });

    // === Export APIs ===
    console.log('\n=== Export APIs ===');
    
    console.log('\n--- POST /admin/export/schedule/csv ---');
    const exportScheduleResult = await fetchAPI(`${API_BASE}/admin/export/schedule/csv`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ center_id: centerId, start_date: '2026-02-01', end_date: '2026-02-28' })
    });
    console.log('Export Schedule CSV:', exportScheduleResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Export', test: 'POST /admin/export/schedule/csv', success: exportScheduleResult.code === 0 });

    console.log('\n--- POST /admin/export/schedule/pdf ---');
    const exportPDFResult = await fetchAPI(`${API_BASE}/admin/export/schedule/pdf`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ center_id: centerId, start_date: '2026-02-01', end_date: '2026-02-28' })
    });
    console.log('Export Schedule PDF:', exportPDFResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Export', test: 'POST /admin/export/schedule/pdf', success: exportPDFResult.code === 0 });

    console.log('\n--- GET /admin/export/centers/:id/exceptions/csv ---');
    const exportExceptionsResult = await fetchAPI(`${API_BASE}/admin/export/centers/${centerId}/exceptions/csv?start_date=2026-01-01&end_date=2026-12-31`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Export Exceptions CSV:', exportExceptionsResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Export', test: 'GET /admin/export/centers/:id/exceptions/csv', success: exportExceptionsResult.code === 0 });

    // === Delete Teacher ===
    console.log('\n=== Delete Teacher ===');
    console.log('\n--- DELETE /teachers/:id ---');
    const deleteTeacherResult = await fetchAPI(`${API_BASE}/teachers/21`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Delete Teacher:', deleteTeacherResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teachers', test: 'DELETE /teachers/:id', success: deleteTeacherResult.code === 0 });

    // === Teacher Session Notes ===
    console.log('\n=== Teacher Session Notes ===');
    console.log('\n--- GET /teacher/sessions/note ---');
    const sessionNoteResult = await fetchAPI(`${API_BASE}/teacher/sessions/note`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Session Note:', sessionNoteResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Session Notes', test: 'GET /teacher/sessions/note', success: sessionNoteResult.code === 0 });

    console.log('\n--- PUT /teacher/sessions/note ---');
    const upsertSessionNoteResult = await fetchAPI(`${API_BASE}/teacher/sessions/note`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ session_id: 1, content: '上課筆記' })
    });
    console.log('Upsert Session Note:', upsertSessionNoteResult.code === 0 ? '✅' : '❌');
    results.push({ category: 'Teacher Session Notes', test: 'PUT /teacher/sessions/note', success: upsertSessionNoteResult.code === 0 });

  } catch (error) {
    console.error('Test Error:', error.message);
    results.push({ category: 'Error', test: error.message, success: false });
  }

  console.log('\n' + '='.repeat(60));
  console.log('All Remaining APIs Test Results Summary');
  console.log('='.repeat(60));
  
  const grouped = results.reduce((acc, r) => {
    if (!acc[r.category]) acc[r.category] = [];
    acc[r.category].push(r);
    return acc;
  }, {});

  for (const [category, tests] of Object.entries(grouped)) {
    console.log(`\n${category}:`);
    tests.forEach(t => {
      console.log(`  ${t.success ? '✅' : '❌'} ${t.test}`);
    });
  }

  const passed = results.filter(r => r.success).length;
  console.log(`\n${'='.repeat(60)}`);
  console.log(`Total: ${passed}/${results.length} tests passed (${Math.round(passed/results.length*100)}%)`);

  await browser.close();
}

testAllRemainingAPIs();
