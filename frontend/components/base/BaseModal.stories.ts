import type { Meta, StoryObj } from '@storybook/vue3'
import BaseModal from './BaseModal.vue'

const meta: Meta<typeof BaseModal> = {
  title: 'Components/BaseModal',
  component: BaseModal,
  tags: ['autodocs'],
  argTypes: {
    modelValue: { control: 'boolean' },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg', 'xl'],
      description: 'Modal 大小',
    },
    mobilePosition: {
      control: 'select',
      options: ['center', 'bottom'],
      description: '行動裝置位置',
    },
    closeOnBackdrop: { control: 'boolean' },
    showCloseButton: { control: 'boolean' },
    close: { action: 'close' },
    'update:modelValue': { action: 'update:modelValue' },
  },
  args: {
    modelValue: true,
    title: '範例標題',
    size: 'md',
    closeOnBackdrop: true,
    showCloseButton: true,
    mobilePosition: 'bottom',
  },
}

export default meta
type Story = StoryObj<typeof BaseModal>

// 預設故事
export const Default: Story = {
  render: (args) => ({
    components: { BaseModal },
    setup() {
      return { args }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <BaseModal v-bind="args">
          <p class="text-slate-300">
            這是 Modal 的內容區域。您可以在這裡放置任何內容，包括表單、文字、或其他元件。
          </p>
          <p class="text-slate-300 mt-2">
            Modal 支援玻璃擬態效果和流暢的動畫過渡。
          </p>
        </BaseModal>
      </div>
    `,
  }),
}

// 不同尺寸展示
export const Sizes: Story = {
  render: () => ({
    components: { BaseModal },
    data() {
      return {
        activeModal: null as string | null,
      }
    },
    methods: {
      openModal(size: string) {
        this.activeModal = size
      },
      closeModal() {
        this.activeModal = null
      },
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <h2 class="text-white text-xl font-bold mb-4">Modal 尺寸展示</h2>
        <div class="flex gap-4 flex-wrap">
          <button
            v-for="size in ['sm', 'md', 'lg', 'xl']"
            :key="size"
            @click="openModal(size)"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            {{ size.toUpperCase() }} ({{ { sm: '320px', md: '384px', lg: '448px', xl: '768px' }[size] }})
          </button>
        </div>

        <BaseModal
          v-if="activeModal"
          v-model="activeModal"
          :title="\`\${activeModal.toUpperCase()} Modal\`"
          :size="activeModal as any"
          @close="closeModal"
        >
          <p class="text-slate-300">
            這是 {{ activeModal }} 尺寸的 Modal 內容。
          </p>
          <p class="text-slate-300 mt-2">
            寬度：{{ { sm: '320px', md: '384px', lg: '448px', xl: '768px' }[activeModal] }}
          </p>
        </BaseModal>
      </div>
    `,
  }),
}

// 無標題 Modal
export const WithoutTitle: Story = {
  render: (args) => ({
    components: { BaseModal },
    setup() {
      return { args }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <BaseModal v-bind="args">
          <p class="text-slate-300">
            這是一個沒有標題的 Modal，只有關閉按鈕和內容。
          </p>
        </BaseModal>
      </div>
    `,
  }),
  args: {
    modelValue: true,
    title: undefined,
    showCloseButton: true,
  },
}

// 自訂 Header Slot
export const CustomHeader: Story = {
  render: () => ({
    components: { BaseModal },
    data() {
      return { showModal: true }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <button
          @click="showModal = true"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg"
        >
          打開自訂 Header Modal
        </button>

        <BaseModal v-model="showModal" size="lg">
          <template #header>
            <div class="flex items-center justify-between w-full">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-indigo-500 rounded-full flex items-center justify-center">
                  <span class="text-white text-lg">🎉</span>
                </div>
                <div>
                  <h3 class="text-white font-semibold">恭喜完成！</h3>
                  <p class="text-slate-400 text-sm">您已成功完成設定</p>
                </div>
              </div>
            </div>
          </template>
          <p class="text-slate-300">
            使用 slot 可以自訂 Modal 的任何部分，包括 header、body 和 footer。
          </p>
          <p class="text-slate-300 mt-2">
            這裡展示了自訂 Header 的範例，包含圖示和副標題。
          </p>
        </BaseModal>
      </div>
    `,
  }),
}

// 自訂 Footer
export const CustomFooter: Story = {
  render: () => ({
    components: { BaseModal },
    data() {
      return { showModal: true }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <button
          @click="showModal = true"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg"
        >
          打開含自訂 Footer 的 Modal
        </button>

        <BaseModal v-model="showModal" title="確認刪除">
          <template #default>
            <p class="text-slate-300">
              您確定要刪除此項目嗎？此操作無法復原。
            </p>
          </template>
          <template #footer>
            <div class="flex justify-end gap-3 w-full">
              <button
                @click="showModal = false"
                class="px-4 py-2 text-slate-300 hover:text-white transition-colors"
              >
                取消
              </button>
              <button
                @click="showModal = false"
                class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
              >
                確認刪除
              </button>
            </div>
          </template>
        </BaseModal>
      </div>
    `,
  }),
}

// 關閉背景互動
export const DisableBackdropClose: Story = {
  render: (args) => ({
    components: { BaseModal },
    setup() {
      return { args }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <BaseModal v-bind="args">
          <p class="text-slate-300">
            點擊背景不會關閉此 Modal。您必須點擊關閉按鈕。
          </p>
        </BaseModal>
      </div>
    `,
  }),
  args: {
    modelValue: true,
    title: '禁止背景關閉',
    closeOnBackdrop: false,
  },
}

// 行動裝置位置
export const MobilePositions: Story = {
  render: () => ({
    components: { BaseModal },
    data() {
      return {
        showBottom: true,
        showCenter: false,
      }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <h2 class="text-white text-xl font-bold mb-4">行動裝置位置</h2>
        <div class="flex gap-4 flex-wrap">
          <button
            @click="showBottom = true"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg"
          >
            底部彈出 (Bottom)
          </button>
          <button
            @click="showCenter = true"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg"
          >
            置中顯示 (Center)
          </button>
        </div>

        <BaseModal
          v-model="showBottom"
          title="底部彈出"
          mobilePosition="bottom"
        >
          <p class="text-slate-300">
            在行動裝置上，Modal 會從底部彈出，提供更好的使用者體驗。
          </p>
        </BaseModal>

        <BaseModal
          v-model="showCenter"
          title="置中顯示"
          mobilePosition="center"
        >
          <p class="text-slate-300">
            置中顯示的 Modal，適用於需要完整視覺焦點的場景。
          </p>
        </BaseModal>
      </div>
    `,
  }),
}

// 複雜內容展示
export const ComplexContent: Story = {
  render: () => ({
    components: { BaseModal },
    data() {
      return { showModal: true }
    },
    template: `
      <div class="p-4 bg-slate-900 min-h-[400px]">
        <button
          @click="showModal = true"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg"
        >
          打開複雜內容 Modal
        </button>

        <BaseModal v-model="showModal" title="表單範例" size="lg">
          <form class="space-y-4">
            <div>
              <label class="block text-slate-300 text-sm mb-1">姓名</label>
              <input
                type="text"
                class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white focus:outline-none focus:border-indigo-500"
                placeholder="請輸入姓名"
              />
            </div>
            <div>
              <label class="block text-slate-300 text-sm mb-1">Email</label>
              <input
                type="email"
                class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white focus:outline-none focus:border-indigo-500"
                placeholder="請輸入 Email"
              />
            </div>
            <div>
              <label class="block text-slate-300 text-sm mb-1">訊息</label>
              <textarea
                rows="3"
                class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white focus:outline-none focus:border-indigo-500"
                placeholder="請輸入訊息"
              ></textarea>
            </div>
          </form>
          <template #footer>
            <div class="flex justify-end gap-3 w-full">
              <button
                @click="showModal = false"
                class="px-4 py-2 text-slate-300 hover:text-white transition-colors"
              >
                取消
              </button>
              <button
                @click="showModal = false"
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                送出
              </button>
            </div>
          </template>
        </BaseModal>
      </div>
    `,
  }),
}
