<script setup lang="ts">
import { computed, ref } from 'vue'
import type { GameDeal } from '../types/deal'
import { formatCurrencyBRL, getRemainingTime, formatUpcomingDate } from '../utils/format'

const props = defineProps<{
  deal: GameDeal
  compact?: boolean
}>()

const imageError = ref(false)

const expirationInfo = computed(() => {
  if (props.deal.status === 'free_now' || props.deal.status === 'deal') {
    return getRemainingTime(props.deal.end_date)
  }
  return null
})

const upcomingInfo = computed(() => {
  if (props.deal.status === 'free_upcoming') {
    return formatUpcomingDate(props.deal.start_date)
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
      ? 'border-slate-300 bg-white/90 text-slate-800 dark:border-slate-700 dark:bg-slate-950/90 dark:text-slate-200'
      : 'border-sky-300 bg-sky-50 text-sky-800 dark:border-sky-800/80 dark:bg-sky-950/90 dark:text-sky-200',
    isEpic
  }
})

const discountBadge = computed(() => {
  if (isFree.value) {
    return {
      text: '100% OFF',
      classes: 'bg-emerald-600 uppercase font-bold text-white shadow-emerald-500/20'
    }
  }
  if (isUpcoming.value) {
    return {
      text: 'Em breve',
      classes: 'bg-violet-600 uppercase font-semibold text-white shadow-violet-500/20'
    }
  }
  return {
    text: `-${props.deal.discount_percent}%`,
    classes: 'bg-emerald-600 font-bold text-white shadow-emerald-500/20'
  }
})

const hasStrikethroughPrice = computed(() => {
  if (isUpcoming.value) return false
  if (isFree.value) return props.deal.original_price > 0
  return props.deal.original_price > props.deal.final_price
})
</script>

<template>
  <article
    class="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-slate-200/90 bg-white/90 dark:border-slate-800/80 dark:bg-slate-900/70 backdrop-blur-sm transition-all duration-200 hover:-translate-y-1 hover:border-slate-300 dark:hover:border-slate-700 hover:shadow-xl hover:shadow-slate-300/40 dark:hover:shadow-black/40"
    :class="{ 'opacity-85 grayscale hover:grayscale-0 hover:opacity-100': isUpcoming }"
  >
    <div class="relative aspect-video w-full overflow-hidden bg-slate-950">
      <img
        v-if="!imageError && deal.image_url"
        :src="deal.image_url"
        :alt="deal.title"
        loading="lazy"
        width="460"
        height="215"
        @error="imageError = true"
        class="size-full object-cover transition-transform duration-300 group-hover:scale-105"
      />
      <div v-else class="flex size-full items-center justify-center text-xs font-medium text-slate-500">
        <span>Imagem indisponível</span>
      </div>

      <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-slate-950/95 via-slate-950/30 to-transparent"></div>

      <div class="absolute top-2.5 left-2.5 flex items-center">
        <span
          class="flex items-center gap-1 rounded-md border px-2 py-0.5 text-[11px] font-semibold uppercase shadow-sm backdrop-blur-md"
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

      <div class="absolute top-2.5 right-2.5">
        <span class="inline-block rounded-md px-2 py-0.5 text-xs shadow-md" :class="discountBadge.classes">
          {{ discountBadge.text }}
        </span>
      </div>

      <div class="absolute bottom-2 inset-x-2.5 flex items-center justify-between text-xs">
        <div v-if="expirationInfo" class="flex items-center gap-1 font-semibold truncate" :class="expirationInfo.isUrgent ? 'text-amber-300' : 'text-slate-300'">
          <span class="inline-block size-1.5 rounded-full" :class="expirationInfo.isUrgent ? 'bg-amber-400 animate-ping' : 'bg-emerald-400'"></span>
          <span>{{ expirationInfo.text }}</span>
        </div>
        <div v-else-if="upcomingInfo" class="text-violet-300 font-medium truncate">
          {{ upcomingInfo }}
        </div>
      </div>
    </div>

    <div class="flex flex-1 flex-col justify-between gap-3.5 p-3.5 sm:p-4">
      <div class="space-y-1.5">
        <h3
          class="line-clamp-1 text-sm sm:text-base font-semibold text-slate-800 dark:text-slate-100 transition-colors group-hover:text-emerald-600 dark:group-hover:text-emerald-400"
          :title="deal.title"
        >
          {{ deal.title }}
        </h3>

        <div class="flex items-baseline gap-2">
          <template v-if="isUpcoming">
            <span class="text-xs sm:text-sm font-semibold text-violet-600 dark:text-violet-300">Em breve de graça</span>
          </template>

          <template v-else>
            <span
              class="text-base sm:text-lg font-black"
              :class="isFree ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-900 dark:text-slate-100'"
            >
              {{ isFree ? 'Grátis' : formatCurrencyBRL(deal.final_price) }}
            </span>
            <span
              v-if="hasStrikethroughPrice"
              class="text-xs text-slate-400 dark:text-slate-500 line-through"
              :title="`Preço original: ${formatCurrencyBRL(deal.original_price)}`"
            >
              {{ formatCurrencyBRL(deal.original_price) }}
            </span>
          </template>
        </div>
      </div>

      <a
        :href="deal.store_url"
        target="_blank"
        rel="noopener noreferrer"
        class="flex min-h-10 w-full items-center justify-center gap-2 rounded-xl border border-slate-200 bg-slate-100/90 text-slate-700 hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-700 dark:border-slate-700/80 dark:bg-slate-800/90 dark:text-slate-200 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-300 px-3.5 py-2 text-xs font-semibold transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 active:scale-[0.98]"
      >
        <span>{{ actionText }}</span>
        <svg class="size-3.5 text-slate-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </a>
    </div>
  </article>
</template>
