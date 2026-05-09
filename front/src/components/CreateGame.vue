<script setup>
import { ref } from 'vue'
import { Plus } from 'lucide-vue-next'

const API = import.meta.env.VITE_API_URL ?? '/api'

const emit = defineEmits(['created'])

const showModal = ref(false)
const newTitle = ref('')
const creating = ref(false)

function openModal() {
  newTitle.value = ''
  showModal.value = true
}

async function createGame() {
  if (!newTitle.value.trim()) return
  creating.value = true
  try {
    const res = await fetch(`${API}/games`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTitle.value.trim(), is_finished: false }),
    })
    const game = await res.json()
    emit('created', game)
    showModal.value = false
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <button
    @click="openModal"
    class="fixed bottom-6 right-6 bg-violet-400 hover:bg-violet-500 text-white rounded-full w-14 h-14 flex items-center justify-center shadow-lg transition-colors"
  >
    <Plus :size="24" />
  </button>

  <Transition name="fade">
    <div
      v-if="showModal"
      class="fixed inset-0 bg-black/30 flex items-end sm:items-center justify-center px-4 pb-4 sm:pb-0 z-50"
      @click.self="showModal = false"
    >
      <div class="bg-white rounded-2xl w-full max-w-md p-6 shadow-xl">
        <h2 class="text-lg font-semibold text-slate-600 mb-5">New game</h2>

        <div class="mb-4">
          <label class="block text-sm text-slate-500 mb-1">Title</label>
          <input
            v-model="newTitle"
            type="text"
            placeholder="Saturday Night"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm text-slate-600 focus:outline-none focus:ring-2 focus:ring-violet-200"
          />
        </div>

        <div class="flex gap-3 justify-end">
          <button
            @click="showModal = false"
            class="px-4 py-2 text-sm text-slate-400 hover:text-slate-600 transition-colors"
          >
            Cancel
          </button>
          <button
            @click="createGame"
            :disabled="creating || !newTitle.trim()"
            class="px-5 py-2 bg-violet-400 hover:bg-violet-500 disabled:opacity-50 text-white text-sm font-medium rounded-xl transition-colors"
          >
            {{ creating ? 'Creating…' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.15s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
