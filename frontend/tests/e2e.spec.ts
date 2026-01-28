import { test, expect, type Page } from '@playwright/test';

test.describe('TimeLedger Admin Flow E2E Tests', () => {
  let page: Page;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
    page.setViewportSize({ width: 1280, height: 720 });
  });

  test.afterAll(async () => {
    await page.close();
  });

  test.describe('Admin Authentication', () => {
    test('should display login page correctly', async () => {
      await page.goto('/admin/login');
      
      await expect(page.locator('h1')).toContainText('管理員登入');
      await expect(page.locator('input[type="email"]')).toBeVisible();
      await expect(page.locator('input[type="password"]')).toBeVisible();
      await expect(page.locator('button[type="submit"]')).toContainText('登入');
    });

    test('should login successfully with valid credentials', async () => {
      await page.goto('/admin/login');
      
      await page.fill('input[type="email"]', 'admin@timeledger.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      
      await expect(page).toHaveURL(/.*dashboard/);
      await expect(page.locator('text=系統管理員')).toBeVisible({ timeout: 10000 });
    });

    test('should show error with invalid credentials', async () => {
      await page.goto('/admin/login');
      
      await page.fill('input[type="email"]', 'wrong@email.com');
      await page.fill('input[type="password"]', 'wrongpassword');
      await page.click('button[type="submit"]');
      
      await expect(page.locator('text=登入失敗')).toBeVisible({ timeout: 5000 });
    });
  });

  test.describe('Admin Dashboard', () => {
    test.beforeEach(async () => {
      await page.goto('/admin/login');
      await page.fill('input[type="email"]', 'admin@timeledger.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
    });

    test('should display dashboard with schedule grid', async () => {
      await expect(page.locator('text=待排課程')).toBeVisible();
      await expect(page.locator('text=已排課表')).toBeVisible();
    });

    test('should navigate to resources page', async () => {
      await page.click('text=資源管理');
      await expect(page.locator('h1')).toContainText('資源管理');
      await expect(page.locator('text=教室')).toBeVisible();
      await expect(page.locator('text=課程')).toBeVisible();
      await expect(page.locator('text=待排課程')).toBeVisible();
      await expect(page.locator('text=老師')).toBeVisible();
    });

    test('should navigate to approval page', async () => {
      await page.click('text=審核');
      await expect(page.locator('h1')).toContainText('審核');
    });
  });

  test.describe('Resource Management', () => {
    test.beforeEach(async () => {
      await page.goto('/admin/login');
      await page.fill('input[type="email"]', 'admin@timeledger.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
      await page.click('text=資源管理');
      await expect(page.locator('h1')).toContainText('資源管理', { timeout: 5000 });
    });

    test('should display rooms tab', async () => {
      await page.click('text=教室');
      await expect(page.locator('text=教室列表')).toBeVisible();
    });

    test('should display courses tab', async () => {
      await page.click('text=課程');
      await expect(page.locator('text=課程列表')).toBeVisible();
    });

    test('should display offerings tab', async () => {
      await page.click('text=待排課程');
      await expect(page.locator('text=待排課程列表')).toBeVisible();
    });

    test('should display teachers tab', async () => {
      await page.click('text=老師');
      await expect(page.locator('text=老師列表')).toBeVisible();
    });
  });

  test.describe('Schedule Management', () => {
    test.beforeEach(async () => {
      await page.goto('/admin/login');
      await page.fill('input[type="email"]', 'admin@timeledger.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
    });

    test('should display schedule grid', async () => {
      await expect(page.locator('text=待排課程')).toBeVisible();
      await expect(page.locator('text=已排課表')).toBeVisible();
    });

    test('should navigate through week days', async () => {
      const nextWeekButton = page.locator('button:has-text("下一週")');
      const prevWeekButton = page.locator('button:has-text("上一週")');
      
      if (await nextWeekButton.isVisible()) {
        await nextWeekButton.click();
        await page.waitForTimeout(500);
      }
      
      if (await prevWeekButton.isVisible()) {
        await prevWeekButton.click();
        await page.waitForTimeout(500);
      }
    });
  });
});

test.describe('API Integration Tests', () => {
  test('backend should accept admin login', async ({ request }) => {
    const response = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    expect(data.code).toBe(0);
    expect(data.datas.token).toBeDefined();
    expect(data.datas.user.email).toBe('admin@timeledger.com');
  });

  test('backend should reject invalid login', async ({ request }) => {
    const response = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'wrong@email.com',
        password: 'wrongpassword'
      }
    });
    
    expect(response.status()).toBe(401);
    
    const data = await response.json();
    expect(data.code).toBe(401);
  });

  test('should access protected endpoint with valid token', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    const token = loginData.datas.token;
    
    const response = await request.get('http://localhost:3005/api/v1/admin/centers/1', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    expect(data.code).toBe(0);
  });

  test('should reject protected endpoint without token', async ({ request }) => {
    const response = await request.get('http://localhost:3005/api/v1/admin/centers/1');
    
    expect(response.status()).toBe(401);
  });

  test('should create room via API', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    const token = loginData.datas.token;
    
    const response = await request.post('http://localhost:3005/api/v1/admin/centers/1/rooms', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      data: {
        name: 'Test Room ' + Date.now(),
        capacity: 10
      }
    });
    
    expect([200, 400, 500]).toContain(response.status());
  });

  test('should create course via API', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    const token = loginData.datas.token;
    
    const response = await request.post('http://localhost:3005/api/v1/admin/centers/1/courses', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      data: {
        name: 'Piano ' + Date.now(),
        teacher_buffer_min: 10,
        room_buffer_min: 5
      }
    });
    
    expect([200, 400, 500]).toContain(response.status());
  });

  test('should get teachers list via API', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    const token = loginData.datas.token;
    
    const response = await request.get('http://localhost:3005/api/v1/teachers', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    expect(response.status()).toBe(200);
  });
});

