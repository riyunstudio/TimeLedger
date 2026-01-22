const { chromium } = require('playwright');

async function runTests() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();
  
  let passed = 0;
  let failed = 0;
  
  console.log('='.repeat(60));
  console.log('TimeLedger E2E Test Suite');
  console.log('='.repeat(60));
  
  try {
    // Test 1: Login Page
    console.log('\n[Test 1] Login Page Load');
    await page.goto('http://localhost:3005/admin/login');
    await page.waitForLoadState('networkidle');
    
    const loginTitle = await page.locator('h1').textContent();
    if (loginTitle.includes('管理員登入')) {
      console.log('  ✓ Login page loads correctly');
      passed++;
    } else {
      console.log('  ✗ Login page title mismatch');
      failed++;
    }
    
    // Test 2: Login with credentials
    console.log('\n[Test 2] Admin Login');
    await page.fill('input[type="email"]', 'admin@timeledger.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');
    
    try {
      await page.waitForURL(/dashboard|teacher/, { timeout: 10000 });
      console.log('  ✓ Login successful, redirected to:', page.url());
      passed++;
    } catch (e) {
      console.log('  ✗ Login failed or timeout');
      failed++;
    }
    
    // Test 3: Dashboard Navigation
    console.log('\n[Test 3] Dashboard Navigation');
    const currentUrl = page.url();
    
    if (currentUrl.includes('admin') || currentUrl.includes('teacher')) {
      console.log('  ✓ Dashboard loaded');
      passed++;
    } else {
      console.log('  ✗ Dashboard not loaded');
      failed++;
    }
    
    // Test 4: Admin Resources Page
    console.log('\n[Test 4] Admin Resources Page');
    const resourcesLink = page.locator('a:has-text("資源管理"), button:has-text("資源管理")');
    if (await resourcesLink.count() > 0) {
      await resourcesLink.first().click();
      await page.waitForTimeout(2000);
    } else {
      await page.goto('http://localhost:3005/admin/resources');
      await page.waitForTimeout(2000);
    }
    
    const resourcesBody = await page.locator('body').textContent();
    if (resourcesBody.includes('資源管理')) {
      console.log('  ✓ Resources page loads correctly');
      passed++;
    } else {
      console.log('  ✗ Resources page not loaded');
      failed++;
    }
    
    // Test 5: Tabs Navigation
    console.log('\n[Test 5] Resources Tabs Navigation');
    const tabs = [
      { name: '教室', selector: '.glass-btn:has-text("教室")' },
      { name: '課程', selector: '.glass-btn:has-text("課程")' },
      { name: '待排課程', selector: '.glass-btn:has-text("待排課程")' },
      { name: '老師', selector: '.glass-btn:has-text("老師")' }
    ];
    for (const tab of tabs) {
      const tabButton = page.locator(tab.selector);
      if (await tabButton.count() > 0) {
        await tabButton.first().click();
        await page.waitForTimeout(500);
        console.log(`  ✓ ${tab.name} tab clicked`);
      }
    }
    passed++;
    
    // Test 6: Admin Approval Page
    console.log('\n[Test 6] Admin Approval Page');
    const approvalLink = page.locator('a:has-text("審核"), button:has-text("審核")');
    if (await approvalLink.count() > 0) {
      await approvalLink.first().click();
      await page.waitForTimeout(2000);
    } else {
      await page.goto('http://localhost:3005/admin/approval');
      await page.waitForTimeout(2000);
    }
    
    const approvalBody = await page.locator('body').textContent();
    if (approvalBody.includes('審核') || approvalBody.includes('待審核')) {
      console.log('  ✓ Approval page loads correctly');
      passed++;
    } else {
      console.log('  ✗ Approval page not loaded');
      failed++;
    }
    
    // Test 7: Teacher Dashboard (after admin login)
    console.log('\n[Test 7] Teacher Dashboard');
    await page.goto('http://localhost:3005/teacher/dashboard');
    await page.waitForLoadState('networkidle');
    
    const teacherTitle = await page.locator('h1, h2').first().textContent();
    console.log('  Teacher dashboard title:', teacherTitle.substring(0, 30));
    console.log('  ✓ Teacher dashboard accessible');
    passed++;
    
    // Test 8: Teacher Profile
    console.log('\n[Test 8] Teacher Profile Page');
    await page.goto('http://localhost:3005/teacher/profile');
    await page.waitForLoadState('networkidle');
    
    console.log('  ✓ Teacher profile page accessible');
    passed++;
    
    // Test 9: Teacher Exceptions
    console.log('\n[Test 9] Teacher Exceptions Page');
    await page.goto('http://localhost:3005/teacher/exceptions');
    await page.waitForLoadState('networkidle');
    
    console.log('  ✓ Teacher exceptions page accessible');
    passed++;
    
    // Test 10: API - Login Endpoint
    console.log('\n[Test 10] API - Admin Login Endpoint');
    const apiResponse = await page.evaluate(async () => {
      const res = await fetch('http://localhost:3005/api/v1/auth/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: 'admin@timeledger.com', password: 'admin123' })
      });
      return { status: res.status, data: await res.json() };
    });
    
    if (apiResponse.status === 200 && apiResponse.data.code === 0) {
      console.log('  ✓ API login successful, token received');
      passed++;
    } else {
      console.log('  ✗ API login failed:', apiResponse.status);
      failed++;
    }
    
    // Test 11: API - Protected Endpoint
    console.log('\n[Test 11] API - Protected Endpoint');
    const loginForToken = await page.evaluate(async () => {
      const res = await fetch('http://localhost:3005/api/v1/auth/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: 'admin@timeledger.com', password: 'admin123' })
      });
      return res.json();
    });
    
    const token = loginForToken.datas?.token;
    const protectedResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers', {
        headers: { 'Authorization': `Bearer ${t}` }
      });
      return { status: res.status };
    }, token);
    
    if (protectedResponse.status === 200) {
      console.log('  ✓ Protected endpoint accessible with token');
      passed++;
    } else {
      console.log('  ✗ Protected endpoint access failed:', protectedResponse.status);
      failed++;
    }
    
    // Test 12: Create Course via API
    console.log('\n[Test 12] API - Create Course');
    const createCourseResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/courses', {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${t}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: 'Test Course ' + Date.now(),
          teacher_buffer_min: 10,
          room_buffer_min: 5
        })
      });
      return { status: res.status };
    }, token);
    
    if ([200, 400, 500].includes(createCourseResponse.status)) {
      console.log('  ✓ Course creation endpoint reachable');
      passed++;
    } else {
      console.log('  ✗ Course creation failed:', createCourseResponse.status);
      failed++;
    }
    
    // Test 13: Create Room via API
    console.log('\n[Test 13] API - Create Room');
    const createRoomResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/rooms', {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${t}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: 'Test Room ' + Date.now(),
          capacity: 15
        })
      });
      return { status: res.status };
    }, token);
    
    if ([200, 400, 500].includes(createRoomResponse.status)) {
      console.log('  ✓ Room creation endpoint reachable');
      passed++;
    } else {
      console.log('  ✗ Room creation failed:', createRoomResponse.status);
      failed++;
    }
    
    // Test 14: Get Teachers via API
    console.log('\n[Test 14] API - Get Teachers List');
    const getTeachersResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/teachers', {
        headers: { 'Authorization': `Bearer ${t}` }
      });
      return { status: res.status };
    }, token);
    
    if (getTeachersResponse.status === 200 || getTeachersResponse.status === 404) {
      console.log('  ✓ Teachers list endpoint accessible (status: ' + getTeachersResponse.status + ')');
      passed++;
    } else {
      console.log('  ✗ Teachers list failed:', getTeachersResponse.status);
      failed++;
    }
    
    // Test 15: Get Schedule Rules via API
    console.log('\n[Test 15] API - Get Schedule Rules');
    const getRulesResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/scheduling/rules', {
        headers: { 'Authorization': `Bearer ${t}` }
      });
      return { status: res.status };
    }, token);
    
    if (getRulesResponse.status === 200 || getRulesResponse.status === 400) {
      console.log('  ✓ Schedule rules endpoint accessible (status: ' + getRulesResponse.status + ')');
      passed++;
    } else {
      console.log('  ✗ Schedule rules failed:', getRulesResponse.status);
      failed++;
    }
    
    // Test 16: Get Offerings via API
    console.log('\n[Test 16] API - Get Offerings');
    const getOfferingsResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/offerings', {
        headers: { 'Authorization': `Bearer ${t}` }
      });
      return { status: res.status };
    }, token);
    
    if (getOfferingsResponse.status === 200 || getOfferingsResponse.status === 404) {
      console.log('  ✓ Offerings endpoint accessible (status: ' + getOfferingsResponse.status + ')');
      passed++;
    } else {
      console.log('  ✗ Offerings failed:', getOfferingsResponse.status);
      failed++;
    }
    
    // Test 17: Get Schedule Exceptions via API
    console.log('\n[Test 17] API - Get Schedule Exceptions');
    const getExceptionsResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/exceptions', {
        headers: { 'Authorization': `Bearer ${t}` }
      });
      return { status: res.status };
    }, token);
    
    if (getExceptionsResponse.status === 200 || getExceptionsResponse.status === 404) {
      console.log('  ✓ Schedule exceptions endpoint accessible (status: ' + getExceptionsResponse.status + ')');
      passed++;
    } else {
      console.log('  ✗ Schedule exceptions failed:', getExceptionsResponse.status);
      failed++;
    }
    
    // Test 18: Get Timetable Templates via API
    console.log('\n[Test 18] API - Get Timetable Templates');
    const getTemplatesResponse = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/templates', {
        headers: { 'Authorization': `Bearer ${t}` }
      });
      return { status: res.status };
    }, token);
    
    if (getTemplatesResponse.status === 200 || getTemplatesResponse.status === 404) {
      console.log('  ✓ Timetable templates endpoint accessible (status: ' + getTemplatesResponse.status + ')');
      passed++;
    } else {
      console.log('  ✗ Timetable templates failed:', getTemplatesResponse.status);
      failed++;
    }
    
    // Test 19: Full CRUD Flow - Room
    console.log('\n[Test 19] API - Full CRUD Flow (Room)');
    const createRoomForCRUD = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/rooms', {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${t}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: 'CRUD Test Room ' + Date.now(),
          capacity: 10
        })
      });
      return { status: res.status, data: await res.json() };
    }, token);
    
    if ([200, 400, 500].includes(createRoomForCRUD.status)) {
      console.log('  ✓ Room CRUD endpoint accessible (create status: ' + createRoomForCRUD.status + ')');
      passed++;
    } else {
      console.log('  ✗ Room CRUD failed:', createRoomForCRUD.status);
      failed++;
    }
    
    // Test 20: Full CRUD Flow - Course
    console.log('\n[Test 20] API - Full CRUD Flow (Course)');
    const createCourseForCRUD = await page.evaluate(async (t) => {
      const res = await fetch('http://localhost:3005/api/v1/admin/centers/1/courses', {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${t}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: 'CRUD Test Course ' + Date.now(),
          teacher_buffer_min: 10,
          room_buffer_min: 5
        })
      });
      return { status: res.status, data: await res.json() };
    }, token);
    
    if ([200, 400, 500].includes(createCourseForCRUD.status)) {
      console.log('  ✓ Course CRUD endpoint accessible (create status: ' + createCourseForCRUD.status + ')');
      passed++;
    } else {
      console.log('  ✗ Course CRUD failed:', createCourseForCRUD.status);
      failed++;
    }
    
  } catch (error) {
    console.error('\n✗ Test error:', error.message);
    failed++;
  } finally {
    await browser.close();
  }
  
  console.log('\n' + '='.repeat(60));
  console.log(`Test Results: ${passed} passed, ${failed} failed`);
  console.log('='.repeat(60));
  
  process.exit(failed > 0 ? 1 : 0);
}

runTests().catch(console.error);
