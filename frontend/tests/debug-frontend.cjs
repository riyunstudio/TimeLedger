const { chromium } = require('playwright');

async function runTests() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();
  
  // Collect console errors
  const errors = [];
  page.on('console', msg => {
    if (msg.type() === 'error') {
      errors.push(msg.text());
    }
  });
  page.on('pageerror', err => {
    errors.push(err.message);
  });
  
  console.log('='.repeat(60));
  console.log('TimeLedger Frontend Debug Test');
  console.log('='.repeat(60));
  
  try {
    console.log('\n[Debug 1] Loading Admin Login Page');
    await page.goto('http://localhost:3005/admin/login', { timeout: 30000 });
    
    // Wait for any JavaScript to execute
    await page.waitForTimeout(5000);
    
    // Log all console messages
    console.log('\n[Debug 2] Console Errors:');
    if (errors.length > 0) {
      errors.forEach((err, i) => console.log(`  ${i + 1}. ${err}`));
    } else {
      console.log('  No console errors detected');
    }
    
    // Check page structure
    console.log('\n[Debug 3] Page Structure:');
    const html = await page.content();
    console.log('  HTML length:', html.length);
    console.log('  Has __nuxt div:', html.includes('__nuxt'));
    console.log('  Has teleports div:', html.includes('teleports'));
    console.log('  Has Nuxt config:', html.includes('__NUXT__.config'));
    
    // Check for Vue app
    console.log('\n[Debug 4] Vue App State:');
    const nuxtData = await page.evaluate(() => {
      return window.__NUXT__ ? 'Nuxt data exists' : 'No Nuxt data';
    });
    console.log('  ', nuxtData);
    
    // Check for mounted app
    const appMounted = await page.evaluate(() => {
      const nuxtDiv = document.getElementById('__nuxt');
      return nuxtDiv && nuxtDiv.children.length > 0 ? 'App mounted with children' : 'App empty';
    });
    console.log('  ', appMounted);
    
    // Get the inner HTML of __nuxt
    const nuxtInnerHtml = await page.evaluate(() => {
      const nuxtDiv = document.getElementById('__nuxt');
      return nuxtDiv ? nuxtDiv.innerHTML.substring(0, 500) : 'No __nuxt div';
    });
    console.log('  __nuxt inner HTML:', nuxtInnerHtml.substring(0, 100) + '...');
    
    // Try to evaluate a simple Vue test
    console.log('\n[Debug 5] Vue Instance Test:');
    const vueTest = await page.evaluate(() => {
      if (typeof Vue !== 'undefined') {
        return 'Vue global exists';
      }
      if (typeof window !== 'undefined' && Object.keys(window).some(k => k.includes('Vue'))) {
        return 'Vue in window';
      }
      return 'No Vue detected';
    });
    console.log('  ', vueTest);
    
    // Check network requests
    console.log('\n[Debug 6] Network Resources:');
    const resources = await page.evaluate(async () => {
      const resources = performance.getEntriesByType('resource');
      return resources.map(r => ({
        name: r.name.substring(0, 80),
        type: r.initiatorType,
        status: r.transferSize > 0 ? 'loaded' : 'pending'
      })).slice(0, 10);
    });
    resources.forEach(r => console.log(`  ${r.status} [${r.type}]: ${r.name}`));
    
    // Summary
    console.log('\n' + '='.repeat(60));
    console.log('Debug Summary:');
    console.log(`  Console Errors: ${errors.length}`);
    console.log(`  Page Rendered: ${html.length > 1000 ? 'Yes' : 'No (minimal HTML)'}`);
    console.log('='.repeat(60));
    
  } catch (error) {
    console.error('\n✗ Test error:', error.message);
  } finally {
    await browser.close();
  }
}

runTests().catch(console.error);
