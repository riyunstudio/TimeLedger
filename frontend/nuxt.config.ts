// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },

  modules: [
    '@pinia/nuxt',
    '@nuxtjs/tailwindcss',
    'nuxt-headlessui',
    '@vite-pwa/nuxt',
  ],

  css: ['~/assets/css/main.css'],

  plugins: [
    { src: '~/plugins/virtual-scroller.ts', mode: 'client' },
  ],

  tailwindcss: {
    configPath: 'tailwind.config.js',
  },

  app: {
    buildAssetsDir: '/_nuxt/',
    buildId: 'dev-' + Date.now(),
    head: {
      title: 'TimeLedger - 教師排課平台',
      meta: [
        { name: 'description', content: '教師中心化多據點排課平台' },
        { name: 'theme-color', content: '#312e81' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1, maximum-scale=1' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' },
      ],
    },
  },

  runtimeConfig: {
    public: {
      apiBase: '/api/v1',
      liffId: process.env.LIFF_ID || '',
      lineOfficialAccountId: process.env.LINE_OFFICIAL_ACCOUNT_ID || '@timeledger',
    },
  },

  ssr: false,

  pwa: {
    registerType: 'autoUpdate',
    manifest: {
      name: 'TimeLedger - 教師排課平台',
      short_name: 'TimeLedger',
      description: '教師中心化多據點排課平台',
      theme_color: '#312e81',
      background_color: '#0f172a',
      display: 'standalone',
      orientation: 'portrait',
      scope: '/',
      start_url: '/',
      icons: [
        {
          src: '/pwa/icon-192.png',
          sizes: '192x192',
          type: 'image/png',
        },
        {
          src: '/pwa/icon-512.png',
          sizes: '512x512',
          type: 'image/png',
        },
        {
          src: '/pwa/icon-512.png',
          sizes: '512x512',
          type: 'image/png',
          purpose: 'any maskable',
        },
      ],
    },
    workbox: {
      navigateFallback: '/',
      globPatterns: ['**/*.{js,css,html,png,svg,ico}'],
      globIgnores: [
        '**/node_modules/**/*',
        '**/sw.js',
        '**/workbox-*.js',
        '**/dev-sw-dist/**/*',
      ],
      runtimeCaching: [
        {
          urlPattern: /^https:\/\/api\.timeledger\.app\/api\/v1\/.*$/,
          handler: 'NetworkFirst',
          options: {
            cacheName: 'api-cache',
            expiration: {
              maxEntries: 100,
              maxAgeSeconds: 60 * 60 * 24, // 24 hours
            },
            cacheableResponse: {
              statuses: [0, 200],
            },
          },
        },
        {
          urlPattern: /^https:\/\/.*\.line\.me\/.*$/,
          handler: 'CacheFirst',
          options: {
            cacheName: 'line-cache',
            expiration: {
              maxEntries: 50,
              maxAgeSeconds: 60 * 60 * 24 * 7, // 7 days
            },
          },
        },
      ],
    },
    client: {
      installPrompt: true,
    },
    devOptions: {
      enabled: true,
      type: 'module',
    },
  },

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
