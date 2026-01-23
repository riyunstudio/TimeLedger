import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseButton from '~/components/base/BaseButton.vue'

describe('BaseButton', () => {
  it('renders with default variant (primary)', () => {
    const wrapper = mount(BaseButton, {
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('btn-primary')
    expect(wrapper.text()).toContain('Click me')
  })

  it('applies secondary variant', () => {
    const wrapper = mount(BaseButton, {
      props: {
        variant: 'secondary',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('btn-secondary')
  })

  it('applies success variant', () => {
    const wrapper = mount(BaseButton, {
      props: {
        variant: 'success',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('btn-success')
  })

  it('applies critical variant', () => {
    const wrapper = mount(BaseButton, {
      props: {
        variant: 'critical',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('btn-critical')
  })

  it('applies warning variant', () => {
    const wrapper = mount(BaseButton, {
      props: {
        variant: 'warning',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('btn-warning')
  })

  it('applies sm size', () => {
    const wrapper = mount(BaseButton, {
      props: {
        size: 'sm',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('px-3')
    expect(wrapper.classes()).toContain('py-1.5')
  })

  it('applies md size', () => {
    const wrapper = mount(BaseButton, {
      props: {
        size: 'md',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('px-6')
    expect(wrapper.classes()).toContain('py-3')
  })

  it('applies lg size', () => {
    const wrapper = mount(BaseButton, {
      props: {
        size: 'lg',
      },
      slots: {
        default: 'Click me',
      },
    })

    expect(wrapper.classes()).toContain('px-8')
    expect(wrapper.classes()).toContain('py-4')
  })

  it('emits click event', async () => {
    const wrapper = mount(BaseButton, {
      slots: {
        default: 'Click me',
      },
    })

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toHaveLength(1)
  })

  it('is disabled when disabled prop is true', () => {
    const wrapper = mount(BaseButton, {
      props: {
        disabled: true,
      },
      slots: {
        default: 'Click me',
      },
    })

    const button = wrapper.find('button')
    expect(button.attributes('disabled')).toBeDefined()
    expect(wrapper.classes()).toContain('opacity-50')
  })

  it('does not emit click when disabled', async () => {
    const wrapper = mount(BaseButton, {
      props: {
        disabled: true,
      },
      slots: {
        default: 'Click me',
      },
    })

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toHaveLength(0)
  })
})
