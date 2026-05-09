<script setup lang="ts">
import type { Game } from '../../types/api.ts'

const props = defineProps<{ game: Game }>()

const date = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
}).format(new Date(props.game.created_at))
</script>

<template>
  <div class="bg-white rounded-2xl px-5 py-4 shadow-sm flex items-center justify-between gap-4 cursor-pointer hover:shadow-md transition-shadow">
    <div>
      <span class="font-medium text-slate-600">{{ game.title }}</span>
      <p class="text-xs text-slate-400 mt-0.5">{{ date }}</p>
    </div>
    <span
      :class="game.is_finished ? 'bg-green-100 text-green-600' : 'bg-violet-100 text-violet-500'"
      class="text-xs font-medium px-3 py-1 rounded-full shrink-0"
    >
      {{ game.is_finished ? 'Finished' : 'In progress' }}
    </span>
  </div>
</template>
