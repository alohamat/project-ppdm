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

const storeBadgeConfig = computed(() => {
  const isEpic = props.deal.store === 'epic'
  return {
    label: isEpic ? 'Epic Games' : 'Steam',
    variantClasses: isEpic
      ? 'border-slate-700 bg-slate-950/90 text-slate-200'
      : 'border-sky-800/80 bg-sky-950/90 text-sky-200',
    isEpic
  }
})

const discountBadge = computed(() => {
  if (isFree.value) {
    return {
      text: '100% OFF',
      classes: 'bg-emerald-600 uppercase font-bold'
    }
  }
  if (isUpcoming.value) {
    return {
      text: 'Em breve',
      classes: 'bg-violet-600 uppercase font-semibold'
    }
  }
  return {
    text: `-${props.deal.discount_percent}%`,
    classes: 'bg-emerald-600 font-bold'
  }
})

const hasStrikethroughPrice = computed(() => {
  if (isUpcoming.value) return false
  if (isFree.value) return props.deal.original_price > 0
  return props.deal.original_price > props.deal.final_price
})

const handleImageError = () => {
  imageError.value = true
}
</script>

<template>
  <article
    class="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-slate-800 bg-slate-900 transition duration-200 hover:-translate-y-0.5 hover:border-slate-700 hover:shadow-xl hover:shadow-black/50"
  >
    <!-- Visual / Media Header -->
    <div class="relative aspect-video w-full overflow-hidden bg-slate-950">
      <img
        v-if="!imageError && deal.image_url"
        :src="deal.image_url"
        :alt="deal.title"
        loading="lazy"
        width="460"
        height="215"
        @error="handleImageError"
        class="size-full object-cover transition-transform duration-300 group-hover:scale-105"
      />
      <div
        v-else
        class="flex size-full items-center justify-center text-xs font-medium text-slate-500"
      >
        <span>Imagem indisponível</span>
      </div>

      <!-- Scrim gradient for contrast -->
      <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-slate-950/90 via-slate-950/30 to-transparent"></div>

      <!-- Store badge -->
      <div class="absolute top-3 left-3 flex items-center">
        <span
          class="flex items-center gap-1 rounded-md border px-2 py-0.5 text-xs font-semibold uppercase shadow-sm"
          :class="storeBadgeConfig.variantClasses"
        >
          <svg class="size-3 fill-current" viewBox="0 0 24 24" aria-hidden="true">
            <path
              v-if="storeBadgeConfig.isEpic"
              d="M4.5 2h15L22 17.5l-10 4.5-10-4.5L4.5 2zm2.2 2.7l-1.8 11.4 7.1 3.2 7.1-3.2-1.8-11.4H6.7zm5.3 2.5h3.4v2.1h-3.4v1.8h2.9v2.1H12V15h-2V7.2h2z"
            />
            <path
              v-else
              d="M12 2a10 10 0 0 0-10 10c0 4.7 3.2 8.6 7.6 9.7l2.8-4a3.3 3.3 0 0 1-.4-1.7 3.3 3.3 0 0 1 .4-1.6l-3.3-2.3a4.7 4.7 0 1 1 8.8-2.6 4.7 4.7 0 0 1-4.7 4.7l-2.4 3.4a3.3 3.3 0 1 1-3.7-2.6l3.5 2.5c.3-.1.6-.2 1-.2a3.3 3.3 0 0 1 3.3 3.3c0 .6-.2 1.2-.5 1.7A10 10 0 0 0 22 12 10 10 0 0 0 12 2z"
            />
          </svg>
          {{ storeBadgeConfig.label }}
        </span>
      </div>

      <!-- Discount badge -->
      <div class="absolute top-3 right-3">
        <span
          class="inline-block rounded-md px-2 py-0.5 text-xs text-white shadow-md"
          :class="discountBadge.classes"
        >
          {{ discountBadge.text }}
        </span>
      </div>

      <!-- Expiration / Availability info -->
      <div v-if="formattedDate" class="absolute bottom-2 inset-x-3 truncate text-xs font-medium text-slate-300">
        {{ formattedDate }}
      </div>
    </div>

    <!-- Card Body & Action -->
    <div class="flex flex-1 flex-col justify-between gap-4 p-4">
      <div class="space-y-2">
        <h3
          class="line-clamp-1 text-base font-semibold text-slate-100 transition-colors group-hover:text-emerald-400"
          :title="deal.title"
        >
          {{ deal.title }}
        </h3>

        <div class="flex items-baseline gap-2">
          <template v-if="isUpcoming">
            <span class="text-sm font-semibold text-violet-300">Em breve na faixa</span>
          </template>

          <template v-else>
            <span
              class="text-lg font-black"
              :class="isFree ? 'text-emerald-400' : 'text-slate-100'"
            >
              {{ isFree ? 'Grátis' : formatPrice(deal.final_price) }}
            </span>
            <span
              v-if="hasStrikethroughPrice"
              class="text-xs text-slate-400 line-through"
            >
              {{ formatPrice(deal.original_price) }}
            </span>
          </template>
        </div>
      </div>

      <a
        :href="deal.store_url"
        target="_blank"
        rel="noopener noreferrer"
        class="flex min-h-11 w-full items-center justify-center gap-2 rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-2.5 text-xs font-semibold text-slate-200 transition-colors duration-200 hover:border-emerald-500 hover:bg-slate-800 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 active:scale-[0.98]"
      >
        <span>{{ actionText }}</span>
        <svg class="size-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </a>
    </div>
  </article>
</template>
