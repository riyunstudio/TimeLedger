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

  tailwindcss: {
    configPath: 'tailwind.config.js',
  },

  app: {
    buildAssetsDir: '/_nuxt/',
    buildId: 'dev-' + Date.now(),
  },

  runtimeConfig: {
    public: {
      apiBase: '/api/v1',
      liffId: process.env.LIFF_ID || '',
      lineOfficialAccountId: process.env.LINE_OFFICIAL_ACCOUNT_ID || '@timeledger',
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
    // Fix Windows path issues
    resolve: {
      alias: {
        // Force relative paths
      }
    },
    // Disable problematic optimizations for Windows
    optimizeDeps: {
      include: [],
    },
    // Fix module preload polyfill on Windows
    modulePreload: {
      polyfill: false,
    },
  },
})
