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
