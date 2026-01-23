import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseModal from '~/components/base/BaseModal.vue'

describe('BaseModal', () => {
  it('does not render when modelValue is false', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: false,
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.find('.fixed.inset-0').exists()).toBe(false)
  })

  it('renders when modelValue is true', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.find('.fixed.inset-0').exists()).toBe(true)
  })

  it('shows title when provided', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        title: 'Test Modal Title',
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.text()).toContain('Test Modal Title')
  })

  it('applies sm size', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        size: 'sm',
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.classes()).toContain('w-full')
    expect(wrapper.classes()).toContain('max-w-sm')
  })

  it('applies md size (default)', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        size: 'md',
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.classes()).toContain('w-full')
    expect(wrapper.classes()).toContain('max-w-md')
  })

  it('applies lg size', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        size: 'lg',
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.classes()).toContain('w-full')
    expect(wrapper.classes()).toContain('max-w-lg')
  })

  it('applies xl size', () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        size: 'xl',
      },
      slots: {
        default: 'Modal Content',
      },
    })

    expect(wrapper.classes()).toContain('w-full')
    expect(wrapper.classes()).toContain('max-w-2xl')
  })

  it('emits update:modelValue when clicking backdrop and closeOnBackdrop is true', async () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        closeOnBackdrop: true,
      },
      slots: {
        default: 'Modal Content',
      },
    })

    await wrapper.find('.absolute.inset-0').trigger('click')
    expect(wrapper.emitted('update:modelValue')).toHaveLength(1)
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
  })

  it('does not emit update:modelValue when clicking backdrop and closeOnBackdrop is false', async () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
        closeOnBackdrop: false,
      },
      slots: {
        default: 'Modal Content',
      },
    })

    await wrapper.find('.absolute.inset-0').trigger('click')
    expect(wrapper.emitted('update:modelValue')).toHaveLength(0)
  })

  it('emits update:modelValue when clicking close button', async () => {
    const wrapper = mount(BaseModal, {
      props: {
        modelValue: true,
      },
      slots: {
        default: 'Modal Content',
      },
    })

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('update:modelValue')).toHaveLength(1)
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
  })
})
