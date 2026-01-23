import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseBadge from '~/components/base/BaseBadge.vue'

describe('BaseBadge', () => {
  it('renders with default variant (primary)', () => {
    const wrapper = mount(BaseBadge, {
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('bg-primary-500/20')
    expect(wrapper.classes()).toContain('text-primary-300')
    expect(wrapper.text()).toContain('Test Badge')
  })

  it('applies secondary variant', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        variant: 'secondary',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('bg-secondary-500/20')
    expect(wrapper.classes()).toContain('text-secondary-300')
  })

  it('applies success variant', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        variant: 'success',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('bg-success-500/20')
    expect(wrapper.classes()).toContain('text-success-300')
  })

  it('applies critical variant', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        variant: 'critical',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('bg-critical-500/20')
    expect(wrapper.classes()).toContain('text-critical-300')
  })

  it('applies warning variant', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        variant: 'warning',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('bg-warning-500/20')
    expect(wrapper.classes()).toContain('text-warning-300')
  })

  it('applies info variant', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        variant: 'info',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('bg-blue-500/20')
    expect(wrapper.classes()).toContain('text-blue-300')
  })

  it('applies sm size', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        size: 'sm',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('px-3')
    expect(wrapper.classes()).toContain('py-1')
  })

  it('applies md size (default)', () => {
    const wrapper = mount(BaseBadge, {
      props: {
        size: 'md',
      },
      slots: {
        default: 'Test Badge',
      },
    })

    expect(wrapper.classes()).toContain('px-3')
    expect(wrapper.classes()).toContain('py-1')
  })
})
