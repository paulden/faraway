<script setup>
import { ref, onMounted, nextTick, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ChevronLeft, Plus } from 'lucide-vue-next'

const API = import.meta.env.VITE_API_URL ?? '/api'
const route = useRoute()
const id = route.params.id

const game = ref(null)
const players = ref([])
const rounds = ref([])
const scores = ref({}) // scores[roundId][playerId] = { id, score } | undefined

onMounted(async () => {
  const [gameRes, playersRes, roundsRes] = await Promise.all([
    fetch(`${API}/games/${id}`),
    fetch(`${API}/games/${id}/players`),
    fetch(`${API}/games/${id}/rounds`),
  ])
  game.value = await gameRes.json()
  players.value = await playersRes.json()
  rounds.value = await roundsRes.json()
  await loadScores()
})

async function loadScores() {
  if (rounds.value.length === 0) return
  const results = await Promise.all(
    rounds.value.map(r =>
      fetch(`${API}/games/${id}/rounds/${r.id}/scores`).then(r => r.json())
    )
  )
  const map = {}
  rounds.value.forEach((round, i) => {
    map[round.id] = {}
    results[i].forEach(s => { map[round.id][s.player_id] = { id: s.id, score: s.score } })
  })
  scores.value = map
}

const totals = computed(() =>
  Object.fromEntries(
    players.value.map(p => [p.id, Object.values(scores.value).reduce((sum, row) => sum + (row[p.id]?.score ?? 0), 0)])
  )
)

const leaderId = computed(() => {
  const entries = Object.entries(totals.value)
  if (entries.every(([, v]) => v === 0)) return null
  return entries.reduce((a, b) => b[1] > a[1] ? b : a)[0]
})

// Editing
const editing = ref(null) // { roundId, playerId }
const editValue = ref('')
const editInput = ref(null)

async function startEdit(roundId, playerId) {
  if (game.value?.is_finished) return
  editing.value = { roundId, playerId }
  editValue.value = scores.value[roundId]?.[playerId]?.score ?? ''
  await nextTick()
  editInput.value?.focus()
  editInput.value?.select()
}

async function saveEdit() {
  if (!editing.value) return
  const { roundId, playerId } = editing.value
  const parsed = parseInt(editValue.value)
  if (isNaN(parsed)) { editing.value = null; return }

  const existing = scores.value[roundId]?.[playerId]

  if (existing) {
    const res = await fetch(`${API}/games/${id}/rounds/${roundId}/scores/${existing.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ score: parsed }),
    })
    const updated = await res.json()
    scores.value[roundId][playerId] = { id: updated.id, score: updated.score }
  } else {
    const res = await fetch(`${API}/games/${id}/rounds/${roundId}/scores`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ player_id: playerId, score: parsed }),
    })
    const created = await res.json()
    if (!scores.value[roundId]) scores.value[roundId] = {}
    scores.value[roundId][playerId] = { id: created.id, score: created.score }
  }

  editing.value = null
}

function cancelEdit() {
  editing.value = null
}

// Remove player
async function removePlayer(player) {
  if (!confirm(`Remove ${player.name} from this game?`)) return
  await fetch(`${API}/games/${id}/players/${player.id}`, { method: 'DELETE' })
  players.value = players.value.filter(p => p.id !== player.id)
}

// Add player
const newPlayerName = ref('')
const showAddPlayer = ref(false)
const addingPlayer = ref(false)

async function addPlayer() {
  if (!newPlayerName.value.trim()) return
  addingPlayer.value = true
  try {
    const res = await fetch(`${API}/games/${id}/players`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: newPlayerName.value.trim() }),
    })
    players.value.push(await res.json())
    newPlayerName.value = ''
    showAddPlayer.value = false
  } finally {
    addingPlayer.value = false
  }
}

// Finish game
async function finishGame() {
  if (!confirm('Mark this game as finished?')) return
  const res = await fetch(`${API}/games/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title: game.value.title, is_finished: true }),
  })
  game.value = await res.json()
}

// Add round
async function addRound() {
  const next = rounds.value.length > 0
    ? Math.max(...rounds.value.map(r => r.number)) + 1
    : 1
  const res = await fetch(`${API}/games/${id}/rounds`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ number: next }),
  })
  const round = await res.json()
  rounds.value.push(round)
  scores.value[round.id] = {}
}
</script>

