import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import VirtualScroll from '~/components/base/VirtualScroll.vue'

describe('VirtualScroll', () => {
  const testItems = Array.from({ length: 100 }, (_, i) => ({
    id: i + 1,
    name: `Item ${i + 1}`,
  }))

  const itemHeight = 50
  const containerHeight = 300

  it('should render nothing when items is empty', () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: [],
        itemHeight,
      },
    })

    expect(wrapper.find('.virtual-scroll-container').exists()).toBe(true)
    expect(wrapper.findAll('.virtual-scroll-item')).toHaveLength(0)
  })

  it('should render visible items only', async () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    // Set container height
    await wrapper.find('.virtual-scroll-container').element.setAttribute('style', `height: ${containerHeight}px`)
    await wrapper.vm.$nextTick()

    // Should render approximately containerHeight / itemHeight + overscan items
    const visibleItems = wrapper.findAll('.virtual-scroll-item')
    expect(visibleItems.length).toBeGreaterThan(0)
    expect(visibleItems.length).toBeLessThan(testItems.length)
  })

  it('should emit scroll event', async () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    const container = wrapper.find('.virtual-scroll-container')
    
    // Trigger scroll event without trying to set target
    await container.trigger('scroll')

    expect(wrapper.emitted('scroll')).toBeTruthy()
  })

  it('should have total height based on items', () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    // Check spacer element has correct height
    const spacer = wrapper.find('.virtual-scroll-spacer')
    expect(spacer.attributes('style')).toContain(`height: ${testItems.length * itemHeight}px`)
  })

  it('should expose scrollToIndex method', async () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    const vm = wrapper.vm as any
    expect(typeof vm.scrollToIndex).toBe('function')
    expect(typeof vm.scrollToTop).toBe('function')
    expect(typeof vm.scrollToBottom).toBe('function')
    expect(typeof vm.refresh).toBe('function')
  })

  it('should have proper ARIA attributes', () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
        ariaLabel: 'Custom List',
      },
    })

    const container = wrapper.find('.virtual-scroll-container')
    expect(container.attributes('role')).toBe('listbox')
    expect(container.attributes('aria-label')).toBe('Custom List')
  })

  it('should have default ARIA label when not provided', () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    const container = wrapper.find('.virtual-scroll-container')
    expect(container.attributes('aria-label')).toBe('虛擬滾動列表')
  })

  it('should apply custom item key', () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
        itemKey: 'name',
      },
    })

    const items = wrapper.findAll('.virtual-scroll-item')
    expect(items.length).toBeGreaterThan(0)
  })

  it('should render items with correct structure', async () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    await wrapper.find('.virtual-scroll-container').element.setAttribute('style', `height: ${containerHeight}px`)
    await wrapper.vm.$nextTick()

    const items = wrapper.findAll('.virtual-scroll-item')
    
    // Each item should have role="option"
    const options = wrapper.findAll('.virtual-scroll-item[role="option"]')
    expect(options.length).toBe(items.length)
  })

  it('should render first few items initially', async () => {
    const wrapper = mount(VirtualScroll, {
      props: {
        items: testItems,
        itemHeight,
      },
    })

    await wrapper.find('.virtual-scroll-container').element.setAttribute('style', `height: ${containerHeight}px`)
    await wrapper.vm.$nextTick()

    // First items should be visible
    const items = wrapper.findAll('.virtual-scroll-item')
    expect(items.length).toBeGreaterThan(0)
  })
})
