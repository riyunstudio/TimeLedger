import liff from '@line/liff'

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()

  try {
    // 檢查 LIFF_ID 是否存在
    if (!config.public.liffId || config.public.liffId === 'YOUR_LIFF_ID_HERE') {
      console.error('LIFF 初始化失敗: LIFF_ID 未正確設定')
      console.error('請檢查 .env.local 文件中的 NUXT_PUBLIC_LIFF_ID 設定')

      // 返回 null lixx 實例
      return {
        provide: {
          liff: null,
        },
      }
    }

    // 初始化 LIFF
    await liff.init({
      liffId: config.public.liffId,
    })

    console.log('[LIFF] 初始化成功，LIFF ID:', config.public.liffId)

    // 注入 liff 實例到 Nuxt App
    return {
      provide: {
        liff,
      },
    }
  } catch (error: any) {
    console.error('LIFF 初始化失敗:', error.message || error)

    // 提供詳細的錯誤資訊
    if (error.message?.includes('null')) {
      console.warn('LIFF SDK 可能未正確載入，請確保在 LINE App 環境中執行')
    }

    // 返回空的 liff 物件，避免後續使用時出錯
    return {
      provide: {
        liff: null,
      },
    }
  }
})
