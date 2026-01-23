import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const logs = [];
  page.on('console', msg => {
    logs.push(`[${msg.type()}] ${msg.text()}`);
  });

  console.log('1. Going to home page...');
  await page.goto('http://localhost:3000/');
  await page.waitForTimeout(2000);

  // Check home page
  const homeTitle = await page.textContent('h1');
  console.log('   Home page title:', homeTitle);

  console.log('2. Clicking mock login...');
  const mockLoginBtn = await page.$('button:has-text("Mock 教師登入")');
  if (mockLoginBtn) {
    await mockLoginBtn.click();
    await page.waitForTimeout(2000);
  }

  // Check current URL
  console.log('3. Current URL:', page.url());

  console.log('4. Checking dashboard...');
  await page.goto('http://localhost:3000/teacher/dashboard');
  await page.waitForTimeout(3000);

  // Check page content
  const pageContent = await page.content();
  console.log('   Has grid class:', pageContent.includes('grid'));
  console.log('   Has schedule items:', pageContent.includes('鋼琴基礎') || pageContent.includes('小提琴'));

  // Check for draggable
  const draggableItems = await page.$$('[draggable="true"]');
  console.log('   Draggable items:', draggableItems.length);

  // Check for debug logs
  console.log('\n=== All console logs ===');
  logs.forEach(log => console.log(log));
  console.log('========================\n');

  await browser.close();
})();
