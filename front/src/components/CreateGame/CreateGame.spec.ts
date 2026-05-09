import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import CreateGame from './CreateGame.vue'

describe('CreateGame', () => {
  it('mounts and renders the open button', () => {
    const wrapper = mount(CreateGame)
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('modal is hidden by default', () => {
    const wrapper = mount(CreateGame)
    expect(wrapper.text()).not.toContain('New game')
  })

  it('shows the modal after clicking the button', async () => {
    const wrapper = mount(CreateGame)
    await wrapper.find('button').trigger('click')
    expect(wrapper.text()).toContain('New game')
  })

  it('hides the modal after clicking Cancel', async () => {
    const wrapper = mount(CreateGame)
    await wrapper.find('button').trigger('click')
    await wrapper.find('button[class*="text-slate-400"]').trigger('click')
    expect(wrapper.text()).not.toContain('New game')
  })
})
