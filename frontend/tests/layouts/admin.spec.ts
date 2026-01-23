import { describe, it, expect } from 'vitest'
import { render } from '@vue/test-utils'

import AdminLayout from '~/layouts/admin.vue'

describe('Admin Layout', () => {
  beforeEach(() => {
    const pinia = createPinia()
    setActivePinia(pinia)
  })

  it('renders layout structure', () => {
    const { html } = render(AdminLayout)

    expect(html()).toContain('min-h-screen')
    expect(html()).toContain('bg-slate-900')
  })

  it('renders admin header', () => {
    const { html } = render(AdminLayout, {
      global: {
        components: {
          AdminHeader: {
            template: '<div class="test-admin-header">Admin Header</div>',
          },
        },
      },
    })

    expect(html()).toContain('Admin Header')
  })

  it('renders main content area', () => {
    const { html } = render(AdminLayout, {
      slots: {
        default: '<div class="test-content">Admin Page</div>',
      },
    })

    expect(html()).toContain('Admin Page')
  })

  it('renders admin sidebar when showSidebar is true', () => {
    const { html } = render(AdminLayout, {
      props: {
        showSidebar: true,
      },
      global: {
        components: {
          AdminSidebar: {
            template: '<div class="test-sidebar">Admin Sidebar</div>',
          },
        },
      },
    })

    expect(html()).toContain('Admin Sidebar')
  })

  it('does not render admin sidebar when showSidebar is false', () => {
    const { html } = render(AdminLayout, {
      props: {
        showSidebar: false,
      },
      global: {
        components: {
          AdminSidebar: {
            template: '<div class="test-sidebar">Admin Sidebar</div>',
          },
        },
      },
    })

    expect(html()).not.toContain('Admin Sidebar')
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
