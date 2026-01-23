import { describe, it, expect } from 'vitest'
import { render } from '@vue/test-utils'

import AdminLayout from '~/layouts/admin.vue'

// Mock Components
const AdminSidebar = {
  template: '<div class="test-admin-sidebar">Admin Sidebar</div>',
}

describe('Admin Layout', () => {
  it('renders layout structure', () => {
    const { html } = render(AdminLayout, {
      global: {
        components: {
          AdminSidebar,
        },
      },
    })

    expect(html()).toContain('min-h-screen')
    expect(html()).toContain('bg-slate-900')
  })

  it('renders sidebar', () => {
    const { html } = render(AdminLayout, {
      global: {
        components: {
          AdminSidebar,
        },
      },
    })

    expect(html()).toContain('Admin Sidebar')
  })

  it('renders main content with padding', () => {
    const { html } = render(AdminLayout, {
      slots: {
        default: '<div class="test-content">Test Content</div>',
      },
    })

    expect(html()).toContain('p-4')
    expect(html()).toContain('md:p-6')
  })
})