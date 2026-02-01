import type { Meta, StoryObj } from '@storybook/vue3'
import Icon from './Icon.vue'

const meta: Meta<typeof Icon> = {
  title: 'Components/Icon',
  component: Icon,
  tags: ['autodocs'],
  argTypes: {
    icon: {
      control: 'select',
      options: [
        'chevron-down', 'chevron-right', 'chevron-left', 'close', 'calendar',
        'check', 'check-circle', 'people', 'user', 'settings', 'template',
        'resource', 'email', 'logout', 'arrow-right', 'plus', 'warning',
        'download', 'bell', 'info', 'error', 'success', 'time', 'search',
        'filter', 'edit', 'delete', 'eye', 'star', 'location', 'refresh',
        'upload', 'menu', 'more', 'link', 'external', 'home', 'office',
        'book', 'tag', 'trending', 'award', 'phone', 'copy', 'checklist',
        'rocket', 'heart', 'thumb', 'chat', 'send', 'image', 'document',
        'folder', 'lock', 'unlock', 'help'
      ],
      description: '圖示名稱',
    },
    size: {
      control: 'select',
      options: ['xs', 'sm', 'md', 'lg', 'xl', '2xl', '3xl', '4xl', '5xl'],
      description: '圖示尺寸',
    },
    color: {
      control: 'text',
      description: '自訂顏色 (CSS class 或顏色值)',
    },
    customSize: {
      control: 'text',
      description: '自訂尺寸 (CSS 值，如 24px)',
    },
  },
  args: {
    icon: 'check',
    size: 'md',
    color: 'currentColor',
  },
}

export default meta
type Story = StoryObj<typeof Icon>

// 所有圖示總覽
export const AllIcons: Story = {
  render: () => ({
    components: { Icon },
    setup() {
      const icons = [
        { name: 'chevron-down', label: '向下' },
        { name: 'chevron-right', label: '向右' },
        { name: 'chevron-left', label: '向左' },
        { name: 'close', label: '關閉' },
        { name: 'calendar', label: '日曆' },
        { name: 'check', label: '勾選' },
        { name: 'check-circle', label: '勾選圓圈' },
        { name: 'people', label: '多人' },
        { name: 'user', label: '使用者' },
        { name: 'settings', label: '設定' },
        { name: 'template', label: '模板' },
        { name: 'resource', label: '資源' },
        { name: 'email', label: 'Email' },
        { name: 'logout', label: '登出' },
        { name: 'arrow-right', label: '箭頭右' },
        { name: 'plus', label: '加號' },
        { name: 'warning', label: '警告' },
        { name: 'download', label: '下載' },
        { name: 'bell', label: '通知' },
        { name: 'info', label: '資訊' },
        { name: 'error', label: '錯誤' },
        { name: 'success', label: '成功' },
        { name: 'time', label: '時間' },
        { name: 'search', label: '搜尋' },
        { name: 'filter', label: '篩選' },
        { name: 'edit', label: '編輯' },
        { name: 'delete', label: '刪除' },
        { name: 'eye', label: '眼睛' },
        { name: 'star', label: '星號' },
        { name: 'location', label: '位置' },
        { name: 'refresh', label: '重新整理' },
        { name: 'upload', label: '上傳' },
        { name: 'menu', label: '選單' },
        { name: 'more', label: '更多' },
        { name: 'link', label: '連結' },
        { name: 'external', label: '外部' },
        { name: 'home', label: '首頁' },
        { name: 'office', label: '辦公室' },
        { name: 'book', label: '書籍' },
        { name: 'tag', label: '標籤' },
        { name: 'trending', label: '趨勢' },
        { name: 'award', label: '獎項' },
        { name: 'phone', label: '電話' },
        { name: 'copy', label: '複製' },
        { name: 'checklist', label: '檢查清單' },
        { name: 'rocket', label: '火箭' },
        { name: 'heart', label: '愛心' },
        { name: 'thumb', label: '拇指' },
        { name: 'chat', label: '對話' },
        { name: 'send', label: '傳送' },
        { name: 'image', label: '圖片' },
        { name: 'document', label: '文件' },
        { name: 'folder', label: '資料夾' },
        { name: 'lock', label: '鎖定' },
        { name: 'unlock', label: '解鎖' },
        { name: 'help', label: '幫助' },
      ]
      return { icons }
    },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen">
        <h2 class="text-white text-xl font-bold mb-6">所有圖示 (52 個)</h2>
        <div class="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 gap-4">
          <div
            v-for="icon in icons"
            :key="icon.name"
            class="flex flex-col items-center gap-2 p-4 bg-slate-800 rounded-lg hover:bg-slate-700 transition-colors"
          >
            <Icon :icon="icon.name" size="lg" class="text-white" />
            <span class="text-xs text-slate-400 text-center">{{ icon.label }}</span>
            <span class="text-[10px] text-slate-500 font-mono">{{ icon.name }}</span>
          </div>
        </div>
      </div>
    `,
  }),
}

// 尺寸展示
export const Sizes: Story = {
  render: () => ({
    components: { Icon },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen">
        <h2 class="text-white text-xl font-bold mb-6">尺寸展示</h2>
        <div class="flex items-end gap-8">
          <div v-for="size in ['xs', 'sm', 'md', 'lg', 'xl', '2xl', '3xl', '4xl', '5xl']" :key="size" class="flex flex-col items-center gap-2">
            <Icon icon="check" :size="size as any" class="text-indigo-400" />
            <span class="text-slate-400 text-sm font-mono">{{ size }}</span>
            <span class="text-slate-500 text-xs">{{ { xs: '12px', sm: '16px', md: '20px', lg: '24px', xl: '28px', '2xl': '32px', '3xl': '40px', '4xl': '48px', '5xl': '64px' }[size] }}</span>
          </div>
        </div>
      </div>
    `,
  }),
}

