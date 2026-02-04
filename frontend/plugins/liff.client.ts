import liff from '@line/liff'

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()

  try {
    if (!config.public.liffId) {
      throw new Error('LIFF_ID is missing in runtimeConfig. Please check your .env.local file.')
    }

    // 初始化 LIFF
    await liff.init({
      liffId: config.public.liffId,
    })

    // 注入 liff 實例到 Nuxt App
    return {
      provide: {
        liff,
      },
    }
  } catch (error) {
    console.error('LIFF init failed:', error)
    // 返回空的 liff 物件，避免後續使用時出錯
    return {
      provide: {
        liff: null,
      },
    }
  }
})
