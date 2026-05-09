import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createRouter, createMemoryHistory } from 'vue-router'
import GamesView from './GamesView.vue'

const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/', component: GamesView },
    { path: '/games/:id', component: { template: '<div/>' } },
  ],
})

describe('GamesView', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn(() =>
      Promise.resolve({ json: () => Promise.resolve([]) })
    ))
  })

  it('mounts and shows loading state initially', () => {
    const wrapper = mount(GamesView, { global: { plugins: [router] } })
    expect(wrapper.text()).toContain('Loading games')
  })

  it('shows "No games yet." when the list is empty', async () => {
    const wrapper = mount(GamesView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('No games yet.')
  })

  it('renders game cards when games are returned', async () => {
    vi.stubGlobal('fetch', vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([
          { id: 1, title: 'Poker Night', is_finished: false, created_at: '2024-01-01T00:00:00Z' },
        ]),
      })
    ))
    const wrapper = mount(GamesView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('Poker Night')
  })

  it('shows an error message when fetch fails', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.reject(new Error('Network error'))))
    const wrapper = mount(GamesView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('Could not load games.')
  })
})
