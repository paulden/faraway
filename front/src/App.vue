<script setup>
import { ref, onMounted } from 'vue'

const API = import.meta.env.VITE_API_URL ?? '/api'

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

      <div v-else-if="games.length === 0" class="text-slate-400 text-sm">
        No games yet.
      </div>

      <div v-else class="flex flex-col gap-3">
        <div
          v-for="game in games"
          :key="game.id"
          class="bg-white rounded-2xl px-5 py-4 shadow-sm flex items-center justify-between gap-4"
        >
          <span class="font-medium text-slate-600">{{ game.title }}</span>
          <span
            :class="game.is_finished
              ? 'bg-green-100 text-green-600'
              : 'bg-violet-100 text-violet-500'"
            class="text-xs font-medium px-3 py-1 rounded-full shrink-0"
          >
            {{ game.is_finished ? 'Finished' : 'In progress' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