<template>
  <div class="min-h-screen px-4 py-10">
    <div class="max-w-4xl mx-auto">

      <!-- Header -->
      <div class="mb-8">
        <router-link to="/" class="inline-flex items-center gap-1 text-sm text-slate-400 hover:text-slate-600 transition-colors mb-4">
          <ChevronLeft :size="16" /> Games
        </router-link>

        <div v-if="game" class="flex items-center gap-3">
          <h1 class="text-3xl font-semibold text-slate-600 tracking-wide">{{ game.title }}</h1>
          <span
            :class="game.is_finished ? 'bg-green-100 text-green-600' : 'bg-violet-100 text-violet-500'"
            class="text-xs font-medium px-3 py-1 rounded-full"
          >
            {{ game.is_finished ? 'Finished' : 'In progress' }}
          </span>
        </div>
      </div>

      <!-- Players -->
      <div class="mb-6">
        <p v-if="players.length === 0" class="text-slate-400 text-sm">No players yet.</p>
        <div v-else class="flex flex-wrap gap-2">
          <span
            v-for="player in players"
            :key="player.id"
            class="inline-flex items-center gap-1.5 text-sm bg-white text-slate-500 pl-3 pr-2 py-1 rounded-full shadow-sm"
          >
            {{ player.name }}
            <button
              v-if="game && !game.is_finished"
              @click="removePlayer(player)"
              class="text-slate-300 hover:text-red-400 transition-colors leading-none"
            >×</button>
          </span>
        </div>
      </div>

      <!-- Scoreboard table -->
      <div v-if="players.length > 0 && rounds.length > 0" class="bg-white rounded-2xl shadow-sm overflow-x-auto mb-6">
        <table class="w-full text-sm border-collapse">
          <thead>
            <tr>
              <th class="text-left px-5 py-3 text-slate-400 font-medium w-20 border-b border-r border-slate-200">Round</th>
              <th
                v-for="player in players"
                :key="player.id"
                class="px-5 py-3 text-slate-600 font-medium text-center border-b border-r border-slate-200 last:border-r-0"
              >
                <span class="inline-flex items-center justify-center gap-1">
                  <span v-if="leaderId == player.id">👑</span>
                  {{ player.name }}
                </span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="round in rounds"
              :key="round.id"
              class="border-b border-slate-200 last:border-b-0"
            >
              <td class="px-5 py-3 text-slate-400 font-medium border-r border-slate-200">{{ round.number }}</td>
              <td
                v-for="player in players"
                :key="player.id"
                class="border-r border-slate-200 last:border-r-0 text-center"
                :class="!game?.is_finished ? 'cursor-pointer hover:bg-violet-50' : ''"
                @click="startEdit(round.id, player.id)"
              >
                <input
                  v-if="editing?.roundId === round.id && editing?.playerId === player.id"
                  ref="editInput"
                  v-model="editValue"
                  type="number"
                  class="w-full px-5 py-3 text-center text-slate-600 focus:outline-none focus:bg-violet-50"
                  @keyup.enter="saveEdit"
                  @keyup.esc="cancelEdit"
                  @blur="saveEdit"
                />
                <span v-else class="block px-5 py-3 text-slate-600">
                  {{ scores[round.id]?.[player.id]?.score ?? '–' }}
                </span>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="border-t-2 border-slate-200 bg-slate-50">
              <td class="px-5 py-3 text-slate-500 font-semibold border-r border-slate-200">Total</td>
              <td
                v-for="player in players"
                :key="player.id"
                class="px-5 py-3 text-center font-semibold text-slate-600 border-r border-slate-200 last:border-r-0"
              >
                {{ totals[player.id] }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <p v-else-if="players.length > 0 && rounds.length === 0" class="text-slate-400 text-sm mb-6">
        No rounds yet.
      </p>

      <!-- Congrats banner -->
      <div v-if="game?.is_finished && leaderId" class="mb-6 bg-amber-50 border border-amber-200 rounded-2xl px-5 py-4 text-amber-700 text-sm font-medium">
        🎉 Congrats {{ players.find(p => p.id == leaderId)?.name }}!
      </div>

      <!-- Actions (only if not finished) -->
      <div v-if="game && !game.is_finished" class="flex flex-wrap gap-3 items-center">
        <button
          @click="addRound"
          class="inline-flex items-center gap-2 px-4 py-2 bg-white text-slate-500 text-sm font-medium rounded-xl shadow-sm hover:shadow-md transition-shadow"
        >
          <Plus :size="15" /> Add round
        </button>

        <div v-if="!showAddPlayer">
          <button
            @click="showAddPlayer = true"
            class="inline-flex items-center gap-2 px-4 py-2 bg-white text-slate-500 text-sm font-medium rounded-xl shadow-sm hover:shadow-md transition-shadow"
          >
            <Plus :size="15" /> Add player
          </button>
        </div>

        <div v-else class="flex items-center gap-2">
          <input
            v-model="newPlayerName"
            type="text"
            placeholder="Player name"
            @keyup.enter="addPlayer"
            @keyup.esc="showAddPlayer = false"
            class="border border-slate-200 rounded-xl px-4 py-2 text-sm text-slate-600 focus:outline-none focus:ring-2 focus:ring-violet-200"
            autofocus
          />
          <button
            @click="addPlayer"
            :disabled="addingPlayer || !newPlayerName.trim()"
            class="px-4 py-2 bg-violet-400 hover:bg-violet-500 disabled:opacity-50 text-white text-sm font-medium rounded-xl transition-colors"
          >
            Add
          </button>
          <button
            @click="showAddPlayer = false"
            class="px-3 py-2 text-sm text-slate-400 hover:text-slate-600 transition-colors"
          >
            Cancel
          </button>
        </div>

        <button
          @click="finishGame"
          class="inline-flex items-center gap-2 px-4 py-2 bg-green-100 text-green-700 hover:bg-green-200 text-sm font-medium rounded-xl transition-colors"
        >
          Finish game
        </button>
      </div>

    </div>
  </div>
</template>
