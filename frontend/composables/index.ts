/**
 * Composables 統一匯出
 *
 * 警告：請勿在此檔案重複匯出已在個別檔案中匯出的函數
 * 這會導致 Nuxt 出現重複匯出警告
 */

// 錯誤處理 Composable（統一從 useErrorHandler 匯出）
export {
  useErrorHandler,
  getErrorMessage,
  getErrorAlertType,
  getErrorToastType,
  isSuccessCode,
  isPermissionError,
  isValidationError,
  isUnauthorizedError,
  handleApiResponse,
  handleError,
  handlePermissionError,
  handleValidationErrors,
  handleScheduleConflict,
  extractApiResponse,
  createAxiosErrorHandler,
  type ApiResponse,
  type ErrorHandlerOptions,
  type ErrorAction,
} from './useErrorHandler'

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

// 錯誤碼對照表
export {
  ERROR_MESSAGES,
  SUCCESS_CODES,
  SYSTEM_ERROR_CODES,
  DATABASE_ERROR_CODES,
  PERMISSION_ERROR_CODES,
  VALIDATION_ERROR_CODES,
  BUSINESS_ERROR_CODES,
  EXCEPTION_ERROR_CODES,
  TEACHER_ERROR_CODES,
  COURSE_ERROR_CODES,
  ROOM_ERROR_CODES,
  NOTIFICATION_ERROR_CODES,
  type ErrorCode,
  isSuccessCode as checkIsSuccessCode,
  isPermissionError as checkIsPermissionError,
  isValidationError as checkIsValidationError,
  isUnauthorizedError as checkIsUnauthorizedError,
} from '~/constants/errorCodes'
