/**
 * Icon Library Composable
 *
 * 提供便捷的圖標使用方式，支援:
 * - 直接使用組件
 * - 動態圖標名稱
 * - 常用圖標集合
 */

import { computed } from 'vue'

/**
 * 可用的圖標名稱類型
 */
export type IconName =
  | 'chevron-down'
  | 'chevron-right'
  | 'chevron-left'
  | 'close'
  | 'calendar'
  | 'check'
  | 'check-circle'
  | 'people'
  | 'user'
  | 'settings'
  | 'template'
  | 'resource'
  | 'email'
  | 'logout'
  | 'arrow-right'
  | 'plus'
  | 'warning'
  | 'download'
  | 'bell'
  | 'info'
  | 'error'
  | 'success'
  | 'time'
  | 'search'
  | 'filter'
  | 'edit'
  | 'delete'
  | 'eye'
  | 'star'
  | 'location'
  | 'refresh'
  | 'upload'
  | 'menu'
  | 'more'
  | 'link'
  | 'external'
  | 'home'
  | 'office'
  | 'book'
  | 'tag'
  | 'trending'
  | 'award'
  | 'phone'
  | 'copy'
  | 'checklist'
  | 'rocket'
  | 'heart'
  | 'thumb'
  | 'chat'
  | 'send'
  | 'image'
  | 'document'
  | 'folder'
  | 'lock'
  | 'unlock'
  | 'help'

/**
 * 圖標大小類型
 */
export type IconSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | '4xl' | '5xl' | 'custom'

/**
 * 常用圖標集合 - 依功能分類
 */
export const IconSets = {
  /** 導航相關 */
  navigation: {
    home: 'home',
    menu: 'menu',
    chevronDown: 'chevron-down',
    chevronRight: 'chevron-right',
    chevronLeft: 'chevron-left',
    arrowRight: 'arrow-right',
  },

  /** 操作相關 */
  actions: {
    close: 'close',
    check: 'check',
    plus: 'plus',
    edit: 'edit',
    delete: 'delete',
    copy: 'copy',
    refresh: 'refresh',
    search: 'search',
    filter: 'filter',
    upload: 'upload',
    download: 'download',
    link: 'link',
    external: 'external',
  },

  /** 狀態相關 */
  status: {
    success: 'success',
    error: 'error',
    warning: 'warning',
    info: 'info',
    checkCircle: 'check-circle',
  },

  /** 用戶相關 */
  user: {
    user: 'user',
    people: 'people',
    logout: 'logout',
    lock: 'lock',
    unlock: 'unlock',
  },

  /** 內容相關 */
  content: {
    calendar: 'calendar',
    time: 'time',
    location: 'location',
    phone: 'phone',
    email: 'email',
    eye: 'eye',
    star: 'star',
  },

  /** 檔案相關 */
  file: {
    image: 'image',
    document: 'document',
    folder: 'folder',
    template: 'template',
    resource: 'resource',
    tag: 'tag',
  },

  /** 互動相關 */
  interaction: {
    bell: 'bell',
    chat: 'chat',
    send: 'send',
    help: 'help',
    more: 'more',
    thumb: 'thumb',
    heart: 'heart',
  },

  /** 商務相關 */
  business: {
    office: 'office',
    book: 'book',
    trending: 'trending',
    award: 'award',
    checklist: 'checklist',
    rocket: 'rocket',
  },
}

/**
 * 圖標名稱對應的中文說明
 */
export const IconDescriptions: Record<IconName, string> = {
  'chevron-down': '展開',
  'chevron-right': '展開',
  'chevron-left': '收起',
  'close': '關閉',
  'calendar': '日曆',
  'check': '勾選',
  'check-circle': '完成',
  'people': '人才',
  'user': '用戶',
  'settings': '設定',
  'template': '範本',
  'resource': '資源',
  'email': '郵件',
  'logout': '登出',
  'arrow-right': '箭頭',
  'plus': '新增',
  'warning': '警告',
  'download': '下載',
  'bell': '通知',
  'info': '資訊',
  'error': '錯誤',
  'success': '成功',
  'time': '時間',
  'search': '搜尋',
  'filter': '篩選',
  'edit': '編輯',
  'delete': '刪除',
  'eye': '查看',
  'star': '星號',
  'location': '位置',
  'refresh': '重新整理',
  'upload': '上傳',
  'menu': '選單',
  'more': '更多',
  'link': '連結',
  'external': '外部連結',
  'home': '首頁',
  'office': '辦公室',
  'book': '課程',
  'tag': '標籤',
  'trending': '趨勢',
  'award': '獎項',
  'phone': '電話',
  'copy': '複製',
  'checklist': '清單',
  'rocket': '快速',
  'heart': '喜歡',
  'thumb': '讚',
  'chat': '對話',
  'send': '傳送',
  'image': '圖片',
  'document': '文件',
  'folder': '資料夾',
  'lock': '鎖定',
  'unlock': '解鎖',
  'help': '說明',
}

/**
 * 使用圖標庫的 composable
 */
export function useIcon() {
  /**
   * 驗證圖標是否存在
   */
  const isValidIcon = (name: string): name is IconName => {
    return name in IconDescriptions
  }

  /**
   * 取得圖標的中文說明
   */
  const getIconDescription = (name: IconName): string => {
    return IconDescriptions[name] || name
  }

  /**
   * 取得所有可用圖標列表
   */
  const getAllIcons = (): IconName[] => {
    return Object.keys(IconDescriptions) as IconName[]
  }

  /**
   * 依類別取得圖標
   */
  const getIconsByCategory = (category: keyof typeof IconSets): IconName[] => {
    return Object.values(IconSets[category]) as IconName[]
  }

  /**
   * 圖標大小的 Tailwind 類別對應
   */
  const sizeClasses: Record<IconSize, string> = {
    xs: 'w-3 h-3',
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
    xl: 'w-7 h-7',
    '2xl': 'w-8 h-8',
    '3xl': 'w-10 h-10',
    '4xl': 'w-12 h-12',
    '5xl': 'w-16 h-16',
    custom: '',
  }

  /**
   * 取得圖標大小的類別
   */
  const getSizeClass = (size: IconSize, customSize?: string): string => {
    if (size === 'custom' && customSize) {
      return customSize
    }
    return sizeClasses[size] || sizeClasses.md
  }

  return {
    IconSets,
    IconDescriptions,
    isValidIcon,
    getIconDescription,
    getAllIcons,
    getIconsByCategory,
    getSizeClass,
  }
}

/**
 * 便捷導出 - 直接使用圖標組件
 */
export function useIcons() {
  const { IconSets, getIconDescription, isValidIcon } = useIcon()

  return {
    IconSets,
    getIconDescription,
    isValidIcon,
  }
}
