<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import GameCard from '../components/GameCard.vue'
import CreateGame from '../components/CreateGame.vue'

const API = import.meta.env.VITE_API_URL ?? '/api'
const router = useRouter()
const games = ref([])
const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  try {
    const res = await fetch(`${API}/games`)
    games.value = await res.json()
  } catch (e) {
    error.value = 'Could not load games.'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen px-4 py-10">
    <div class="max-w-2xl mx-auto">
      <h1 class="text-3xl font-semibold text-slate-600 mb-8 tracking-wide">Faraway</h1>

      <p v-if="loading" class="text-slate-400 text-sm">Loading games…</p>
      <p v-else-if="error" class="text-red-400 text-sm">{{ error }}</p>
      <p v-else-if="games.length === 0" class="text-slate-400 text-sm">No games yet.</p>

      <div v-else class="flex flex-col gap-3">
        <GameCard
          v-for="game in games"
          :key="game.id"
          :game="game"
          @click="router.push(`/games/${game.id}`)"
        />
      </div>
    </div>
  </div>

  <CreateGame @created="games.unshift($event)" />
</template>
