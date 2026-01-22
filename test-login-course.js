const { chromium } = require('@playwright/test');

(async () => {
  console.log('啟動瀏覽器...');
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();

  page.on('console', msg => console.log(`[Console ${msg.type()}] ${msg.text()}`));
  page.on('pageerror', err => console.log(`[Page Error] ${err.message}`));

  try {
    // 步驟 1: 前往登入頁面
    console.log('\n=== 步驟 1: 前往登入頁面 ===');
    await page.goto('http://localhost:3000/admin/login', { waitUntil: 'networkidle' });
    console.log(`URL: ${page.url()}`);
    console.log(`標題: ${await page.title()}`);
    await page.screenshot({ path: 'D:/project/TimeLedger/test-results/01-login-page.png' });

    // 步驟 2: 填寫登入資訊
    console.log('\n=== 步驟 2: 填寫登入資訊 ===');
    await page.waitForSelector('input[type="email"], input[name*="email"], input[name*="account"]');
    await page.fill('input[type="email"], input[name*="email"], input[name*="account"]', 'admin@timeledger.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.screenshot({ path: 'D:/project/TimeLedger/test-results/02-login-filled.png' });
    console.log('已填寫登入資訊');

    // 步驟 3: 點擊登入
    console.log('\n=== 步驟 3: 登入 ===');
    await page.click('button[type="submit"]');
    await page.waitForTimeout(3000);
    console.log(`登入後 URL: ${page.url()}`);
    await page.screenshot({ path: 'D:/project/TimeLedger/test-results/03-after-login.png' });

    // 步驟 4: 前往資源管理頁面
    console.log('\n=== 步驟 4: 前往資源管理頁面 ===');
    await page.goto('http://localhost:3000/admin/resources', { waitUntil: 'networkidle' });
    console.log(`URL: ${page.url()}`);
    await page.screenshot({ path: 'D:/project/TimeLedger/test-results/04-resources-page.png' });

    // 步驟 5: 點擊課程分頁
    console.log('\n=== 步驟 5: 點擊課程分頁 ===');
    const coursesTab = page.locator('button:has-text("課程"):not(:has-text("待排"))');
    if (await coursesTab.isVisible()) {
      await coursesTab.click();
      await page.waitForTimeout(1000);
      console.log('已點擊課程分頁');
      await page.screenshot({ path: 'D:/project/TimeLedger/test-results/05-courses-tab.png' });
    }

    // 步驟 6: 新增課程
    console.log('\n=== 步驟 6: 新增課程 ===');
    const addBtn = page.locator('button:has-text("新增")');
    if (await addBtn.first().isVisible()) {
      await addBtn.first().click();
      console.log('已點擊新增按鈕');
      await page.waitForTimeout(1000);
      await page.screenshot({ path: 'D:/project/TimeLedger/test-results/06-add-course-form.png' });
    }

    // 步驟 7: 填寫課程資料
    console.log('\n=== 步驟 7: 填寫課程資料 ===');
    const nameInput = page.locator('input[name*="name"], input[placeholder*="名稱"], input[placeholder*="名子"]');
    if (await nameInput.isVisible()) {
      await nameInput.fill('測試課程');
      console.log('已填寫課程名稱');
      await page.screenshot({ path: 'D:/project/TimeLedger/test-results/07-course-filled.png' });
    }

    console.log('\n=== 測試完成 ===');
    console.log(`最終 URL: ${page.url()}`);

  } catch (error) {
    console.error('錯誤:', error.message);
  } finally {
    await page.waitForTimeout(2000);
    console.log('瀏覽器保持開啟中...');
  }
})();
