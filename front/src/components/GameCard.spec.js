import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import GameCard from './GameCard.vue'

const game = {
  id: 1,
  title: 'Test Game',
  is_finished: false,
  created_at: '2024-01-01T00:00:00Z',
}

describe('GameCard', () => {
  it('mounts and renders the game title', () => {
    const wrapper = mount(GameCard, { props: { game } })
    expect(wrapper.text()).toContain('Test Game')
  })

  it('shows "In progress" for an unfinished game', () => {
    const wrapper = mount(GameCard, { props: { game } })
    expect(wrapper.text()).toContain('In progress')
  })

  it('shows "Finished" for a finished game', () => {
    const wrapper = mount(GameCard, { props: { game: { ...game, is_finished: true } } })
    expect(wrapper.text()).toContain('Finished')
  })
})