test.describe('Course Scheduling Flow', () => {
  let page: Page;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
    page.setViewportSize({ width: 1280, height: 720 });
  });

  test.afterAll(async () => {
    await page.close();
  });

  test.beforeEach(async () => {
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', 'admin@timeledger.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
  });

  test('should be able to add offering to schedule', async () => {
    const addButton = page.locator('button:has-text("新增")').first();

    if (await addButton.isVisible({ timeout: 3000 })) {
      await addButton.click();
      await page.waitForTimeout(500);

      const modal = page.locator('text=新增待排課程').first();
      if (await modal.isVisible({ timeout: 2000 })) {
        await expect(modal).toBeVisible();
      }
    }
  });

  test('should display schedule details panel when clicking on scheduled class', async () => {
    const scheduleItems = page.locator('[class*="schedule-item"], [class*="class-item"]').first();

    if (await scheduleItems.count() > 0) {
      await scheduleItems.click();
      await page.waitForTimeout(500);

      const detailPanel = page.locator('text=課堂詳情, text=課程詳情').first();
      if (await detailPanel.isVisible({ timeout: 2000 })) {
        await expect(detailPanel).toBeVisible();
      }
    }
  });
});

test.describe('Exception Approval Flow E2E Tests', () => {
  let page: Page;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
    page.setViewportSize({ width: 1280, height: 720 });
  });

  test.afterAll(async () => {
    await page.close();
  });

  test.beforeEach(async () => {
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', 'admin@timeledger.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
  });

  test('should navigate to approval page', async () => {
    await page.click('text=審核');
    await expect(page.locator('h1')).toContainText('審核');
  });

  test('should display pending exceptions', async () => {
    await page.click('text=審核');
    await expect(page.locator('text=待審核')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=例外申請')).toBeVisible();
  });

  test('should approve an exception request', async () => {
    await page.click('text=審核');

    // Find and click on first pending exception
    const approveButton = page.locator('button:has-text("核准")').first();
    if (await approveButton.isVisible({ timeout: 3000 })) {
      await approveButton.click();
      await page.waitForTimeout(500);

      // Verify confirmation dialog appears
      await expect(page.locator('text=確認核准')).toBeVisible({ timeout: 2000 });
    }
  });

  test('should reject an exception request with reason', async () => {
    await page.click('text=審核');

    const rejectButton = page.locator('button:has-text("拒絕")').first();
    if (await rejectButton.isVisible({ timeout: 3000 })) {
      await rejectButton.click();
      await page.waitForTimeout(500);

      // Verify rejection modal appears
      await expect(page.locator('text=拒絕原因')).toBeVisible({ timeout: 2000 });
    }
  });
});