// 顏色展示
export const Colors: Story = {
  render: () => ({
    components: { Icon },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen">
        <h2 class="text-white text-xl font-bold mb-6">顏色展示</h2>
        <div class="flex flex-wrap gap-6">
          <div class="flex flex-col items-center gap-2">
            <Icon icon="check-circle" size="2xl" class="text-green-500" />
            <span class="text-slate-400 text-sm">text-green-500</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="error" size="2xl" class="text-red-500" />
            <span class="text-slate-400 text-sm">text-red-500</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="warning" size="2xl" class="text-yellow-500" />
            <span class="text-slate-400 text-sm">text-yellow-500</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="info" size="2xl" class="text-blue-500" />
            <span class="text-slate-400 text-sm">text-blue-500</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="bell" size="2xl" class="text-purple-500" />
            <span class="text-slate-400 text-sm">text-purple-500</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="2xl" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">text-yellow-400</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="heart" size="2xl" class="text-pink-500" />
            <span class="text-slate-400 text-sm">text-pink-500</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="lock" size="2xl" class="text-slate-400" />
            <span class="text-slate-400 text-sm">text-slate-400</span>
          </div>
        </div>
      </div>
    `,
  }),
}

// 尺寸與顏色組合矩陣
export const SizeColorMatrix: Story = {
  render: () => ({
    components: { Icon },
    setup() {
      const sizes = ['sm', 'md', 'lg', 'xl']
      const colors = [
        { class: 'text-indigo-400', name: 'indigo' },
        { class: 'text-emerald-400', name: 'emerald' },
        { class: 'text-amber-400', name: 'amber' },
        { class: 'text-rose-400', name: 'rose' },
      ]
      return { sizes, colors }
    },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen">
        <h2 class="text-white text-xl font-bold mb-6">尺寸與顏色組合矩陣</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr>
                <th class="text-slate-400 text-left pb-4 pr-4"></th>
                <th v-for="color in colors" :key="color.name" class="text-slate-400 text-center pb-4 px-4">
                  {{ color.name }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="size in sizes" :key="size">
                <td class="text-slate-400 py-2 pr-4 font-mono text-sm">{{ size }}</td>
                <td v-for="color in colors" :key="color.name" class="text-center py-2 px-4">
                  <div class="flex justify-center">
                    <Icon icon="check" :size="size as any" :class="color.class" />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    `,
  }),
}

