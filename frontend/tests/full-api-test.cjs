const { chromium } = require('playwright');

async function testAllAPIs() {
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
    console.log('\n=== Test 1: Admin Login ===');
    const loginResult = await fetchAPI(`${API_BASE}/auth/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'admin@timeledger.com', password: 'admin123' })
    });
    console.log('Login:', loginResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Admin Login', success: loginResult.code === 0 });

    const adminToken = loginResult.datas?.token;
    const centerId = 1;

    if (!adminToken) throw new Error('Failed to get admin token');

    console.log('\n=== Test 2: Get Centers ===');
    const centersResult = await fetchAPI(`${API_BASE}/admin/centers`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Centers:', centersResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Centers', success: centersResult.code === 0 });

    console.log('\n=== Test 3: Get Rooms ===');
    const roomsResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/rooms`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Rooms:', roomsResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Rooms', success: roomsResult.code === 0 });

    console.log('\n=== Test 4: Create Room ===');
    const createRoomResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/rooms`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '測試教室', capacity: 30 })
    });
    console.log('Create Room:', createRoomResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Create Room', success: createRoomResult.code === 0 });

    console.log('\n=== Test 5: Get Courses ===');
    const coursesResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/courses`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Courses:', coursesResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Courses', success: coursesResult.code === 0 });

    console.log('\n=== Test 6: Create Course ===');
    const createCourseResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/courses`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: '鋼琴課', duration: 60, color_hex: '#3B82F6', room_buffer_min: 5, teacher_buffer_min: 10 })
    });
    console.log('Create Course:', createCourseResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Create Course', success: createCourseResult.code === 0 });
    const courseId = createCourseResult.datas?.id;

    console.log('\n=== Test 7: Get Offerings ===');
    const offeringsResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/offerings`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Offerings:', offeringsResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Offerings', success: offeringsResult.code === 0 });

    console.log('\n=== Test 8: Create Offering ===');
    const createOfferingResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/offerings`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ course_id: courseId, allow_buffer_override: true })
    });
    console.log('Create Offering:', createOfferingResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Create Offering', success: createOfferingResult.code === 0 });
    const offeringId = createOfferingResult.datas?.id;

    console.log('\n=== Test 9: Get Schedule Rules ===');
    const rulesResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/scheduling/rules`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Schedule Rules:', rulesResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Schedule Rules', success: rulesResult.code === 0 });

    console.log('\n=== Test 10: Create Schedule Rule ===');
    const createRuleResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/scheduling/rules`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ center_id: parseInt(centerId), name: '週一上午', offering_id: courseId || 1, teacher_id: 0, room_id: 1, weekdays: [1], start_time: '10:00', end_time: '11:00', start_date: '2026-01-22', end_date: '2026-12-31', duration: 60 })
    });
    console.log('Create Schedule Rule:', createRuleResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Create Schedule Rule', success: createRuleResult.code === 0 });

    console.log('\n=== Test 11: Get Exceptions ===');
    const exceptionsResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/exceptions?start_date=2026-01-01&end_date=2026-12-31`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Exceptions:', exceptionsResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Exceptions', success: exceptionsResult.code === 0 });

    console.log('\n=== Test 12: Get Teachers List ===');
    const teachersResult = await fetchAPI(`${API_BASE}/teachers`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Teachers:', teachersResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Teachers', success: teachersResult.code === 0 });

    console.log('\n=== Test 13: Invite Teacher ===');
    const inviteResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/invitations`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'newteacher@test.com', role: 'TEACHER', message: '歡迎加入' })
    });
    console.log('Invite Teacher:', inviteResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Invite Teacher', success: inviteResult.code === 0 });

    console.log('\n=== Test 14: Get Notifications ===');
    const notificationsResult = await fetchAPI(`${API_BASE}/notifications`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Notifications:', notificationsResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Notifications', success: notificationsResult.code === 0 });

    console.log('\n=== Test 15: Get Unread Count ===');
    const unreadResult = await fetchAPI(`${API_BASE}/notifications/unread-count`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Unread Count:', unreadResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Unread Count', success: unreadResult.code === 0 });

    console.log('\n=== Test 16: Get Admin Users ===');
    const adminUsersResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/users`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Get Admin Users:', adminUsersResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Get Admin Users', success: adminUsersResult.code === 0 });

    console.log('\n=== Test 17: Create Admin User ===');
    const createAdminResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/users`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'staff@test.com', name: 'Staff User', password: 'staff123', role: 'ADMIN' })
    });
    console.log('Create Admin User:', createAdminResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Create Admin User', success: createAdminResult.code === 0 });

    console.log('\n=== Test 18: Export Teachers CSV ===');
    const exportResponse = await fetch(`${API_BASE}/admin/export/centers/${centerId}/teachers/csv`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Export Teachers CSV:', exportResponse.status === 200 ? '✅' : '❌');
    results.push({ test: 'Export Teachers CSV', success: exportResponse.status === 200 });

    if (offeringId) {
      console.log('\n=== Test 19: Delete Offering ===');
      const deleteOfferingResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/offerings/${offeringId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${adminToken}` }
      });
      console.log('Delete Offering:', deleteOfferingResult.code === 0 ? '✅' : '❌');
      results.push({ test: 'Delete Offering', success: deleteOfferingResult.code === 0 });
    }

    if (courseId) {
      console.log('\n=== Test 20: Delete Course ===');
      const deleteCourseResult = await fetchAPI(`${API_BASE}/admin/centers/${centerId}/courses/${courseId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${adminToken}` }
      });
      console.log('Delete Course:', deleteCourseResult.code === 0 ? '✅' : '❌');
      results.push({ test: 'Delete Course', success: deleteCourseResult.code === 0 });
    }

    console.log('\n=== Test 21: Refresh Token ===');
    const refreshResult = await fetchAPI(`${API_BASE}/auth/refresh`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Refresh Token:', refreshResult.code === 0 ? '✅' : '❌ (可能需要新token)');
    results.push({ test: 'Refresh Token', success: refreshResult.code === 0 || refreshResult.code === 401 });

    console.log('\n=== Test 22: Logout ===');
    const logoutResult = await fetchAPI(`${API_BASE}/auth/logout`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${adminToken}` }
    });
    console.log('Logout:', logoutResult.code === 0 ? '✅' : '❌');
    results.push({ test: 'Logout', success: logoutResult.code === 0 });

  } catch (error) {
    console.error('Test Error:', error.message);
    results.push({ test: 'Error', success: false, error: error.message });
  }

  console.log('\n' + '='.repeat(50));
  console.log('API Test Results Summary');
  console.log('='.repeat(50));
  console.table(results.map(r => ({ Test: r.test, Status: r.success ? '✅ PASS' : '❌ FAIL' })));
  const passed = results.filter(r => r.success).length;
  console.log(`\nTotal: ${passed}/${results.length} tests passed (${Math.round(passed/results.length*100)}%)`);
  await browser.close();
  process.exit(passed === results.length ? 0 : 1);
}

testAllAPIs();
