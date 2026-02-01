import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    include: [
      'tests/**/*.test.ts',
      'tests/**/*.spec.ts',
      'tests/components/**/*.spec.ts',
      'tests/pages/**/*.spec.ts',
      'tests/layouts/**/*.spec.ts',
      'tests/bench/**/*.test.ts',
    ],
    exclude: [
      'tests/e2e.spec.ts',
      'tests/login-flow.spec.ts',
      'tests/approval-flow.spec.ts',
      'tests/admin-course-test.spec.ts',
    ],
    benchmark: {
      enabled: true,
      include: ['tests/bench/**/*.test.ts'],
      reporter: ['text', 'json', 'html'],
      outputFile: {
        text: 'tests/benchmark-results.txt',
        json: 'tests/benchmark-results.json',
        html: 'tests/benchmark-results.html',
      },
    },
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'tests/',
        '**/*.d.ts',
        '**/*.vue',
      ],
    },
    setupFiles: ['tests/setup.ts'],
    deps: {
      inline: ['vue', 'pinia', '@vueuse/core'],
    },
    transformMode: {
      web: [/\.vue$/],
    },
  },
  resolve: {
    alias: {
      '~': resolve(__dirname, './'),
      '@': resolve(__dirname, './'),
    },
  },
  define: {
    'process.env': {},
  },
})
