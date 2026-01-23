import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GridSkeleton from '~/components/base/GridSkeleton.vue'

describe('GridSkeleton', () => {
  it('renders with default rows (5)', () => {
    const wrapper = mount(GridSkeleton)

    expect(wrapper.findAll('.space-y-4 > div').length).toBe(5)
  })

  it('renders with custom rows', () => {
    const wrapper = mount(GridSkeleton, {
      props: {
        rows: 3,
      },
    })

    expect(wrapper.findAll('.space-y-4 > div').length).toBe(3)
  })

  it('renders avatar with default class', () => {
    const wrapper = mount(GridSkeleton)

    expect(wrapper.find('.w-12.h-12').exists()).toBe(true)
    expect(wrapper.find('.flex-shrink-0').exists()).toBe(true)
  })

  it('renders avatar with custom class', () => {
    const wrapper = mount(GridSkeleton, {
      props: {
        avatarClass: 'w-8 h-8',
      },
    })

    expect(wrapper.find('.w-8.h-8').exists()).toBe(true)
  })

  it('renders title with default width (40%)', () => {
    const wrapper = mount(GridSkeleton)

    const titles = wrapper.findAll('.flex-1.space-y-2 > div:nth-child(1)')
    expect(titles.length).toBeGreaterThan(0)
  })

  it('renders subtitle with default width (70%)', () => {
    const wrapper = mount(GridSkeleton)

    const subtitles = wrapper.findAll('.flex-1.space-y-2 > div:nth-child(2)')
    expect(subtitles.length).toBeGreaterThan(0)
  })

  it('renders with custom title width', () => {
    const wrapper = mount(GridSkeleton, {
      props: {
        titleWidth: '60%',
      },
    })

    expect(wrapper.html()).toContain('width: 60%')
  })

  it('renders with custom subtitle width', () => {
    const wrapper = mount(GridSkeleton, {
      props: {
        subtitleWidth: '80%',
      },
    })

    expect(wrapper.html()).toContain('width: 80%')
  })

  it('has animation class', () => {
    const wrapper = mount(GridSkeleton)

    expect(wrapper.find('.animate-pulse').exists()).toBe(true)
  })

  it('renders skeleton bars', () => {
    const wrapper = mount(GridSkeleton)

    expect(wrapper.findAll('.rounded-full').length).toBeGreaterThan(0)
    expect(wrapper.findAll('.rounded-xl').length).toBeGreaterThan(0)
  })
})
