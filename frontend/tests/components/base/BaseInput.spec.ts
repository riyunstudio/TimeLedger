import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseInput from '~/components/base/BaseInput.vue'

describe('BaseInput', () => {
  it('renders with default props', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
      },
    })

    expect(wrapper.find('input').exists()).toBe(true)
  })

  it('binds modelValue', async () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: 'test',
      },
    })

    const input = wrapper.find('input')
    await input.setValue('new value')
    expect(wrapper.emitted('update:modelValue')).toHaveLength(1)
    expect(wrapper.emitted('update:modelValue')[0]).toEqual(['new value'])
  })

  it('shows label when provided', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        label: 'Test Label',
      },
    })

    expect(wrapper.text()).toContain('Test Label')
  })

  it('shows required asterisk when required', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        label: 'Required Field',
        required: true,
      },
    })

    expect(wrapper.text()).toContain('*')
  })

  it('shows error message', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        error: 'This field is required',
      },
    })

    expect(wrapper.text()).toContain('This field is required')
  })

  it('shows helper text', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        helper: 'Enter your email address',
      },
    })

    expect(wrapper.text()).toContain('Enter your email address')
  })

  it('applies sm size', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        size: 'sm',
      },
    })

    expect(wrapper.find('input').classes()).toContain('px-3')
    expect(wrapper.find('input').classes()).toContain('py-1.5')
  })

  it('applies md size', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        size: 'md',
      },
    })

    expect(wrapper.find('input').classes()).toContain('px-4')
    expect(wrapper.find('input').classes()).toContain('py-3')
  })

  it('applies lg size', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: '',
        size: 'lg',
      },
    })

    expect(wrapper.find('input').classes()).toContain('px-5')
    expect(wrapper.find('input').classes()).toContain('py-4')
  })

  it('is disabled when disabled prop is true', () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: 'test',
        disabled: true,
      },
    })

    const input = wrapper.find('input')
    expect(input.attributes('disabled')).toBeDefined()
  })

  it('emits blur event', async () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: 'test',
      },
    })

    await wrapper.find('input').trigger('blur')
    expect(wrapper.emitted('blur')).toHaveLength(1)
  })

  it('emits focus event', async () => {
    const wrapper = mount(BaseInput, {
      props: {
        modelValue: 'test',
      },
    })

    await wrapper.find('input').trigger('focus')
    expect(wrapper.emitted('focus')).toHaveLength(1)
  })
})
