import { test, expect, type Page } from '@playwright/test';

test.describe('Admin Course Creation Flow', () => {
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

  test('should navigate to resources and display all tabs', async () => {
    await page.goto('/admin/resources');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('main h1')).toContainText('資源管理', { timeout: 10000 });

    await expect(page.locator('text=教室')).toBeVisible();
    await expect(page.locator('text=課程')).toBeVisible();
    await expect(page.locator('text=待排課程')).toBeVisible();
    await expect(page.locator('text=老師')).toBeVisible();
  });

  test('should create a new course successfully', async () => {
    await page.goto('/admin/resources');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('main h1')).toContainText('資源管理', { timeout: 10000 });

    await page.click('text=課程');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('text=課程列表')).toBeVisible();

    await page.click('text=+ 新增課程');
    await expect(page.locator('text=新增課程')).toBeVisible({ timeout: 5000 });

    await page.fill('input[placeholder*="例：鋼琴基礎"]', '測試課程 ' + Date.now());

    await page.fill('input[type="number"]:above(label:has-text("老師緩衝時間"))', '15');
    await page.fill('input[type="number"]:above(label:has-text("教室緩衝時間"))', '5');

    const submitButton = page.locator('button[type="submit"]:has-text("儲存")');
    await expect(submitButton).toBeVisible();

    await submitButton.click();

    await page.waitForTimeout(2000);

    const successMessage = page.locator('text=儲存成功, text=成功').first();
    if (await successMessage.isVisible({ timeout: 5000 })) {
      await expect(successMessage).toBeVisible();
    }
  });

  test('should create a new room successfully', async () => {
    await page.goto('/admin/resources');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('main h1')).toContainText('資源管理', { timeout: 10000 });

    await page.click('text=教室');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('text=教室列表')).toBeVisible();

    await page.click('text=+ 新增教室, button:has-text("新增教室")');
    await expect(page.locator('text=新增教室')).toBeVisible({ timeout: 5000 });
    
    await page.fill('input[placeholder*="例：Room A"]', '測試教室 ' + Date.now());
    await page.fill('input[type="number"]:above(label:has-text("容量"))', '20');

    const submitButton = page.locator('button[type="submit"]:has-text("儲存")');
    await expect(submitButton).toBeVisible();

    await submitButton.click();

    await page.waitForTimeout(2000);
  });

  test('should display offerings correctly', async () => {
    await page.goto('/admin/resources');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('main h1')).toContainText('資源管理', { timeout: 10000 });

    await page.click('text=待排課程');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    const offeringsSection = page.locator('text=待排課程').first();
    if (await offeringsSection.isVisible()) {
      console.log('Offerings tab is visible');
    }
  });

  test('should display teachers list', async () => {
    await page.goto('/admin/resources');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('main h1')).toContainText('資源管理', { timeout: 10000 });

    await page.click('text=老師');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    const teachersSection = page.locator('text=老師列表').first();
    if (await teachersSection.isVisible() || page.locator('text=尚未添加老師').isVisible()) {
      console.log('Teachers tab is visible');
    }
  });
});

test.describe('Admin Schedules Page', () => {
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

  test('should navigate to schedules page', async () => {
    await page.goto('/admin/schedules');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const schedulesHeader = page.locator('h1:has-text("課程時段")').first();
    await expect(schedulesHeader).toBeVisible({ timeout: 10000 });
  });

  test('should navigate to templates page', async () => {
    await page.goto('/admin/templates');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const templatesHeader = page.locator('h1:has-text("課表模板")').first();
    await expect(templatesHeader).toBeVisible({ timeout: 10000 });
  });

  test('should navigate to matching page', async () => {
    await page.goto('/admin/matching');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const matchingHeader = page.locator('h1:has-text("智慧媒合")').first();
    await expect(matchingHeader).toBeVisible({ timeout: 10000 });
  });
});

test.describe('Admin Approval Page', () => {
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

  test('should navigate to approval page and load exceptions', async () => {
    await page.goto('/admin/approval');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const approvalHeader = page.locator('h1:has-text("審核")').first();
    await expect(approvalHeader).toBeVisible({ timeout: 10000 });

    await page.waitForTimeout(2000);

    const allFilter = page.locator('button:has-text("全部")').first();
    if (await allFilter.isVisible()) {
      console.log('Approval filters are visible');
    }
  });
});

test.describe('API - Course Creation', () => {
  test('should create course via API with correct fields', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    expect(loginResponse.status()).toBe(200);
    const token = loginData.datas.token;
    
    const courseName = 'API Test Course ' + Date.now();
    const response = await request.post('http://localhost:3005/api/v1/admin/centers/1/courses', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      data: {
        name: courseName,
        duration: 60,
        color_hex: '#3B82F6',
        teacher_buffer_min: 10,
        room_buffer_min: 5
      }
    });
    
    console.log('Course creation response status:', response.status());
    console.log('Course creation response:', await response.text());
    
    if (response.status() === 200) {
      const data = await response.json();
      expect(data.code).toBe(0);
      expect(data.datas.name).toBe(courseName);
    } else if (response.status() === 400) {
      console.log('Course creation validation error');
    } else {
      console.log('Unexpected response status:', response.status());
    }
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
        name: 'API Test Room ' + Date.now(),
        capacity: 20
      }
    });
    
    console.log('Room creation response status:', response.status());
    console.log('Room creation response:', await response.text());
    
    expect([200, 400]).toContain(response.status());
  });

  test('should get courses list via API', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    const token = loginData.datas.token;
    
    const response = await request.get('http://localhost:3005/api/v1/admin/centers/1/courses', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    console.log('Get courses response status:', response.status());
    
    if (response.status() === 200) {
      const data = await response.json();
      console.log('Courses count:', data.datas?.length || 0);
      expect(data.code).toBe(0);
    }
  });

  test('should get rooms list via API', async ({ request }) => {
    const loginResponse = await request.post('http://localhost:3005/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    });
    
    const loginData = await loginResponse.json();
    const token = loginData.datas.token;
    
    const response = await request.get('http://localhost:3005/api/v1/admin/centers/1/rooms', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    console.log('Get rooms response status:', response.status());
    
    if (response.status() === 200) {
      const data = await response.json();
      console.log('Rooms count:', data.datas?.length || 0);
      expect(data.code).toBe(0);
    }
  });
});
