<script setup lang="ts">
import { computed } from 'vue'
import type { GameDeal } from '../types/deal'

const props = defineProps<{
  deal: GameDeal
}>()

// formatador
const formatPrice = (val: number) => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  }).format(val)
}

// data
const formattedDate = computed(() => {
  if (props.deal.status === 'free_upcoming' && props.deal.start_date) {
    const d = new Date(props.deal.start_date)
    return `A partir de ${d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' })}`
  }
  if (props.deal.end_date) {
    const d = new Date(props.deal.end_date)
    return `Até ${d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' })} às ${d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}`
  }
  return null
})

// status
const isFree = computed(() => props.deal.status === 'free_now')
const isUpcoming = computed(() => props.deal.status === 'free_upcoming')
</script>

<template>
  <div class="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-white/15 bg-white/10 backdrop-blur-md transition-all duration-300 hover:-translate-y-1 hover:border-white/30 hover:bg-white/15 hover:shadow-2xl hover:shadow-purple-500/10">
    <div>
      <!-- banner -->
      <div class="relative aspect-16/9 w-full overflow-hidden bg-black/40">
        <img
          :src="deal.image_url || '/placeholder.png'"
          :alt="deal.title"
          loading="lazy"
          class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
        />

        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>

        <!-- badge loja -->
        <div class="absolute top-3 left-3 flex items-center gap-1.5">
          <span
            v-if="deal.store === 'epic'"
            class="flex items-center gap-1 rounded-md border border-white/20 bg-zinc-950/85 px-2 py-0.5 text-[11px] font-semibold tracking-wider text-white uppercase backdrop-blur-sm shadow-sm"
          >
            <!-- epic -->
            <svg class="h-3 w-3 fill-current" viewBox="0 0 24 24">
              <path d="M4.5 2h15L22 17.5l-10 4.5-10-4.5L4.5 2zm2.2 2.7l-1.8 11.4 7.1 3.2 7.1-3.2-1.8-11.4H6.7zm5.3 2.5h3.4v2.1h-3.4v1.8h2.9v2.1H12V15h-2V7.2h2z" />
            </svg>
            Epic Games
          </span>
          <span
            v-else
            class="flex items-center gap-1 rounded-md border border-sky-400/30 bg-sky-950/85 px-2 py-0.5 text-[11px] font-semibold tracking-wider text-sky-200 uppercase backdrop-blur-sm shadow-sm"
          >
            <!-- steam -->
            <svg class="h-3 w-3 fill-current" viewBox="0 0 24 24">
              <path d="M12 2a10 10 0 0 0-10 10c0 4.7 3.2 8.6 7.6 9.7l2.8-4a3.3 3.3 0 0 1-.4-1.7 3.3 3.3 0 0 1 .4-1.6l-3.3-2.3a4.7 4.7 0 1 1 8.8-2.6 4.7 4.7 0 0 1-4.7 4.7l-2.4 3.4a3.3 3.3 0 1 1-3.7-2.6l3.5 2.5c.3-.1.6-.2 1-.2a3.3 3.3 0 0 1 3.3 3.3c0 .6-.2 1.2-.5 1.7A10 10 0 0 0 22 12 10 10 0 0 0 12 2z" />
            </svg>
            Steam
          </span>
        </div>

        <!-- badge desconto -->
        <div class="absolute top-3 right-3">
          <span
            v-if="isFree"
            class="inline-block rounded-md bg-emerald-500 px-2 py-0.5 text-xs font-black tracking-wide text-white uppercase shadow-lg shadow-emerald-500/30 animate-pulse"
          >
            100% OFF
          </span>
          <span
            v-else-if="isUpcoming"
            class="inline-block rounded-md bg-purple-500/90 px-2 py-0.5 text-xs font-bold tracking-wide text-white uppercase shadow-md"
          >
            Em Breve
          </span>
          <span
            v-else
            class="inline-block rounded-md bg-emerald-500/90 px-2 py-0.5 text-xs font-extrabold text-white shadow-md"
          >
            -{{ deal.discount_percent }}%
          </span>
        </div>

        <!-- data limite -->
        <div v-if="formattedDate" class="absolute bottom-2 left-3 right-3 text-[11px] text-white/80 font-medium truncate">
          {{ formattedDate }}
        </div>
      </div>

      <!-- info -->
      <div class="p-4 space-y-2">
        <h3 class="line-clamp-1 font-semibold text-base text-white group-hover:text-amber-300 transition-colors" :title="deal.title">
          {{ deal.title }}
        </h3>

        <!-- precos -->
        <div class="flex items-baseline gap-2">
          <template v-if="isFree">
            <span class="text-lg font-black text-emerald-400">Grátis</span>
            <span v-if="deal.original_price > 0" class="text-xs text-white/50 line-through">
              {{ formatPrice(deal.original_price) }}
            </span>
          </template>

          <template v-else-if="isUpcoming">
            <span class="text-sm font-semibold text-purple-300">Em Breve Grátis</span>
          </template>

          <template v-else>
            <span class="text-lg font-black text-white">
              {{ formatPrice(deal.final_price) }}
            </span>
            <span v-if="deal.original_price > deal.final_price" class="text-xs text-white/50 line-through">
              {{ formatPrice(deal.original_price) }}
            </span>
          </template>
        </div>
      </div>
    </div>

    <!-- botao -->
    <div class="p-4 pt-0">
      <a
        :href="deal.store_url"
        target="_blank"
        rel="noopener noreferrer"
        class="flex w-full items-center justify-center gap-1.5 rounded-xl border border-white/20 bg-white/10 py-2 px-3 text-xs font-semibold text-white transition-all duration-200 hover:scale-[1.02] hover:bg-white/25 active:scale-95 cursor-pointer shadow-sm"
      >
        <span>{{ isFree ? 'Resgatar Grátis' : isUpcoming ? 'Ver na Loja' : 'Aproveitar Oferta' }}</span>
        <svg class="h-3.5 w-3.5 opacity-80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </a>
    </div>
  </div>
</template>
