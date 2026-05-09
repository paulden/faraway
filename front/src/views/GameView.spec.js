import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createRouter, createMemoryHistory } from 'vue-router'
import GameView from './GameView.vue'

const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/', component: { template: '<div/>' } },
    { path: '/games/:id', component: GameView },
  ],
})

const mockGame = { id: 1, title: 'Poker Night', is_finished: false }

function stubFetch({ game = mockGame, players = [], rounds = [] } = {}) {
  vi.stubGlobal('fetch', vi.fn((url) => {
    if (url.endsWith('/players')) return Promise.resolve({ json: () => Promise.resolve(players) })
    if (url.endsWith('/rounds')) return Promise.resolve({ json: () => Promise.resolve(rounds) })
    return Promise.resolve({ json: () => Promise.resolve(game) })
  }))
}

describe('GameView', () => {
  beforeEach(async () => {
    stubFetch()
    await router.push('/games/1')
    await router.isReady()
  })

  it('mounts without errors', async () => {
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.exists()).toBe(true)
  })

  it('renders the game title after loading', async () => {
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('Poker Night')
  })

  it('shows "In progress" badge for an active game', async () => {
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('In progress')
  })

  it('shows "No players yet." when there are no players', async () => {
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('No players yet.')
  })

  it('shows action buttons for an in-progress game', async () => {
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('Add round')
    expect(wrapper.text()).toContain('Add player')
    expect(wrapper.text()).toContain('Finish game')
  })

  it('hides action buttons for a finished game', async () => {
    stubFetch({ game: { ...mockGame, is_finished: true } })
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).not.toContain('Finish game')
  })

  it('renders player names in the scoreboard header', async () => {
    stubFetch({
      players: [{ id: 10, name: 'Alice' }, { id: 11, name: 'Bob' }],
      rounds: [{ id: 100, number: 1 }],
    })
    vi.stubGlobal('fetch', vi.fn((url) => {
      if (url.endsWith('/players')) return Promise.resolve({ json: () => Promise.resolve([{ id: 10, name: 'Alice' }, { id: 11, name: 'Bob' }]) })
      if (url.endsWith('/rounds')) return Promise.resolve({ json: () => Promise.resolve([{ id: 100, number: 1 }]) })
      if (url.includes('/scores')) return Promise.resolve({ json: () => Promise.resolve([]) })
      return Promise.resolve({ json: () => Promise.resolve(mockGame) })
    }))
    const wrapper = mount(GameView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('Alice')
    expect(wrapper.text()).toContain('Bob')
  })
})
