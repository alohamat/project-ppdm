<script setup lang="ts">
import { computed, ref } from 'vue'
import type { GameDeal } from '../types/deal'

const props = defineProps<{
  deal: GameDeal
}>()

const imageError = ref(false)

const formatPrice = (val: number) => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  }).format(val)
}

const formattedDate = computed(() => {
  if (props.deal.status === 'free_upcoming' && props.deal.start_date) {
    const d = new Date(props.deal.start_date)
    return `Libera em ${d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' })}`
  }
  if (props.deal.end_date) {
    const d = new Date(props.deal.end_date)
    return `Até ${d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' })} às ${d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}`
  }
  return null
})

const isFree = computed(() => props.deal.status === 'free_now')
const isUpcoming = computed(() => props.deal.status === 'free_upcoming')

const actionText = computed(() => {
  if (isFree.value) return 'Pegar de graça'
  if (isUpcoming.value) return 'Ver na loja'
  return 'Ir pra promoção'
})

const handleImageError = () => {
  imageError.value = true
}
</script>

<template>
  <article
    class="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-slate-800 bg-slate-900 transition-all duration-200 hover:-translate-y-0.5 hover:border-slate-700 hover:shadow-xl hover:shadow-black/50"
  >
    <div>
      <div class="relative aspect-video w-full overflow-hidden bg-slate-950">
        <img
          v-if="!imageError && deal.image_url"
          :src="deal.image_url"
          :alt="deal.title"
          loading="lazy"
          width="460"
          height="215"
          @error="handleImageError"
          class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />
        <div
          v-else
          class="flex h-full w-full items-center justify-center bg-slate-950 text-slate-500 text-xs font-medium"
        >
          <span>Imagem indisponível</span>
        </div>

        <div class="absolute inset-0 bg-linear-to-t from-slate-950/90 via-slate-950/30 to-transparent pointer-events-none"></div>

        <div class="absolute top-3 left-3 flex items-center gap-1.5">
          <span
            v-if="deal.store === 'epic'"
            class="flex items-center gap-1 rounded-md border border-slate-700 bg-slate-950/90 px-2 py-0.5 text-[11px] font-semibold text-slate-200 uppercase shadow-sm"
          >
            <svg class="h-3 w-3 fill-current" viewBox="0 0 24 24" aria-hidden="true">
              <path d="M4.5 2h15L22 17.5l-10 4.5-10-4.5L4.5 2zm2.2 2.7l-1.8 11.4 7.1 3.2 7.1-3.2-1.8-11.4H6.7zm5.3 2.5h3.4v2.1h-3.4v1.8h2.9v2.1H12V15h-2V7.2h2z" />
            </svg>
            Epic Games
          </span>
          <span
            v-else
            class="flex items-center gap-1 rounded-md border border-sky-800/80 bg-sky-950/90 px-2 py-0.5 text-[11px] font-semibold text-sky-200 uppercase shadow-sm"
          >
            <svg class="h-3 w-3 fill-current" viewBox="0 0 24 24" aria-hidden="true">
              <path d="M12 2a10 10 0 0 0-10 10c0 4.7 3.2 8.6 7.6 9.7l2.8-4a3.3 3.3 0 0 1-.4-1.7 3.3 3.3 0 0 1 .4-1.6l-3.3-2.3a4.7 4.7 0 1 1 8.8-2.6 4.7 4.7 0 0 1-4.7 4.7l-2.4 3.4a3.3 3.3 0 1 1-3.7-2.6l3.5 2.5c.3-.1.6-.2 1-.2a3.3 3.3 0 0 1 3.3 3.3c0 .6-.2 1.2-.5 1.7A10 10 0 0 0 22 12 10 10 0 0 0 12 2z" />
            </svg>
            Steam
          </span>
        </div>

        <div class="absolute top-3 right-3">
          <span
            v-if="isFree"
            class="inline-block rounded-md bg-emerald-600 px-2 py-0.5 text-xs font-bold text-white uppercase shadow-md"
          >
            100% OFF
          </span>
          <span
            v-else-if="isUpcoming"
            class="inline-block rounded-md bg-violet-600 px-2 py-0.5 text-xs font-semibold text-white uppercase shadow-md"
          >
            Em breve
          </span>
          <span
            v-else
            class="inline-block rounded-md bg-emerald-600 px-2 py-0.5 text-xs font-bold text-white shadow-md"
          >
            -{{ deal.discount_percent }}%
          </span>
        </div>

        <div v-if="formattedDate" class="absolute bottom-2 left-3 right-3 text-[11px] text-slate-300 font-medium truncate">
          {{ formattedDate }}
        </div>
      </div>

      <div class="p-4 space-y-2">
        <h3
          class="line-clamp-1 font-semibold text-base text-slate-100 group-hover:text-emerald-400 transition-colors"
          :title="deal.title"
        >
          {{ deal.title }}
        </h3>

        <div class="flex items-baseline gap-2">
          <template v-if="isFree">
            <span class="text-lg font-black text-emerald-400">Grátis</span>
            <span v-if="deal.original_price > 0" class="text-xs text-slate-400 line-through">
              {{ formatPrice(deal.original_price) }}
            </span>
          </template>

          <template v-else-if="isUpcoming">
            <span class="text-sm font-semibold text-violet-300">Em breve na faixa</span>
          </template>

          <template v-else>
            <span class="text-lg font-black text-slate-100">
              {{ formatPrice(deal.final_price) }}
            </span>
            <span v-if="deal.original_price > deal.final_price" class="text-xs text-slate-400 line-through">
              {{ formatPrice(deal.original_price) }}
            </span>
          </template>
        </div>
      </div>
    </div>

    <div class="p-4 pt-0">
      <a
        :href="deal.store_url"
        target="_blank"
        rel="noopener noreferrer"
        class="flex min-h-[44px] w-full items-center justify-center gap-2 rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-2.5 text-xs font-semibold text-slate-200 transition-colors duration-200 hover:border-emerald-500 hover:bg-slate-800 hover:text-white focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 focus-visible:outline-none active:scale-[0.98]"
      >
        <span>{{ actionText }}</span>
        <svg class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </a>
    </div>
  </article>
</template>
