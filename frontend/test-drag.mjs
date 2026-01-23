import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const logs = [];
  const errors = [];

  page.on('console', msg => {
    const text = msg.text();
    if (msg.type() === 'error') {
      errors.push(text);
    } else {
      logs.push(`[${msg.type()}] ${text}`);
    }
  });
  page.on('pageerror', err => errors.push(err.message));

  console.log('Navigating to dashboard...');
  await page.goto('http://localhost:3000/teacher/dashboard');
  await page.waitForTimeout(3000);

  console.log('Checking page title:', await page.title());

  // Check if schedule grid exists
  const gridExists = await page.$('.glass-card');
  console.log('Grid exists:', !!gridExists);

  // Check for any draggable elements
  const draggableItems = await page.$$('[draggable="true"]');
  console.log('Draggable items found:', draggableItems.length);

  // Log any errors
  console.log('\n=== Console Errors ===');
  console.log(errors.length > 0 ? errors.join('\n') : 'None');
  console.log('======================\n');

  // Log first 10 console messages
  console.log('=== Console Logs (first 10) ===');
  logs.slice(0, 10).forEach(log => console.log(log));
  console.log('======================\n');

  await page.screenshot({ path: 'C:/tmp/dashboard2.png' });
  console.log('Screenshot saved to C:/tmp/dashboard2.png');

  await browser.close();
})();
