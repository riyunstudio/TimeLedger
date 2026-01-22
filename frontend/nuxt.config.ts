// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },

  modules: [
    '@pinia/nuxt',
    '@nuxtjs/tailwindcss',
    'nuxt-headlessui',
  ],

  css: ['~/assets/css/main.css'],

  app: {
    buildAssetsDir: '/_nuxt/',
    buildId: 'dev-' + Date.now(),
  },

  runtimeConfig: {
    public: {
      apiBase: '/api/v1',
      liffId: process.env.LIFF_ID || '',
    },
  },

  ssr: false,

  vite: {
    server: {
      hmr: {
        overlay: false,
      },
      proxy: {
        '/api': {
          target: process.env.API_TARGET || 'http://localhost:8888',
          changeOrigin: true,
          secure: false,
        },
      },
    },
  },
})
