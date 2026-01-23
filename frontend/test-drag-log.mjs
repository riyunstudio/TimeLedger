import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const logs = [];
  page.on('console', msg => {
    logs.push(`[${msg.type()}] ${msg.text()}`);
  });
  page.on('pageerror', err => logs.push(`[error] ${err.message}`));

  console.log('Starting test...');
  await page.goto('http://localhost:3000/');
  await page.waitForTimeout(1000);

  // Mock login
  const mockLoginBtn = await page.$('button:has-text("Mock 教師登入")');
  if (mockLoginBtn) {
    await mockLoginBtn.click();
    await page.waitForTimeout(2000);
  }

  // Go to dashboard
  await page.goto('http://localhost:3000/teacher/dashboard');
  await page.waitForTimeout(3000);

  // Get draggable items
  const draggableItems = await page.$$('[draggable="true"]');
  console.log(`Found ${draggableItems.length} draggable items`);

  if (draggableItems.length > 0) {
    const item = draggableItems[0];
    const itemBox = await item.boundingBox();

    const cells = await page.$$('.border-t.border-l');
    if (cells.length > 1 && itemBox) {
      const targetCell = cells[1];
      const targetBox = await targetCell.boundingBox();

      if (targetBox) {
        // Perform drag
        await page.mouse.move(itemBox.x + itemBox.width / 2, itemBox.y + itemBox.height / 2);
        await page.mouse.down();
        await page.mouse.move(targetBox.x + targetBox.width / 2, targetBox.y + targetBox.height / 2, { steps: 10 });
        await page.mouse.up();

        console.log('Drag completed, waiting 3 seconds...');
        await page.waitForTimeout(3000);
      }
    }
  }

  // Filter logs for API-related messages
  console.log('\n=== API-related logs ===');
  const apiLogs = logs.filter(log =>
    log.includes('moveScheduleItem') ||
    log.includes('Failed to move schedule') ||
    log.includes('更新失敗') ||
    log.includes('POST') ||
    log.includes('patch') ||
    log.includes('teacher/schedule') ||
    log.includes('personal-events')
  );
  apiLogs.forEach(log => console.log(log));
  console.log('========================\n');

  await browser.close();
})();
