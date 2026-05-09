import { createRouter, createWebHistory } from 'vue-router'
import GamesView from './views/GamesView.vue'
import GameView from './views/GameView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: GamesView },
    { path: '/games/:id', component: GameView },
  ],
})
