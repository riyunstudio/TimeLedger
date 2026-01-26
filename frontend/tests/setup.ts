// Vitest setup file
// Basic window mocks for jsdom environment

// Mock localStorage
Object.defineProperty(window, 'localStorage', {
  value: {
    getItem: () => null,
    setItem: () => {},
    removeItem: () => {},
    clear: () => {},
  },
  writable: true,
})

// Mock sessionStorage
Object.defineProperty(window, 'sessionStorage', {
  value: {
    getItem: () => null,
    setItem: () => {},
    removeItem: () => {},
    clear: () => {},
  },
  writable: true,
})

// Mock scrollTo
window.scrollTo = () => {}

// Mock $alert for toast fallback
window.$alert = () => {}

// Mock ResizeObserver
global.ResizeObserver = vi?.fn?.()?.mockImplementation(() => ({
  observe: () => {},
  unobserve: () => {},
  disconnect: () => {},
})) || function() {}

// Mock IntersectionObserver  
global.IntersectionObserver = vi?.fn?.()?.mockImplementation(() => ({
  observe: () => {},
  unobserve: () => {},
  disconnect: () => {},
})) || function() {}

// Global component stubs for Transition and Teleport
import { config } from '@vue/test-utils'

config.global.stubs = {
  Transition: {
    template: '<slot />',
  },
  Teleport: {
    template: '<slot />',
  },
}

console.log('Vitest setup loaded')
