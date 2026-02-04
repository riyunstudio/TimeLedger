/**
 * Composables 統一匯出
 *
 * 警告：請勿在此檔案重複匯出已在個別檔案中匯出的函數
 * 這會導致 Nuxt 出現重複匯出警告
 *
 * Nuxt 會自動掃描 composables 目錄，
 * 此檔案僅作為統一路由，不建議在此重新匯出已存在的項目
 */

// 快速錯誤處理 Composable
export {
  useQuickError,
  useFormErrors,
  useAsyncOperation,
  useErrorBoundary,
  createApiRequest,
} from './useQuickError'

// Alert 函數
export {
  alertError,
  alertWarning,
  alertSuccess,
  alertInfo,
  alertConfirm,
  showAlert,
  type AlertOptions,
  type AlertButton,
  type AlertType,
} from './useAlert'

// Toast 函數
export {
  useToast,
  type ToastType,
} from './useToast'
