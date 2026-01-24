const { chromium } = require('playwright');

async function testScheduleRuleModal() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();

  // Collect ALL console messages
  const consoleLogs = [];
  const jsErrors = [];

  page.on('console', msg => {
    consoleLogs.push({ type: msg.type(), text: msg.text(), location: msg.location() });
  });

  page.on('pageerror', err => {
    jsErrors.push(err.message);
  });

  console.log('='.repeat(60));
  console.log('ScheduleRuleModal Debug Test');
  console.log('='.repeat(60));

  try {
    console.log('\n[Step 1] Navigating to Admin Dashboard...');
    await page.goto('http://localhost:3000/admin', { timeout: 30000, waitUntil: 'domcontentloaded' });

    console.log('\n[Step 2] Waiting for page to render...');
    await page.waitForTimeout(5000);

    // Print ALL console messages
    console.log('\n[Step 3] All Console Messages:');
    consoleLogs.forEach(log => {
      console.log(`  [${log.type}] ${log.text.substring(0, 200)}`);
    });

    // Print JS errors
    console.log('\n[Step 4] JavaScript Errors:');
    if (jsErrors.length > 0) {
      jsErrors.forEach((err, i) => console.log(`  ${i + 1}. ${err}`));
    } else {
      console.log('  No JavaScript errors detected');
    }

    // Check Nuxt state
    console.log('\n[Step 5] Nuxt State:');
    const nuxtState = await page.evaluate(() => {
      const nuxtDiv = document.getElementById('__nuxt');
      return {
        exists: !!nuxtDiv,
        innerHTML: nuxtDiv ? nuxtDiv.innerHTML.substring(0, 300) : 'No __nuxt div',
        childrenCount: nuxtDiv ? nuxtDiv.children.length : 0,
        hasNuxtData: !!window.__NUXT__,
        nuxtDataKeys: window.__NUXT__ ? Object.keys(window.__NUXT__).length : 0,
      };
    });
    console.log('  __nuxt exists:', nuxtState.exists);
    console.log('  __nuxt children:', nuxtState.childrenCount);
    console.log('  __nuxt innerHTML:', nuxtState.innerHTML.substring(0, 100) + '...');
    console.log('  __NUXT__ data exists:', nuxtState.hasNuxtData);
    console.log('  __NUXT__ keys:', nuxtState.nuxtDataKeys);

    // Check if Vue is loaded
    console.log('\n[Step 6] Vue Loading State:');
    const vueState = await page.evaluate(() => {
      return {
        hasVueGlobal: typeof Vue !== 'undefined',
        vueVersion: window.Vue?.version || 'not found',
        hasNuxtVue: window.nuxtApp?.vueApp !== undefined,
      };
    });
    console.log('  Vue global:', vueState.hasVueGlobal);
    console.log('  Vue version:', vueState.vueVersion);
    console.log('  Nuxt Vue app:', vueState.hasNuxtVue);

    // Get page structure
    console.log('\n[Step 7] Page Structure:');
    const pageStructure = await page.evaluate(() => {
      const body = document.body;
      return {
        bodyHTML: body.innerHTML.substring(0, 500),
        bodyChildren: body.children.length,
      };
    });
    console.log('  Body children:', pageStructure.bodyChildren);
    console.log('  Body HTML:', pageStructure.bodyHTML.substring(0, 200) + '...');

    // Check network requests for failed resources
    console.log('\n[Step 8] Checking for failed resources...');
    const failedResources = await page.evaluate(async () => {
      const performance = window.performance;
      const resources = performance.getEntriesByType('resource');
      return resources
        .filter(r => r.transferSize === 0 && r.decodedBodySize === 0)
        .map(r => ({
          name: r.name.substring(0, 80),
          type: r.initiatorType,
        }));
    });
    console.log('  Failed resources:', failedResources.length);
    failedResources.slice(0, 5).forEach(r => console.log(`    - ${r.name}`));

    console.log('\n' + '='.repeat(60));
    console.log('Debug Complete');
    console.log('='.repeat(60));

  } catch (error) {
    console.error('\n✗ Test error:', error.message);
    console.error(error.stack);
  } finally {
    await browser.close();
  }
}

testScheduleRuleModal().catch(console.error);
