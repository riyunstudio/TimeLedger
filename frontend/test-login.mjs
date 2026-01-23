import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const errors = [];
  page.on('console', msg => {
    if (msg.type() === 'error') errors.push(msg.text());
  });
  page.on('pageerror', err => errors.push(err.message));

  // Go to home page and click mock login
  console.log('Going to home page...');
  await page.goto('http://localhost:3000/');
  await page.waitForTimeout(2000);

  // Click mock login button
  console.log('Clicking mock login...');
  const mockLoginBtn = await page.$('button:has-text("Mock 教師登入")');
  if (mockLoginBtn) {
    await mockLoginBtn.click();
    await page.waitForTimeout(2000);
  } else {
    console.log('Mock login button not found');
  }

  // Navigate to dashboard
  console.log('Navigating to dashboard...');
  await page.goto('http://localhost:3000/teacher/dashboard');
  await page.waitForTimeout(3000);

  // Check for loading text
  const loadingText = await page.textContent('body');
  console.log('Page contains 載入中:', loadingText.includes('載入中'));

  // Check for schedule grid
  const gridExists = await page.$('.glass-card');
  console.log('Grid exists:', !!gridExists);

  console.log('\n=== Console Errors ===');
  console.log(errors.length > 0 ? errors.join('\n') : 'None');
  console.log('======================\n');

  await page.screenshot({ path: 'C:/tmp/after-login.png' });
  console.log('Screenshot saved');

  await browser.close();
})();
