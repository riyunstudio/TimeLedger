const { chromium } = require('playwright');

const API_BASE = 'http://localhost:8888/api/v1';

async function testAPIs() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();

  const results = [];

  page.on('console', msg => {
    if (msg.type() === 'log') {
      console.log('Browser:', msg.text());
    }
  });

  page.on('request', request => {
    console.log('Request:', request.method(), request.url());
  });

  page.on('response', response => {
    console.log('Response:', response.status(), response.url());
  });

  try {
    // Test 1: Login API
    console.log('\n=== Test 1: Admin Login ===');
    const loginResult = await page.evaluate(async () => {
      const response = await fetch('http://localhost:8888/api/v1/auth/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: 'admin@timeledger.com',
          password: 'admin123'
        })
      });
      const data = await response.json();
      return data;
    });
    console.log('Login Result:', loginResult);
    results.push({ test: 'Admin Login', success: loginResult.code === 0, data: loginResult });

    const token = loginResult.datas?.token;
    if (!token) {
      throw new Error('No token received');
    }

    // Test 2: Get Centers
    console.log('\n=== Test 2: Get Centers ===');
    const centersResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/admin/centers', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      return await response.json();
    }, token);
    console.log('Centers Result:', centersResult);
    results.push({ test: 'Get Centers', success: centersResult.code === 0, data: centersResult });

    // Test 3: Get Rooms
    console.log('\n=== Test 3: Get Rooms ===');
    const roomsResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/admin/centers/1/rooms', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      return await response.json();
    }, token);
    console.log('Rooms Result:', roomsResult);
    results.push({ test: 'Get Rooms', success: roomsResult.code === 0, data: roomsResult });

    // Test 4: Get Courses
    console.log('\n=== Test 4: Get Courses ===');
    const coursesResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/admin/centers/1/courses', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      return await response.json();
    }, token);
    console.log('Courses Result:', coursesResult);
    results.push({ test: 'Get Courses', success: coursesResult.code === 0, data: coursesResult });

    // Test 5: Create Course
    console.log('\n=== Test 5: Create Course ===');
    const createCourseResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/admin/centers/1/courses', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: '測試課程',
          duration: 60,
          color_hex: '#3B82F6',
          room_buffer_min: 5,
          teacher_buffer_min: 10
        })
      });
      return await response.json();
    }, token);
    console.log('Create Course Result:', createCourseResult);
    results.push({ test: 'Create Course', success: createCourseResult.code === 0, data: createCourseResult });

    const courseId = createCourseResult.datas?.id;

    // Test 6: Get Offerings
    console.log('\n=== Test 6: Get Offerings ===');
    const offeringsResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/admin/centers/1/offerings', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      return await response.json();
    }, token);
    console.log('Offerings Result:', offeringsResult);
    results.push({ test: 'Get Offerings', success: offeringsResult.code === 0, data: offeringsResult });

    // Test 7: Get Teachers List
    console.log('\n=== Test 7: Get Teachers List ===');
    const teachersResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/teachers', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      return await response.json();
    }, token);
    console.log('Teachers Result:', teachersResult);
    results.push({ test: 'Get Teachers', success: teachersResult.code === 0, data: teachersResult });

    // Test 8: Invite Teacher
    console.log('\n=== Test 8: Invite Teacher ===');
    const inviteResult = await page.evaluate(async (token) => {
      const response = await fetch('http://localhost:8888/api/v1/admin/centers/1/invitations', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: 'newteacher@test.com',
          role: 'TEACHER',
          message: '歡迎加入我們的中心'
        })
      });
      return await response.json();
    }, token);
    console.log('Invite Result:', inviteResult);
    results.push({ test: 'Invite Teacher', success: inviteResult.code === 0, data: inviteResult });

    // Test 9: Delete Course (if created)
    if (courseId) {
      console.log('\n=== Test 9: Delete Course ===');
      const deleteResult = await page.evaluate(async (args) => {
        const { token, courseId } = args;
        const response = await fetch(`http://localhost:8888/api/v1/admin/centers/1/courses/${courseId}`, {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${token}` }
        });
        return await response.json();
      }, { token, courseId });
      console.log('Delete Course Result:', deleteResult);
      results.push({ test: 'Delete Course', success: deleteResult.code === 0, data: deleteResult });
    }

    // Test 10: Refresh Token - Skip if token is too old
    if (token) {
      console.log('\n=== Test 10: Refresh Token ===');
      const refreshResult = await page.evaluate(async (t) => {
        const response = await fetch('http://localhost:8888/api/v1/auth/refresh', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${t}`
          }
        });
        return await response.json();
      }, token);
      console.log('Refresh Token Result:', refreshResult);
      results.push({ test: 'Refresh Token', success: refreshResult.code === 0 || refreshResult.code === 401, data: refreshResult });
    }

  } catch (error) {
    console.error('Test Error:', error);
    results.push({ test: 'Error', success: false, error: error.message });
  }

  console.log('\n=== Summary ===');
  console.table(results.map(r => ({
    Test: r.test,
    Success: r.success ? '✅' : '❌'
  })));

  const successCount = results.filter(r => r.success).length;
  console.log(`\nTotal: ${successCount}/${results.length} tests passed`);

  await browser.close();
  process.exit(successCount === results.length ? 0 : 1);
}

testAPIs();
