import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'

describe('BaseGlassCard', () => {
  it('renders with default classes', () => {
    const wrapper = mount(BaseGlassCard, {
      slots: {
        default: 'Test Content',
      },
    })

    expect(wrapper.classes()).toContain('glass-card')
    expect(wrapper.text()).toContain('Test Content')
  })

  it('applies custom padding', () => {
    const wrapper = mount(BaseGlassCard, {
      props: {
        padding: 'p-8',
      },
      slots: {
        default: 'Test Content',
      },
    })

    expect(wrapper.classes()).toContain('p-8')
  })

  it('applies custom margin', () => {
    const wrapper = mount(BaseGlassCard, {
      props: {
        margin: 'm-4',
      },
      slots: {
        default: 'Test Content',
      },
    })

    expect(wrapper.classes()).toContain('m-4')
  })
})