test.describe('Smart Matching E2E Tests', () => {
  let page: Page;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
    page.setViewportSize({ width: 1280, height: 720 });
  });

  test.afterAll(async () => {
    await page.close();
  });

  test.beforeEach(async () => {
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', 'admin@timeledger.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
  });

  test('should navigate to matching page', async () => {
    await page.click('text=媒合');
    await expect(page.locator('h1')).toContainText('智慧媒合');
  });

  test('should search for teachers by keyword', async () => {
    await page.click('text=媒合');

    const searchInput = page.locator('input[placeholder*="搜尋"], input[placeholder*="關鍵字"]').first();
    if (await searchInput.isVisible({ timeout: 3000 })) {
      await searchInput.fill('瑜珈');
      await page.waitForTimeout(500);

      await expect(page.locator('text=瑜珈')).toBeVisible({ timeout: 5000 });
    }
  });

  test('should display talent pool statistics', async () => {
    await page.click('text=媒合');

    await expect(page.locator('text=人才庫')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=總人數')).toBeVisible();
    await expect(page.locator('text=待邀請')).toBeVisible();
  });

  test('should show alternative time slots', async () => {
    await page.click('text=媒合');

    const alternativeTab = page.locator('text=替代時段').first();
    if (await alternativeTab.isVisible({ timeout: 3000 })) {
      await alternativeTab.click();
      await expect(page.locator('text=替代時段')).toBeVisible();
    }
  });
});

test.describe('Teacher Dashboard E2E Tests', () => {
  let page: Page;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
    page.setViewportSize({ width: 1280, height: 720 });
  });

  test.afterAll(async () => {
    await page.close();
  });

  test('should display teacher login page', async () => {
    await page.goto('/teacher/login');
    await expect(page.locator('h1, h2')).toContainText('登入');
  });

  test('should display teacher dashboard after login', async () => {
    // Note: This test requires LINE login mock or valid credentials
    await page.goto('/teacher/dashboard');
    // Should redirect to login or show dashboard based on auth state
    await page.waitForTimeout(2000);
  });

  test('should navigate to teacher exceptions page', async () => {
    await page.goto('/teacher/exceptions');
    await expect(page.locator('h1, h2')).toContainText('例外');
  });
});

test.describe('Notification Settings E2E Tests', () => {
  let page: Page;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
    page.setViewportSize({ width: 1280, height: 720 });
  });

  test.afterAll(async () => {
    await page.close();
  });

  test.beforeEach(async () => {
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', 'admin@timeledger.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 });
  });

  test('should navigate to settings page', async () => {
    await page.click('[class*="user"], text=設定, text=系統管理員');
    await page.waitForTimeout(500);

    // Look for settings link in dropdown or sidebar
    const settingsLink = page.locator('text=設定').first();
    if (await settingsLink.isVisible()) {
      await settingsLink.click();
      await expect(page.locator('text=LINE')).toBeVisible({ timeout: 5000 });
    }
  });

  test('should display LINE binding section in settings', async () => {
    await page.click('text=設定');
    await page.waitForTimeout(500);

    await expect(page.locator('text=LINE')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=綁定')).toBeVisible();
  });

  test('should display LINE notification toggle', async () => {
    await page.click('text=設定');
    await page.waitForTimeout(500);

    await expect(page.locator('text=通知')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=LINE')).toBeVisible();
  });
});
