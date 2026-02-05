<template>
  <div class="flex justify-center py-4">
    <!-- 手機框容器 -->
    <div class="line-phone-frame">
      <!-- 手機頂部狀態列 -->
      <div class="phone-status-bar">
        <span class="time">{{ currentTime }}</span>
        <div class="status-icons">
          <svg class="signal-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M2 22h2V12H2v10zm4 0h2V9H6v13zm4 0h2V6h-2v16zm4 0h2V3h-2v19zm4 0h2V0h-2v22z"/>
          </svg>
          <svg class="wifi-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M1 9l2 2c4.97-4.97 13.03-4.97 18 0l2-2C16.93 2.93 7.08 2.93 1 9zm8 8l3 3 3-3c-1.65-1.66-4.34-1.66-6 0zm-4-4l2 2c2.76-2.76 7.24-2.76 10 0l2-2C15.14 9.14 8.87 9.14 5 13z"/>
          </svg>
          <svg class="battery-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M15.67 4H14V2h-4v2H8.33C7.6 4 7 4.6 7 5.33v15.33C7 21.4 7.6 22 8.33 22h7.33c.74 0 1.34-.6 1.34-1.33V5.33C17 4.6 16.4 4 15.67 4z"/>
          </svg>
        </div>
      </div>

      <!-- LINE 標題列 -->
      <div class="line-header">
        <div class="line-title">
          <svg class="line-logo" viewBox="0 0 24 24" fill="currentColor">
            <circle cx="12" cy="12" r="10" fill="#06C755"/>
            <path d="M17 12c0-2.76-2.24-5-5-5s-5 2.24-5 5c0 1.66.81 3.13 2.04 4.02L8 17.41l.96-.41c.69.24 1.45.37 2.24.37.79 0 1.55-.13 2.24-.37l.96.41-1.04-1.39C16.19 15.13 17 13.66 17 12z" fill="white"/>
          </svg>
          <span class="official-account">TimeLedger</span>
        </div>
        <div class="line-menu">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <circle cx="12" cy="5" r="2"/>
            <circle cx="12" cy="12" r="2"/>
            <circle cx="12" cy="19" r="2"/>
          </svg>
        </div>
      </div>

      <!-- 氣泡訊息容器 -->
      <div class="message-container">
        <!-- 日期標籤 -->
        <div v-if="showDate" class="date-label">
          {{ formattedDate }}
        </div>

        <!-- Flex Message 氣泡 -->
        <div class="flex-message-bubble">
          <!-- 氣泡箭頭 -->
          <div class="bubble-arrow"></div>

          <!-- 氣泡內容 -->
          <div class="bubble-content">
            <!-- 標題區塊 -->
            <div class="bubble-header">
              <div class="notification-icon">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/>
                </svg>
              </div>
              <h3 class="bubble-title">{{ title }}</h3>
            </div>

            <!-- 分隔線 -->
            <div class="bubble-divider"></div>

            <!-- 內容區塊 -->
            <div class="bubble-body">
              <slot>
                <div v-if="content" class="bubble-text">
                  {{ content }}
                </div>
                <div v-else class="bubble-fields">
                  <div v-for="(field, index) in fields" :key="index" class="bubble-field">
                    <span class="field-icon">{{ field.icon }}</span>
                    <span class="field-label">{{ field.label }}</span>
                    <span class="field-value">{{ field.value }}</span>
                  </div>
                </div>
              </slot>
            </div>

            <!-- 警告區塊（可選） -->
            <div v-if="warning" class="bubble-warning">
              <div class="warning-divider"></div>
              <div class="warning-content">
                <span class="warning-icon">⚠️</span>
                <p class="warning-text">{{ warning }}</p>
              </div>
            </div>

            <!-- 動作按鈕（可選） -->
            <div v-if="actionLabel" class="bubble-footer">
              <a :href="actionUrl" class="bubble-button" :class="{ 'button-disabled': disabled }">
                {{ actionLabel }}
              </a>
            </div>
          </div>
        </div>

        <!-- 用戶輸入區域（裝飾） -->
        <div class="input-area">
          <div class="input-placeholder">
            <span>輸入訊息...</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Field {
  icon: string
  label: string
  value: string
}

interface Props {
  title?: string
  content?: string
  fields?: Field[]
  warning?: string
  actionLabel?: string
  actionUrl?: string
  showDate?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '🔔 新的通知',
  fields: () => [],
  showDate: true,
  disabled: false,
})

// 格式化時間
const currentTime = computed(() => {
  const now = new Date()
  return now.toLocaleTimeString('zh-TW', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
})

// 格式化日期
const formattedDate = computed(() => {
  const now = new Date()
  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  const month = now.getMonth() + 1
  const day = now.getDate()
  const weekDay = weekDays[now.getDay()]
  return `${month}月${day}日 星期${weekDay}`
})
</script>

<style scoped>
/* 手機框樣式 */
.line-phone-frame {
  width: 320px;
  min-height: 500px;
  background: linear-gradient(180deg, #f5f5f5 0%, #e8e8e8 100%);
  border-radius: 36px;
  padding: 12px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.25),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  position: relative;
  overflow: hidden;
}

/* 劉海/相機區域 */
.line-phone-frame::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 120px;
  height: 28px;
  background: #e8e8e8;
  border-radius: 0 0 16px 16px;
  z-index: 10;
}

/* 狀態列 */
.phone-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px 4px;
  font-size: 12px;
  color: #333;
}

