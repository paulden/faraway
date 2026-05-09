import { createRouter, createWebHistory } from 'vue-router'
import GamesView from './views/GamesView/GamesView.vue'
import GameView from './views/GameView/GameView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: GamesView },
    { path: '/games/:id', component: GameView },
  ],
})
