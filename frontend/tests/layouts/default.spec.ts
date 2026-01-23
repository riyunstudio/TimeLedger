import { describe, it, expect } from 'vitest'
import { render } from '@vue/test-utils'

import DefaultLayout from '~/layouts/default.vue'

// Mock Components
const TeacherHeader = {
  template: '<div class="test-teacher-header">Teacher Header</div>',
}

// Mock Composables
const notificationUI = {
  show: ref(false),
  close: () => notificationUI.show.value = false,
  toggle: () => notificationUI.show.value = !notificationUI.show.value,
}

describe('Default Layout', () => {
  beforeEach(() => {
    const pinia = createPinia()
    setActivePinia(pinia)
  })

  it('renders layout structure', () => {
    const { html } = render(DefaultLayout, {
      global: {
        components: {
          TeacherHeader,
        },
      },
    })

    expect(html()).toContain('min-h-screen')
    expect(html()).toContain('bg-slate-900')
  })

  it('renders teacher header when showHeader is true', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showHeader: true,
      },
      global: {
        components: {
          TeacherHeader,
        },
      },
    })

    expect(html()).toContain('Teacher Header')
  })

  it('does not render teacher header when showHeader is false', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showHeader: false,
      },
      global: {
        components: {
          TeacherHeader,
        },
      },
    })

    expect(html()).not.toContain('Teacher Header')
  })

  it('renders bottom navigation when showNav is true', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showNav: true,
      },
      global: {
        components: {
          TeacherHeader,
        },
      },
    })

    expect(html()).toContain('nav.fixed.bottom-0')
  })

  it('does not render bottom navigation when showNav is false', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showNav: false,
      },
      global: {
        components: {
          TeacherHeader,
        },
      },
    })

    expect(html()).not.toContain('nav.fixed.bottom-0')
  })
})

  it('renders layout structure', () => {
    const { html } = render(DefaultLayout)

    expect(html()).toContain('min-h-screen')
    expect(html()).toContain('bg-slate-900')
  })

  it('renders main content area', () => {
    const { html } = render(DefaultLayout, {
      slots: {
        default: '<div class="test-content">Test Page</div>',
      },
    })

    expect(html()).toContain('Test Page')
    expect(html()).toContain('p-4')
  })

  it('has bottom navigation for mobile', () => {
    const { html } = render(DefaultLayout)

    const bottomNav = html().match(/<nav[^>]*fixed bottom-0[^>]*><\/nav>/s)
    expect(bottomNav).toBeTruthy()
  })

  it('bottom navigation has correct links', () => {
    const { html } = render(DefaultLayout)

    const links = html().match(/<a[^>]*href="([^"]*)"[^>]*>/g)
    expect(links).toBeTruthy()
    expect(links[1]).toContain('/teacher/dashboard')
    expect(links[2]).toContain('/teacher/profile')
  })

  it('renders teacher header when showHeader is true', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showHeader: true,
      },
      global: {
        components: {
          TeacherHeader: {
            template: '<div class="test-header">Teacher Header</div>',
          },
        },
      },
    })

    expect(html()).toContain('Teacher Header')
  })

  it('does not render teacher header when showHeader is false', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showHeader: false,
      },
      global: {
        components: {
          TeacherHeader: {
            template: '<div class="test-header">Teacher Header</div>',
          },
        },
      },
    })

    expect(html()).not.toContain('Teacher Header')
  })

  it('renders bottom navigation when showNav is true', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showNav: true,
      },
    })

    expect(html()).toContain('nav.fixed.bottom-0')
  })

  it('does not render bottom navigation when showNav is false', () => {
    const { html } = render(DefaultLayout, {
      props: {
        showNav: false,
      },
    })

    expect(html()).not.toContain('nav.fixed.bottom-0')
  })
})
