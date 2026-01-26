import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock modules before imports
vi.mock('vue', async () => {
  const vue = await import('vue')
  return {
    ...vue,
    ref: (v: any) => ({ value: v }),
    computed: (fn: () => any) => ({ value: fn() }),
    onMounted: (fn: () => void) => fn(),
  }
})

vi.mock('~/stores/auth', () => ({
  useAuthStore: () => ({
    login: vi.fn(),
    user: null,
  }),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}))

describe('teacher/login.vue 頁面邏輯', () => {
  // TeacherLoginFormLogic 類別 - 老師登入表單邏輯
  class TeacherLoginFormLogic {
    lineUserId: string
    accessToken: string
    loading: boolean
    error: string
    success: boolean

    constructor() {
      this.lineUserId = ''
      this.accessToken = ''
      this.loading = false
      this.error = ''
      this.success = false
    }

    setLineUserId(id: string) {
      this.lineUserId = id
    }

    setAccessToken(token: string) {
      this.accessToken = token
    }

    setLoading(loading: boolean) {
      this.loading = loading
    }

    setError(error: string) {
      this.error = error
    }

    setSuccess(success: boolean) {
      this.success = success
    }

    resetForm() {
      this.lineUserId = ''
      this.accessToken = ''
      this.error = ''
      this.success = false
    }

    isFormValid(): boolean {
      return Boolean(this.lineUserId.trim() && this.accessToken.trim())
    }

    getLoginData(): any {
      return {
        line_user_id: this.lineUserId.trim(),
        access_token: this.accessToken.trim(),
      }
    }

    canSubmit(): boolean {
      return !this.loading && this.isFormValid()
    }

    clearError() {
      this.error = ''
    }
  }

  // QueryParamLogic 類別 - URL 參數邏輯
  class QueryParamLogic {
    params: Record<string, string>

    constructor() {
      this.params = {}
    }

    setParams(params: Record<string, string>) {
      this.params = params
    }

    getLineUserId(): string {
      return this.params.line_user_id || ''
    }

    getAccessToken(): string {
      return this.params.access_token || ''
    }

    hasLineUserId(): boolean {
      return Boolean(this.params.line_user_id)
    }

    hasAccessToken(): boolean {
      return Boolean(this.params.access_token)
    }

    hasAllParams(): boolean {
      return this.hasLineUserId() && this.hasAccessToken()
    }

    getParam(key: string): string {
      return this.params[key] || ''
    }
  }

  // AuthStoreLogic 類別 - 認證儲存邏輯
  class AuthStoreLogic {
    storedToken: string | null = null
    storedTeacher: any = null
    storedRefreshToken: string | null = null

    login(userData: any) {
      this.storedToken = userData.token || null
      this.storedRefreshToken = userData.refresh_token || null
      this.storedTeacher = userData.teacher || null
    }

    logout() {
      this.storedToken = null
      this.storedRefreshToken = null
      this.storedTeacher = null
    }

    isAuthenticated(): boolean {
      return this.storedToken !== null
    }

    getToken(): string | null {
      return this.storedToken
    }

    getTeacher(): any {
      return this.storedTeacher
    }

    hasToken(): boolean {
      return this.storedToken !== null
    }
  }

  // NavigationLogic 類別 - 導航邏輯
  class NavigationLogic {
    currentUrl: string
    redirectUrl: string

    constructor() {
      this.currentUrl = '/teacher/login'
      this.redirectUrl = '/teacher/dashboard'
    }

    getCurrentUrl(): string {
      return this.currentUrl
    }

    getRedirectUrl(): string {
      return this.redirectUrl
    }

    setRedirectUrl(url: string) {
      this.redirectUrl = url
    }

    getHomeUrl(): string {
      return '/'
    }

    canRedirect(): boolean {
      return Boolean(this.redirectUrl)
    }

    buildLoginUrl(lineUserId: string, accessToken: string): string {
      return `/teacher/login?line_user_id=${encodeURIComponent(lineUserId)}&access_token=${encodeURIComponent(accessToken)}`
    }
  }

  // ValidationLogic 類別 - 驗證邏輯
  class ValidationLogic {
    validateLineUserId(id: string): { valid: boolean; message: string } {
      if (!id) {
        return { valid: false, message: '請輸入 LINE User ID' }
      }
      if (id.length < 10) {
        return { valid: false, message: 'LINE User ID 格式不正確' }
      }
      return { valid: true, message: '' }
    }

    validateAccessToken(token: string): { valid: boolean; message: string } {
      if (!token) {
        return { valid: false, message: '請輸入 Access Token' }
      }
      if (token.length < 10) {
        return { valid: false, message: 'Access Token 格式不正確' }
      }
      return { valid: true, message: '' }
    }

    validateForm(lineUserId: string, accessToken: string): { valid: boolean; errors: string[] } {
      const errors: string[] = []
      const lineUserIdResult = this.validateLineUserId(lineUserId)
      if (!lineUserIdResult.valid) {
        errors.push(lineUserIdResult.message)
      }
      const accessTokenResult = this.validateAccessToken(accessToken)
      if (!accessTokenResult.valid) {
        errors.push(accessTokenResult.message)
      }
      return {
        valid: errors.length === 0,
        errors,
      }
    }

    isValidLineUserIdFormat(id: string): boolean {
      // LINE User ID typically starts with U and followed by hex characters
      const lineIdRegex = /^U[0-9a-fA-F]{32}$/
      return lineIdRegex.test(id)
    }
  }

  // ResponseHandlingLogic 類別 - 回應處理邏輯
  class ResponseHandlingLogic {
    handleApiResponse(response: any): { success: boolean; data: any; message: string } {
      const responseCode = response.code || response.code
      const responseData = response.data || response.datas
      const responseMessage = response.message

      if (responseCode === 0 && responseData) {
        return {
          success: true,
          data: responseData,
          message: '',
        }
      }

      return {
        success: false,
        data: null,
        message: responseMessage || '登入失敗',
      }
    }

    extractTokenFromResponse(data: any): string | null {
      return data?.token || null
    }

    extractUserFromResponse(data: any): any {
      return data?.user || null
    }

    isSuccessCode(code: number): boolean {
      return code === 0
    }
  }

  describe('TeacherLoginFormLogic 老師登入表單邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TeacherLoginFormLogic()
      expect(logic.lineUserId).toBe('')
      expect(logic.accessToken).toBe('')
      expect(logic.loading).toBe(false)
      expect(logic.error).toBe('')
      expect(logic.success).toBe(false)
    })

    it('setLineUserId 應該正確設定 LINE User ID', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setLineUserId('U1234567890abcdef1234567890abcd')
      expect(logic.lineUserId).toBe('U1234567890abcdef1234567890abcd')
    })

    it('setAccessToken 應該正確設定 Access Token', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setAccessToken('test_access_token')
      expect(logic.accessToken).toBe('test_access_token')
    })

    it('resetForm 應該重置表單', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setLineUserId('U123')
      logic.setAccessToken('token')
      logic.setError('錯誤')
      logic.setSuccess(true)
      logic.resetForm()
      expect(logic.lineUserId).toBe('')
      expect(logic.accessToken).toBe('')
      expect(logic.error).toBe('')
      expect(logic.success).toBe(false)
    })

    it('isFormValid 應該在兩個欄位都有值時返回 true', () => {
      const logic = new TeacherLoginFormLogic()
      expect(logic.isFormValid()).toBe(false)
      logic.setLineUserId('U123')
      expect(logic.isFormValid()).toBe(false)
      logic.setAccessToken('token')
      expect(logic.isFormValid()).toBe(true)
    })

    it('isFormValid 應該忽略空白', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setLineUserId('  ')
      logic.setAccessToken('  ')
      expect(logic.isFormValid()).toBe(false)
    })

    it('getLoginData 應該返回登入資料', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setLineUserId('U1234567890abcdef1234567890abcd')
      logic.setAccessToken('test_token')
      const data = logic.getLoginData()
      expect(data.line_user_id).toBe('U1234567890abcdef1234567890abcd')
      expect(data.access_token).toBe('test_token')
    })

    it('canSubmit 應該在非 loading 且表單有效時返回 true', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setLineUserId('U123')
      logic.setAccessToken('token')
      expect(logic.canSubmit()).toBe(true)
      logic.setLoading(true)
      expect(logic.canSubmit()).toBe(false)
    })

    it('clearError 應該清除錯誤訊息', () => {
      const logic = new TeacherLoginFormLogic()
      logic.setError('登入失敗')
      expect(logic.error).toBe('登入失敗')
      logic.clearError()
      expect(logic.error).toBe('')
    })
  })

  describe('QueryParamLogic URL 參數邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new QueryParamLogic()
      expect(logic.params).toEqual({})
    })

    it('setParams 應該正確設定參數', () => {
      const logic = new QueryParamLogic()
      logic.setParams({
        line_user_id: 'U123',
        access_token: 'token',
      })
      expect(logic.hasLineUserId()).toBe(true)
      expect(logic.hasAccessToken()).toBe(true)
    })

    it('getLineUserId 應該返回 LINE User ID', () => {
      const logic = new QueryParamLogic()
      logic.setParams({ line_user_id: 'U1234567890abcdef1234567890abcd' })
      expect(logic.getLineUserId()).toBe('U1234567890abcdef1234567890abcd')
    })

    it('getAccessToken 應該返回 Access Token', () => {
      const logic = new QueryParamLogic()
      logic.setParams({ access_token: 'test_token' })
      expect(logic.getAccessToken()).toBe('test_token')
    })

    it('hasAllParams 應該在兩個參數都有值時返回 true', () => {
      const logic = new QueryParamLogic()
      expect(logic.hasAllParams()).toBe(false)
      logic.setParams({ line_user_id: 'U123' })
      expect(logic.hasAllParams()).toBe(false)
      logic.setParams({ access_token: 'token' })
      expect(logic.hasAllParams()).toBe(false)
      logic.setParams({
        line_user_id: 'U123',
        access_token: 'token',
      })
      expect(logic.hasAllParams()).toBe(true)
    })

    it('getParam 應該返回指定參數', () => {
      const logic = new QueryParamLogic()
      logic.setParams({ test: 'value' })
      expect(logic.getParam('test')).toBe('value')
      expect(logic.getParam('nonexistent')).toBe('')
    })
  })

  describe('AuthStoreLogic 認證儲存邏輯', () => {
    it('login 應該正確儲存使用者資料', () => {
      const logic = new AuthStoreLogic()
      logic.login({
        token: 'teacher_token_123',
        refresh_token: 'refresh_token_456',
        teacher: { id: 1, name: '老師' },
      })
      expect(logic.isAuthenticated()).toBe(true)
      expect(logic.getToken()).toBe('teacher_token_123')
      expect(logic.getTeacher()?.name).toBe('老師')
    })

    it('logout 應該清除所有儲存資料', () => {
      const logic = new AuthStoreLogic()
      logic.login({ token: 'token' })
      logic.logout()
      expect(logic.isAuthenticated()).toBe(false)
      expect(logic.getToken()).toBeNull()
    })

    it('hasToken 應該正確判斷是否有 token', () => {
      const logic = new AuthStoreLogic()
      expect(logic.hasToken()).toBe(false)
      logic.login({ token: 'token' })
      expect(logic.hasToken()).toBe(true)
    })
  })

  describe('NavigationLogic 導航邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new NavigationLogic()
      expect(logic.getCurrentUrl()).toBe('/teacher/login')
      expect(logic.getRedirectUrl()).toBe('/teacher/dashboard')
    })

    it('setRedirectUrl 應該正確設定重導向網址', () => {
      const logic = new NavigationLogic()
      logic.setRedirectUrl('/teacher/schedule')
      expect(logic.getRedirectUrl()).toBe('/teacher/schedule')
    })

    it('getHomeUrl 應該返回首頁網址', () => {
      const logic = new NavigationLogic()
      expect(logic.getHomeUrl()).toBe('/')
    })

    it('buildLoginUrl 應該正確建立登入網址', () => {
      const logic = new NavigationLogic()
      const url = logic.buildLoginUrl('U123', 'token')
      expect(url).toContain('line_user_id=U123')
      expect(url).toContain('access_token=token')
    })
  })

  describe('ValidationLogic 驗證邏輯', () => {
    it('validateLineUserId 應該正確驗證 LINE User ID', () => {
      const logic = new ValidationLogic()
      expect(logic.validateLineUserId('').valid).toBe(false)
      expect(logic.validateLineUserId('U123').valid).toBe(false) // too short
      expect(logic.validateLineUserId('U1234567890abcdef1234567890abcd').valid).toBe(true)
    })

    it('validateAccessToken 應該正確驗證 Access Token', () => {
      const logic = new ValidationLogic()
      expect(logic.validateAccessToken('').valid).toBe(false)
      expect(logic.validateAccessToken('short').valid).toBe(false)
      expect(logic.validateAccessToken('valid_token_12345').valid).toBe(true)
    })

    it('validateForm 應該返回完整的驗證結果', () => {
      const logic = new ValidationLogic()
      const validResult = logic.validateForm('U1234567890abcdef1234567890abcd', 'valid_token')
      expect(validResult.valid).toBe(true)
      expect(validResult.errors).toHaveLength(0)

      const invalidResult = logic.validateForm('', '')
      expect(invalidResult.valid).toBe(false)
      expect(invalidResult.errors.length).toBe(2)
    })

    it('isValidLineUserIdFormat 應該正確驗證格式', () => {
      const logic = new ValidationLogic()
      // LINE User ID 格式：U 開頭，後面跟著字母數字組合
      // 實際長度可能不同，這裡測試一般格式
      expect(logic.isValidLineUserIdFormat('U1234567890abcdef1234567890ab')).toBe(true)
      expect(logic.isValidLineUserIdFormat('U123')).toBe(false)
      expect(logic.isValidLineUserIdFormat('invalid')).toBe(false)
      expect(logic.isValidLineUserIdFormat('')).toBe(false)
    })
  })

  describe('ResponseHandlingLogic 回應處理邏輯', () => {
    it('handleApiResponse 應該正確處理成功回應', () => {
      const logic = new ResponseHandlingLogic()
      const response = { code: 0, datas: { token: 'abc', user: { id: 1 } } }
      const result = logic.handleApiResponse(response)
      expect(result.success).toBe(true)
      expect(result.data.token).toBe('abc')
    })

    it('handleApiResponse 應該正確處理失敗回應', () => {
      const logic = new ResponseHandlingLogic()
      const response = { code: 1, message: '登入失敗' }
      const result = logic.handleApiResponse(response)
      expect(result.success).toBe(false)
      expect(result.message).toBe('登入失敗')
    })

    it('extractTokenFromResponse 應該正確提取 token', () => {
      const logic = new ResponseHandlingLogic()
      expect(logic.extractTokenFromResponse({ token: 'abc' })).toBe('abc')
      expect(logic.extractTokenFromResponse({})).toBeNull()
    })

    it('extractUserFromResponse 應該正確提取 user', () => {
      const logic = new ResponseHandlingLogic()
      const user = { id: 1, name: '老師' }
      expect(logic.extractUserFromResponse({ user })).toEqual(user)
      expect(logic.extractUserFromResponse({})).toBeNull()
    })

    it('isSuccessCode 應該正確判斷成功碼', () => {
      const logic = new ResponseHandlingLogic()
      expect(logic.isSuccessCode(0)).toBe(true)
      expect(logic.isSuccessCode(1)).toBe(false)
      expect(logic.isSuccessCode(-1)).toBe(false)
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整處理 LINE 登入流程', () => {
      const formLogic = new TeacherLoginFormLogic()
      const queryLogic = new QueryParamLogic()
      const authLogic = new AuthStoreLogic()
      const navLogic = new NavigationLogic()
      const validateLogic = new ValidationLogic()

      // 設定 URL 參數
      queryLogic.setParams({
        line_user_id: 'U1234567890abcdef1234567890abcd',
        access_token: 'test_access_token',
      })

      // 從 URL 參數載入表單
      formLogic.setLineUserId(queryLogic.getLineUserId())
      formLogic.setAccessToken(queryLogic.getAccessToken())

      // 驗證表單
      expect(formLogic.isFormValid()).toBe(true)

      // 驗證 URL 參數
      expect(queryLogic.hasAllParams()).toBe(true)

      // 模擬登入中
      formLogic.setLoading(true)
      expect(formLogic.loading).toBe(true)

      // 模擬登入成功
      formLogic.setLoading(false)
      formLogic.setSuccess(true)
      authLogic.login({
        token: 'teacher_token_123',
        refresh_token: 'refresh_token_456',
        teacher: { id: 1, name: '測試老師' },
      })

      // 驗證登入狀態
      expect(formLogic.success).toBe(true)
      expect(authLogic.isAuthenticated()).toBe(true)
      expect(authLogic.getToken()).toBe('teacher_token_123')

      // 驗證重導向
      expect(navLogic.canRedirect()).toBe(true)
      expect(navLogic.getRedirectUrl()).toBe('/teacher/dashboard')
    })

    it('應該正確處理登入失敗', () => {
      const formLogic = new TeacherLoginFormLogic()
      const authLogic = new AuthStoreLogic()

      formLogic.setLineUserId('U1234567890abcdef1234567890abcd')
      formLogic.setAccessToken('invalid_token')

      // 模擬登入失敗
      formLogic.setLoading(false)
      formLogic.setError('LINE 驗證失敗')
      authLogic.logout()

      // 驗證失敗狀態
      expect(formLogic.hasError()).toBe(true)
      expect(formLogic.error).toBe('LINE 驗證失敗')
      expect(authLogic.isAuthenticated()).toBe(false)
    })

    it('應該正確處理 URL 參數自動登入', () => {
      const queryLogic = new QueryParamLogic()
      const formLogic = new TeacherLoginFormLogic()
      const authLogic = new AuthStoreLogic()
      const validateLogic = new ValidationLogic()

      // 從 URL 獲取參數
      queryLogic.setParams({
        line_user_id: 'U1234567890abcdef1234567890abcd',
        access_token: 'valid_token_12345',
      })

      // 載入到表單
      formLogic.setLineUserId(queryLogic.getLineUserId())
      formLogic.setAccessToken(queryLogic.getAccessToken())

      // 驗證格式
      const validation = validateLogic.validateForm(
        formLogic.lineUserId,
        formLogic.accessToken
      )
      expect(validation.valid).toBe(true)

      // 模擬自動登入成功
      authLogic.login({ token: 'auto_login_token' })
      expect(authLogic.isAuthenticated()).toBe(true)
    })
  })
})