.status-icons {
  display: flex;
  gap: 4px;
  align-items: center;
}

.signal-icon, .wifi-icon {
  width: 16px;
  height: 16px;
}

.battery-icon {
  width: 20px;
  height: 20px;
}

/* LINE 標題列 */
.line-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px 12px;
  border-bottom: 1px solid #e8e8e8;
  background: linear-gradient(180deg, #ffffff 0%, #f9f9f9 100%);
}

.line-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.line-logo {
  width: 28px;
  height: 28px;
}

.official-account {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.line-menu svg {
  width: 24px;
  height: 24px;
  color: #888;
}

/* 訊息容器 */
.message-container {
  padding: 16px;
  min-height: 360px;
  display: flex;
  flex-direction: column;
}

/* 日期標籤 */
.date-label {
  text-align: center;
  font-size: 12px;
  color: #888;
  margin-bottom: 16px;
}

/* Flex Message 氣泡 */
.flex-message-bubble {
  position: relative;
  max-width: 100%;
}

/* 氣泡箭頭 */
.bubble-arrow {
  position: absolute;
  top: 0;
  left: 16px;
  width: 0;
  height: 0;
  border-left: 12px solid transparent;
  border-right: 12px solid transparent;
  border-top: 12px solid #ffffff;
  transform: translateY(-100%);
}

/* 氣泡內容 */
.bubble-content {
  background: #ffffff;
  border-radius: 16px;
  box-shadow:
    0 2px 8px rgba(0, 0, 0, 0.08),
    0 0 0 1px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

/* 氣泡標題區塊 */
.bubble-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 16px 12px;
}

.notification-icon {
  width: 24px;
  height: 24px;
  color: #06C755;
}

.notification-icon svg {
  width: 100%;
  height: 100%;
}

.bubble-title {
  font-size: 16px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0;
  letter-spacing: 0.3px;
}

/* 氣泡分隔線 */
.bubble-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, #e8e8e8 50%, transparent 100%);
  margin: 0 16px;
}

/* 氣泡主體 */
.bubble-body {
  padding: 16px;
}

/* 氣泡文字 */
.bubble-text {
  font-size: 14px;
  color: #444;
  line-height: 1.6;
}

/* 氣泡欄位 */
.bubble-fields {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.bubble-field {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.field-icon {
  width: 20px;
  text-align: center;
  font-size: 14px;
}

.field-label {
  color: #666;
  min-width: 60px;
}

.field-value {
  color: #1a1a1a;
  font-weight: 500;
}

/* 警告區塊 */
.bubble-warning {
  margin-top: 4px;
}

.warning-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, #fbbf24 50%, transparent 100%);
  margin: 0 16px;
}

.warning-content {
  padding: 12px 16px;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.warning-icon {
  font-size: 14px;
  line-height: 1.4;
}

.warning-text {
  font-size: 13px;
  color: #92400e;
  margin: 0;
  line-height: 1.5;
}

/* 氣泡底部按鈕 */
.bubble-footer {
  padding: 12px 16px 16px;
}

.bubble-button {
  display: block;
  text-align: center;
  padding: 12px 24px;
  background: linear-gradient(135deg, #06C755 0%, #05b547 100%);
  color: white;
  text-decoration: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(6, 199, 85, 0.3);
  transition: all 0.3s ease;
}

.bubble-button:hover:not(.button-disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(6, 199, 85, 0.4);
}

.bubble-button:active:not(.button-disabled) {
  transform: translateY(0);
}

.button-disabled {
  background: linear-gradient(135deg, #a8a8a8 0%, #999999 100%);
  cursor: not-allowed;
  box-shadow: none;
}

/* 輸入區域 */
.input-area {
  margin-top: auto;
  padding-top: 12px;
}

.input-placeholder {
  background: #f0f0f0;
  border-radius: 20px;
  padding: 10px 16px;
  color: #888;
  font-size: 14px;
}

/* 深色模式支援 */
@media (prefers-color-scheme: dark) {
  .line-phone-frame {
    background: linear-gradient(180deg, #1a1a1a 0%, #0d0d0d 100%);
    box-shadow:
      0 25px 50px -12px rgba(0, 0, 0, 0.5),
      0 0 0 1px rgba(255, 255, 255, 0.05);
  }

  .line-phone-frame::before {
    background: #0d0d0d;
  }

  .phone-status-bar {
    color: #888;
  }

  .line-header {
    border-bottom-color: #333;
    background: linear-gradient(180deg, #1f1f1f 0%, #1a1a1a 100%);
  }

  .official-account {
    color: #e0e0e0;
  }

  .line-menu svg {
    color: #666;
  }

  .message-container {
    background: #141414;
  }

  .date-label {
    color: #666;
  }

  .bubble-content {
    background: #2a2a2a;
    box-shadow:
      0 2px 8px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.05);
  }

  .bubble-divider {
    background: linear-gradient(90deg, transparent 0%, #444 50%, transparent 100%);
  }

  .bubble-title {
    color: #ffffff;
  }

  .field-label {
    color: #999;
  }

  .field-value {
    color: #e0e0e0;
  }

  .warning-content {
    background: linear-gradient(135deg, #451a03 0%, #78350f 100%);
  }

  .warning-text {
    color: #fcd34d;
  }

  .input-area {
    background: #1a1a1a;
  }

  .input-placeholder {
    background: #2a2a2a;
    color: #666;
  }
}
</style>
