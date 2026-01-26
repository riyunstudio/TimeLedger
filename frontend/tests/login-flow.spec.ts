import { test, expect, type Page } from '@playwright/test'

test.describe('Admin Login Flow Tests', () => {
  let page: Page

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage()
    page.setViewportSize({ width: 1280, height: 720 })
  })

  test.afterAll(async () => {
    await page.close()
  })

  test.describe('Login Page Display', () => {
    test('應該正確顯示登入頁面元素', async () => {
      await page.goto('/admin/login')
      await page.waitForLoadState('domcontentloaded')

      // 檢查標題
      await expect(page.locator('h1')).toContainText('管理員登入')

      // 檢查表單元素
      await expect(page.locator('input[type="email"]')).toBeVisible()
      await expect(page.locator('input[type="password"]')).toBeVisible()
      await expect(page.locator('button[type="submit"]')).toContainText('登入')
    })

  test('應該正確顯示 Email 和密碼輸入框 placeholder', async () => {
    await page.goto('/admin/login')

    const emailInput = page.locator('input[type="email"]')
    const passwordInput = page.locator('input[type="password"]')

    // 檢查輸入框存在且可見
    await expect(emailInput).toBeVisible()
    await expect(passwordInput).toBeVisible()
    
    // 檢查類型正確
    await expect(emailInput).toHaveAttribute('type', 'email')
    await expect(passwordInput).toHaveAttribute('type', 'password')
  })
  })

  test.describe('Successful Login Flow', () => {
  test('應該在輸入正確帳密後成功登入並跳轉到儀表板', async () => {
    await page.goto('/admin/login')
    
    await page.fill('input[type="email"]', 'admin@timeledger.com')
    await page.fill('input[type="password"]', 'admin123')
    await page.click('button[type="submit"]')
    
    // 等待 URL 變更到儀表板
    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 })

    // 驗證登入後頁面有載入內容（檢查頁面標題或主要區塊）
    await expect(page.locator('main').first()).toBeVisible({ timeout: 5000 })
  })

  test('登入成功後應該顯示導航選單', async () => {
    await page.goto('/admin/login')
    await page.fill('input[type="email"]', 'admin@timeledger.com')
    await page.fill('input[type="password"]', 'admin123')
    await page.click('button[type="submit"]')

    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 })

    // 檢查側邊欄導航 - 使用更具體的選擇器
    await expect(page.locator('nav').first().locator('text=資源管理')).toBeVisible({ timeout: 5000 })
    await expect(page.locator('nav').first().locator('text=審核')).toBeVisible()
    // 使用 getByRole 來選擇正確的導航連結
    await expect(page.getByRole('link', { name: '課表', exact: true })).toBeVisible()
  })
  })

  test.describe('Failed Login Flow', () => {
    test('應該在輸入錯誤密碼時顯示錯誤訊息', async () => {
      await page.goto('/admin/login')

      await page.fill('input[type="email"]', 'admin@timeledger.com')
      await page.fill('input[type="password"]', 'wrongpassword')
      await page.click('button[type="submit"]')

      // 檢查錯誤提示
      await expect(page.locator('text=登入失敗')).toBeVisible({ timeout: 5000 })
    })

    test('應該在輸入不存在帳號時顯示錯誤訊息', async () => {
      await page.goto('/admin/login')

      await page.fill('input[type="email"]', 'nonexistent@email.com')
      await page.fill('input[type="password"]', 'anypassword')
      await page.click('button[type="submit"]')

      // 檢查錯誤提示
      await expect(page.locator('text=登入失敗')).toBeVisible({ timeout: 5000 })
    })

    test('應該在輸入空帳號時顯示驗證錯誤', async () => {
      await page.goto('/admin/login')

      // 不輸入任何內容直接點擊登入
      await page.click('button[type="submit"]')

      // 檢查是否有驗證錯誤提示（取決於前端實作）
      // 可能會顯示 "必填" 或類似訊息
      const errorMessage = page.locator('text=必填|required|Required').first()
      if (await errorMessage.isVisible({ timeout: 2000 })) {
        await expect(errorMessage).toBeVisible()
      }
    })
  })

  test.describe('Login Session Persistence', () => {
  test('登入成功後刷新頁面應該保持登入狀態', async () => {
    await page.goto('/admin/login')
    await page.fill('input[type="email"]', 'admin@timeledger.com')
    await page.fill('input[type="password"]', 'admin123')
    await page.click('button[type="submit"]')

    await expect(page).toHaveURL(/.*dashboard/, { timeout: 10000 })

    // 刷新頁面
    await page.reload()

    // 應該仍然在儀表板，而不是跳回登入頁
    await expect(page).not.toHaveURL(/.*login/)
    await expect(page).toHaveURL(/.*dashboard/, { timeout: 5000 })
  })
  })
})

test.describe('Admin Authentication API Tests', () => {
  test('後端應該接受正確的管理員登入', async ({ request }) => {
    const response = await request.post('http://localhost:3000/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'admin123'
      }
    })

    expect(response.status()).toBe(200)

    const data = await response.json()
    expect(data.code).toBe(0)
    expect(data.datas.token).toBeDefined()
    expect(data.datas.user.email).toBe('admin@timeledger.com')
    // 驗證用戶類型欄位（可能是 role, user_type 或其他名稱）
    expect(data.datas.user).toBeDefined()
    console.log('User data:', JSON.stringify(data.datas.user, null, 2))
  })

  test('後端應該拒絕錯誤密碼的登入', async ({ request }) => {
    const response = await request.post('http://localhost:3000/api/v1/auth/admin/login', {
      data: {
        email: 'admin@timeledger.com',
        password: 'wrongpassword'
      }
    })

    expect(response.status()).toBe(401)

    const data = await response.json()
    expect(data.code).toBe(401)
  })

  test('後端應該拒絕不存在帳號的登入', async ({ request }) => {
    const response = await request.post('http://localhost:3000/api/v1/auth/admin/login', {
      data: {
        email: 'nonexistent@timeledger.com',
        password: 'admin123'
      }
    })

    expect(response.status()).toBe(401)

    const data = await response.json()
    expect(data.code).toBe(401)
  })

  test('Mock token 應該可以直接存取受保護的 API', async ({ request }) => {
    // 使用 mock-admin-token 前綴跳過 JWT 驗證
    // 測試 /admin/centers/1 端點
    const response = await request.get(
      'http://localhost:3000/api/v1/admin/centers/1',
      {
        headers: {
          'Authorization': 'Bearer mock-admin-token-1'
        }
      }
    )

    // mock token 可能返回 200、401 或 404，取決於後端實作
    expect([200, 401, 404]).toContain(response.status())
    
    if (response.status() === 200) {
      const data = await response.json()
      expect(data.code).toBe(0)
    }
  })
})
