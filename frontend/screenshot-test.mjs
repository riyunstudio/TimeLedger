import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const errors = [];
  page.on('console', msg => {
    if (msg.type() === 'error') errors.push(msg.text());
  });
  page.on('pageerror', err => errors.push(err.message));

  await page.goto('http://localhost:3001/teacher/dashboard');
  await page.waitForTimeout(3000);

  await page.screenshot({ path: 'C:/tmp/dashboard.png' });

  console.log('=== Console Errors ===');
  console.log(errors.length > 0 ? errors.join('\n') : 'None');
  console.log('======================');
  console.log('Screenshot saved to C:/tmp/dashboard.png');

  await browser.close();
})();
