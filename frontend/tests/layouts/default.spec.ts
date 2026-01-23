import { describe, it, expect } from 'vitest'
import { render } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'

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