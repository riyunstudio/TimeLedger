import { describe, it, expect } from 'vitest'
import { render } from '@vue/test-utils'

import BlankLayout from '~/layouts/blank.vue'

describe('Blank Layout', () => {
  beforeEach(() => {
    const pinia = createPinia()
    setActivePinia(pinia)
  })

  it('renders layout structure', () => {
    const { html } = render(BlankLayout)

    expect(html()).toContain('min-h-screen')
  })

  it('renders with gradient background', () => {
    const { html } = render(BlankLayout)

    expect(html()).toContain('bg-gradient-mesh')
  })

  it('renders content in center', () => {
    const { html } = render(BlankLayout, {
      slots: {
        default: '<div class="test-content">Test Page</div>',
      },
    })

    expect(html()).toContain('Test Page')
    expect(html()).toContain('min-h-screen')
    expect(html()).toContain('items-center')
    expect(html()).toContain('justify-center')
  })

  it('renders content with padding', () => {
    const { html } = render(BlankLayout, {
      slots: {
        default: '<div class="test-content">Test Page</div>',
      },
    })

    expect(html()).toContain('p-4')
  })
})
