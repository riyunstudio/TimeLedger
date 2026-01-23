import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  const errors = [];
  const apiCalls = [];

  page.on('console', msg => {
    const text = msg.text();
    if (msg.type() === 'error') {
      errors.push(text);
    } else if (text.includes('POST') || text.includes('moveScheduleItem') || text.includes('fetchSchedule')) {
      apiCalls.push(text);
    }
  });
  page.on('pageerror', err => errors.push(err.message));

  console.log('1. Going to home page...');
  await page.goto('http://localhost:3000/');
  await page.waitForTimeout(1000);

  console.log('2. Clicking mock login...');
  const mockLoginBtn = await page.$('button:has-text("Mock 教師登入")');
  if (mockLoginBtn) {
    await mockLoginBtn.click();
    await page.waitForTimeout(2000);
  }

  console.log('3. Navigating to dashboard...');
  await page.goto('http://localhost:3000/teacher/dashboard');
  await page.waitForTimeout(3000);

  console.log('4. Checking for draggable items...');
  const draggableItems = await page.$$('[draggable="true"]');
  console.log('   Draggable items found:', draggableItems.length);

  if (draggableItems.length > 0) {
    console.log('5. Testing drag and drop...');

    // Get the first draggable item
    const item = draggableItems[0];
    const itemBox = await item.boundingBox();

    // Get a target cell (first empty cell)
    const cells = await page.$$('.border-t.border-l');
    if (cells.length > 1 && itemBox) {
      const targetCell = cells[1];
      const targetBox = await targetCell.boundingBox();

      if (targetBox) {
        // Perform drag and drop
        await page.mouse.move(itemBox.x + itemBox.width / 2, itemBox.y + itemBox.height / 2);
        await page.mouse.down();
        await page.mouse.move(targetBox.x + targetBox.width / 2, targetBox.y + targetBox.height / 2, { steps: 10 });
        await page.mouse.up();

        console.log('   Drag completed, waiting for API call...');
        await page.waitForTimeout(2000);
      }
    }
  }

  console.log('\n=== Console Errors ===');
  console.log(errors.length > 0 ? errors.join('\n') : 'None');
  console.log('======================\n');

  await page.screenshot({ path: 'C:/tmp/drag-test.png' });
  console.log('Screenshot saved to C:/tmp/drag-test.png');

  await browser.close();
})();