// 圖示與文字搭配
export const WithText: Story = {
  render: () => ({
    components: { Icon },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen space-y-6">
        <h2 class="text-white text-xl font-bold mb-4">圖示與文字搭配</h2>

        <div class="space-y-4">
          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="user" size="lg" class="text-indigo-400" />
            <span class="text-white">個人資料</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="settings" size="lg" class="text-slate-400" />
            <span class="text-white">設定</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="bell" size="lg" class="text-yellow-500" />
            <span class="text-white">通知</span>
            <span class="ml-auto bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">3</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="email" size="lg" class="text-blue-400" />
            <span class="text-white">郵件</span>
            <span class="ml-auto bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">99+</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="check-circle" size="lg" class="text-green-500" />
            <span class="text-white">已完成</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="error" size="lg" class="text-red-500" />
            <span class="text-white">錯誤</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="warning" size="lg" class="text-yellow-500" />
            <span class="text-white">警告</span>
          </div>

          <div class="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
            <Icon icon="help" size="lg" class="text-purple-400" />
            <span class="text-white">說明</span>
          </div>
        </div>
      </div>
    `,
  }),
}

// 按鈕中的圖示
export const InButtons: Story = {
  render: () => ({
    components: { Icon },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen space-y-6">
        <h2 class="text-white text-xl font-bold mb-4">按鈕中的圖示</h2>

        <div class="flex flex-wrap gap-4">
          <button class="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
            <Icon icon="plus" size="sm" />
            <span>新增</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors">
            <Icon icon="check" size="sm" />
            <span>確認</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 bg-amber-600 text-white rounded-lg hover:bg-amber-700 transition-colors">
            <Icon icon="edit" size="sm" />
            <span>編輯</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors">
            <Icon icon="delete" size="sm" />
            <span>刪除</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 bg-slate-700 text-white rounded-lg hover:bg-slate-600 transition-colors">
            <Icon icon="download" size="sm" />
            <span>下載</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
            <Icon icon="upload" size="sm" />
            <span>上傳</span>
          </button>
        </div>

        <div class="flex flex-wrap gap-4">
          <button class="flex items-center justify-center w-10 h-10 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
            <Icon icon="chevron-left" size="lg" />
          </button>

          <button class="flex items-center justify-center w-10 h-10 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
            <Icon icon="chevron-right" size="lg" />
          </button>

          <button class="flex items-center justify-center w-10 h-10 bg-slate-700 text-white rounded-lg hover:bg-slate-600 transition-colors">
            <Icon icon="refresh" size="lg" />
          </button>

          <button class="flex items-center justify-center w-10 h-10 bg-slate-700 text-white rounded-lg hover:bg-slate-600 transition-colors">
            <Icon icon="search" size="lg" />
          </button>

          <button class="flex items-center justify-center w-10 h-10 bg-slate-700 text-white rounded-lg hover:bg-slate-600 transition-colors">
            <Icon icon="close" size="lg" />
          </button>
        </div>

        <div class="flex flex-wrap gap-4">
          <button class="flex items-center gap-2 px-4 py-2 border border-slate-600 text-slate-300 rounded-lg hover:bg-slate-800 transition-colors">
            <Icon icon="arrow-right" size="sm" />
            <span>前往</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 border border-slate-600 text-slate-300 rounded-lg hover:bg-slate-800 transition-colors">
            <Icon icon="external" size="sm" />
            <span>開啟連結</span>
          </button>

          <button class="flex items-center gap-2 px-4 py-2 border border-slate-600 text-slate-300 rounded-lg hover:bg-slate-800 transition-colors">
            <Icon icon="link" size="sm" />
            <span>複製連結</span>
          </button>
        </div>
      </div>
    `,
  }),
}

// 狀態圖示
export const StatusIcons: Story = {
  render: () => ({
    components: { Icon },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen space-y-6">
        <h2 class="text-white text-xl font-bold mb-4">狀態圖示</h2>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div class="p-4 bg-slate-800 rounded-lg flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-green-500/20 flex items-center justify-center">
              <Icon icon="check-circle" size="2xl" class="text-green-500" />
            </div>
            <div>
              <h3 class="text-white font-semibold">成功</h3>
              <p class="text-slate-400 text-sm">操作已完成</p>
            </div>
          </div>

          <div class="p-4 bg-slate-800 rounded-lg flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-red-500/20 flex items-center justify-center">
              <Icon icon="error" size="2xl" class="text-red-500" />
            </div>
            <div>
              <h3 class="text-white font-semibold">錯誤</h3>
              <p class="text-slate-400 text-sm">發生錯誤</p>
            </div>
          </div>

          <div class="p-4 bg-slate-800 rounded-lg flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-yellow-500/20 flex items-center justify-center">
              <Icon icon="warning" size="2xl" class="text-yellow-500" />
            </div>
            <div>
              <h3 class="text-white font-semibold">警告</h3>
              <p class="text-slate-400 text-sm">請注意</p>
            </div>
          </div>

          <div class="p-4 bg-slate-800 rounded-lg flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-blue-500/20 flex items-center justify-center">
              <Icon icon="info" size="2xl" class="text-blue-500" />
            </div>
            <div>
              <h3 class="text-white font-semibold">資訊</h3>
              <p class="text-slate-400 text-sm">提示訊息</p>
            </div>
          </div>
        </div>

        <div class="p-4 bg-slate-800 rounded-lg">
          <h3 class="text-white font-semibold mb-4">載入狀態</h3>
          <div class="flex items-center gap-4">
            <Icon icon="refresh" size="lg" class="text-indigo-400 animate-spin" />
            <span class="text-slate-300">載入中...</span>
          </div>
        </div>
      </div>
    `,
  }),
}

// 自訂尺寸
export const CustomSize: Story = {
  render: () => ({
    components: { Icon },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen">
        <h2 class="text-white text-xl font-bold mb-6">自訂尺寸</h2>
        <div class="flex items-end gap-8">
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="custom" customSize="16px" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">16px</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="custom" customSize="24px" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">24px</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="custom" customSize="32px" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">32px</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="custom" customSize="48px" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">48px</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="custom" customSize="64px" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">64px</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <Icon icon="star" size="custom" customSize="96px" class="text-yellow-400" />
            <span class="text-slate-400 text-sm">96px</span>
          </div>
        </div>
      </div>
    `,
  }),
}

// 互動式 Playground
export const Playground: Story = {
  render: (args) => ({
    components: { Icon },
    setup() {
      return { args }
    },
    template: `
      <div class="p-6 bg-slate-900 min-h-screen">
        <h2 class="text-white text-xl font-bold mb-6">Playground</h2>
        <div class="flex justify-center py-12">
          <Icon v-bind="args" />
        </div>
      </div>
    `,
  }),
}
