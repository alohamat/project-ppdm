<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useAuthStore } from '../auth'
import GameCard from '../components/GameCard.vue'
import type { GameDeal, DealsResponse } from '../types/deal'
import { formatCurrencyBRL } from '../utils/format'

const router = useRouter()
const authStore = useAuthStore()

const dealsData = ref<DealsResponse | null>(null)
const loading = ref(true)
const error = ref('')

const activePill = ref<'all' | 'huge_discount' | 'under_15' | 'steam_deck' | 'free_now'>('all')
const selectedStore = ref<'all' | 'epic' | 'steam'>('all')
const searchQuery = ref('')
const sortBy = ref<'discount' | 'price' | 'name'>('discount')
const isDarkMode = ref(true)

const pageSize = 16
const visibleCount = ref(pageSize)
const infiniteSentinel = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)
const steamCarouselRef = ref<HTMLElement | null>(null)
const nowTimestamp = ref(Date.now())

let observer: IntersectionObserver | null = null
let countdownInterval: number | null = null

const storeFilters = [
  { id: 'all', label: 'Todas' },
  { id: 'epic', label: 'Epic' },
  { id: 'steam', label: 'Steam' }
] as const

const pills = [
  { id: 'all', label: '⭐ Todas as Ofertas', active: 'bg-slate-100 text-slate-950 shadow-md', inactive: 'hover:border-slate-700 hover:text-white' },
  { id: 'huge_discount', label: '🔥 >75% OFF', active: 'bg-rose-600 text-white shadow-md border-rose-500', inactive: 'hover:border-rose-900/50 hover:text-rose-300' },
  { id: 'under_15', label: '🪙 Até R$ 15', active: 'bg-amber-500 text-slate-950 shadow-md border-amber-400', inactive: 'hover:border-amber-900/50 hover:text-amber-300' },
  { id: 'steam_deck', label: '🐧 Deck / Steam', active: 'bg-sky-600 text-white shadow-md border-sky-500', inactive: 'hover:border-sky-900/50 hover:text-sky-300' },
  { id: 'free_now', label: '🎁 100% Grátis', active: 'bg-emerald-600 text-white shadow-md border-emerald-500', inactive: 'hover:border-emerald-900/50 hover:text-emerald-300' }
] as const

const scrollCarousel = (dir: 'left' | 'right') => {
  steamCarouselRef.value?.scrollBy({ left: dir === 'left' ? -340 : 340, behavior: 'smooth' })
}

const fetchDeals = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get<DealsResponse>('/api/deals')
    dealsData.value = res.data
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Não foi possível carregar as ofertas. Verifique a conexão com o servidor.'
  } finally {
    loading.value = false
  }
}

const handleGlobalKeydown = (e: KeyboardEvent) => {
  if (e.key === '/' && document.activeElement?.tagName !== 'INPUT' && document.activeElement?.tagName !== 'TEXTAREA') {
    e.preventDefault()
    searchInputRef.value?.focus()
  }
}

onMounted(() => {
  fetchDeals()
  countdownInterval = window.setInterval(() => { nowTimestamp.value = Date.now() }, 1000)
  observer = new IntersectionObserver(([e]) => { if (e?.isIntersecting) loadMore() }, { rootMargin: '200px' })
  if (infiniteSentinel.value) observer.observe(infiniteSentinel.value)
  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  if (countdownInterval) clearInterval(countdownInterval)
  if (observer) observer.disconnect()
  window.removeEventListener('keydown', handleGlobalKeydown)
})

const allDeals = computed<GameDeal[]>(() => {
  if (!dealsData.value) return []
  return [...dealsData.value.free_now, ...dealsData.value.free_upcoming, ...dealsData.value.big_deals]
})

const econometer = computed(() => {
  if (!dealsData.value) return { totalSavings: 0, freebiesCount: 0, dealsCount: 0, maxDiscount: 90 }
  let totalSavings = dealsData.value.free_now.reduce((acc, d) => acc + (d.original_price > 0 ? d.original_price : 49.9), 0)
  let maxDiscount = 0
  dealsData.value.big_deals.forEach((d) => {
    if (d.original_price > d.final_price) totalSavings += (d.original_price - d.final_price)
    if (d.discount_percent > maxDiscount) maxDiscount = d.discount_percent
  })
  return {
    totalSavings,
    freebiesCount: dealsData.value.free_now.length,
    dealsCount: dealsData.value.big_deals.length,
    maxDiscount: maxDiscount || 90
  }
})

const heroFreeGame = computed<GameDeal | null>(() => {
  if (!dealsData.value?.free_now.length) return null
  return [...dealsData.value.free_now].sort((a, b) => b.original_price - a.original_price)[0]
})

const countdownUnits = computed(() => {
  if (!heroFreeGame.value?.end_date) return []
  const diff = Math.max(0, new Date(heroFreeGame.value.end_date).getTime() - nowTimestamp.value)
  const d = Math.floor(diff / 86400000)
  const h = Math.floor((diff % 86400000) / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  const s = Math.floor((diff % 60000) / 1000)

  return [
    { label: 'Dias', val: String(d).padStart(2, '0'), highlight: false },
    { label: 'Horas', val: String(h).padStart(2, '0'), highlight: false },
    { label: 'Min', val: String(m).padStart(2, '0'), highlight: false },
    { label: 'Seg', val: String(s).padStart(2, '0'), highlight: true }
  ]
})

const epicFreebiesCurrent = computed(() => dealsData.value?.free_now.filter((d) => d.store === 'epic') || [])
const epicFreebiesUpcoming = computed(() => dealsData.value?.free_upcoming.filter((d) => d.store === 'epic') || [])
const steamSpecials = computed(() => dealsData.value?.big_deals.filter((d) => d.store === 'steam') || [])

const filteredDeals = computed(() => {
  let list = [...allDeals.value]
  if (selectedStore.value !== 'all') list = list.filter((i) => i.store === selectedStore.value)

  if (activePill.value === 'huge_discount') list = list.filter((i) => i.discount_percent >= 75)
  else if (activePill.value === 'under_15') list = list.filter((i) => i.final_price > 0 && i.final_price <= 15)
  else if (activePill.value === 'steam_deck') list = list.filter((i) => i.store === 'steam')
  else if (activePill.value === 'free_now') list = list.filter((i) => i.status === 'free_now')

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter((i) => i.title.toLowerCase().includes(q))
  }

  return list.sort((a, b) => {
    if (sortBy.value === 'discount') return b.discount_percent - a.discount_percent
    if (sortBy.value === 'price') return a.final_price - b.final_price
    if (sortBy.value === 'name') return a.title.localeCompare(b.title)
    return 0
  })
})

const visibleDeals = computed(() => filteredDeals.value.slice(0, visibleCount.value))
const hasMoreDeals = computed(() => visibleCount.value < filteredDeals.value.length)
const loadMore = () => { if (hasMoreDeals.value) visibleCount.value += pageSize }

const resetAllFilters = () => {
  searchQuery.value = ''
  selectedStore.value = 'all'
  activePill.value = 'all'
  sortBy.value = 'discount'
  visibleCount.value = pageSize
}
</script>

<template>
  <div class="min-h-screen bg-slate-950 font-sf text-slate-100 flex flex-col selection:bg-emerald-500 selection:text-slate-950">
    <header class="sticky top-0 z-50 border-b border-slate-800/80 bg-slate-950/85 backdrop-blur-md">
      <div class="mx-auto flex max-w-7xl items-center justify-between gap-3 px-4 py-3 sm:px-6 lg:px-8">
        <div class="flex items-center gap-3">
          <div class="flex size-12 shrink-0 items-center justify-center rounded-xl border border-emerald-500/30 bg-emerald-500/10 text-xl shadow-inner shadow-emerald-500/20" aria-hidden="true"><img src="/miteido.png"></div>
          <div>
            <div class="flex items-center gap-2">
              <span class="text-base sm:text-lg font-black tracking-tight text-white">BadukiSales</span>
            </div>
            <p class="hidden md:block text-[11px] text-slate-400">Garimpe jogos grátis e promoções insanas</p>
          </div>
        </div>

        <div class="relative flex-1 max-w-md hidden sm:flex items-center">
          <label for="navbar-search" class="sr-only">Busca rápida</label>
          <svg class="pointer-events-none absolute left-3 size-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          <input id="navbar-search" ref="searchInputRef" v-model="searchQuery" type="search" placeholder="Buscar jogo rápido (pressione /)..." class="min-h-10 w-full rounded-xl border border-slate-800 bg-slate-900/90 py-2 pr-8 pl-9 text-xs text-slate-100 placeholder:text-slate-500 focus:border-emerald-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/30" />
          <button v-if="searchQuery" type="button" @click="searchQuery = ''" class="absolute right-2.5 p-1 text-xs text-slate-400 hover:text-white" aria-label="Limpar">✕</button>
        </div>

        <div class="flex items-center gap-2">
          <div class="hidden lg:flex items-center rounded-xl border border-slate-800 bg-slate-900/90 p-1 text-xs">
            <button v-for="s in storeFilters" :key="s.id" type="button" @click="selectedStore = s.id" class="rounded-lg px-2.5 py-1 font-semibold transition-colors" :class="selectedStore === s.id ? (s.id === 'steam' ? 'bg-sky-950/80 text-sky-300 border border-sky-800/60' : 'bg-slate-800 text-white shadow-sm') : 'text-slate-400 hover:text-slate-200'">{{ s.label }}</button>
          </div>

          <button type="button" @click="isDarkMode = !isDarkMode" class="flex size-10 items-center justify-center rounded-xl border border-slate-800 bg-slate-900/80 text-amber-300 transition-colors hover:border-slate-700 hover:bg-slate-850 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500" title="Tema Dark Ativo" aria-label="Tema escuro">
            <svg class="size-4.5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true"><path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" /></svg>
          </button>

          <template v-if="authStore.isAuthenticated">
            <button type="button" @click="router.push('/dashboard')" class="flex min-h-10 items-center gap-1.5 rounded-xl border border-slate-700 bg-slate-800 px-3.5 py-1.5 text-xs font-semibold text-slate-200 transition-colors hover:bg-slate-700 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500">
              <span>{{ authStore.isAdmin ? 'Admin' : 'Painel' }}</span>
            </button>
            <button type="button" @click="authStore.logout()" class="min-h-10 rounded-xl border border-rose-900/50 bg-rose-950/40 px-3 py-1.5 text-xs font-semibold text-rose-300 transition-colors hover:bg-rose-900/60 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-500">Sair</button>
          </template>
          <template v-else>
            <button type="button" @click="router.push('/login')" class="min-h-10 rounded-xl bg-emerald-500 px-3.5 py-1.5 text-xs font-bold text-slate-950 shadow-sm shadow-emerald-500/20 transition-colors hover:bg-emerald-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-300">Entrar</button>
          </template>
        </div>
      </div>

      <div class="sm:hidden px-4 pb-3">
        <div class="relative flex items-center">
          <svg class="pointer-events-none absolute left-3 size-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          <input v-model="searchQuery" type="search" placeholder="Buscar ofertas de jogos..." class="min-h-10 w-full rounded-xl border border-slate-800 bg-slate-900/90 py-2 pr-8 pl-9 text-xs text-slate-100 placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/30" />
        </div>
      </div>
    </header>

    <main class="mx-auto w-full max-w-7xl flex-1 space-y-10 p-4 sm:p-6 lg:p-8">
      <section aria-label="Destaque principal e Economômetro" class="space-y-4">
        <div class="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-emerald-500/20 bg-linear-to-r from-emerald-950/40 via-slate-900/80 to-slate-900/70 p-4 backdrop-blur-md shadow-lg shadow-emerald-950/20">
          <div class="flex items-center gap-3">
            <div class="flex size-10 items-center justify-center rounded-xl bg-emerald-500/20 text-emerald-400 font-black text-lg border border-emerald-500/30">🪙</div>
            <div>
              <div class="flex items-center gap-2">
                <span class="text-xs font-bold uppercase tracking-wider text-emerald-400">Economômetro Ativo</span>
                <span class="inline-flex size-2 rounded-full bg-emerald-400 animate-pulse"></span>
              </div>
              <p class="text-lg sm:text-2xl font-black text-white">
                {{ formatCurrencyBRL(econometer.totalSavings) }}
                <span class="text-xs font-normal text-slate-300">economizados em promoções ativas</span>
              </p>
            </div>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <span class="inline-flex items-center gap-1.5 rounded-xl border border-slate-800 bg-slate-950/70 px-3 py-1.5 text-xs font-semibold text-slate-300">
              <span>🎁</span><span><strong>{{ econometer.freebiesCount }}</strong> Grátis para resgatar</span>
            </span>
            <span class="inline-flex items-center gap-1.5 rounded-xl border border-slate-800 bg-slate-950/70 px-3 py-1.5 text-xs font-semibold text-amber-300">
              <span>🔥</span><span>Até <strong>-{{ econometer.maxDiscount }}%</strong> OFF</span>
            </span>
          </div>
        </div>

        <div v-if="heroFreeGame" class="relative overflow-hidden rounded-3xl border border-slate-800 bg-slate-900/80 backdrop-blur-md shadow-2xl transition-all duration-300 hover:border-slate-700">
          <div class="flex flex-col lg:flex-row">
            <div class="relative aspect-video overflow-hidden bg-slate-950">
              <img :src="heroFreeGame.image_url" :alt="heroFreeGame.title" class="size-full transition-transform duration-500 hover:scale-105" />
              <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-slate-950 via-slate-950/40 to-transparent lg:bg-linear-to-r lg:from-transparent lg:to-slate-900/95"></div>
              <div class="absolute top-4 left-4 flex items-center gap-2">
                <span class="rounded-xl border border-emerald-500/40 bg-emerald-600/90 px-3 py-1 text-xs font-black uppercase tracking-wider text-white shadow-lg backdrop-blur-md">⚡ Grátis da Semana</span>
                <span class="rounded-xl border border-slate-700/80 bg-slate-950/80 px-2.5 py-1 text-xs font-bold uppercase tracking-wider text-slate-200 backdrop-blur-md">{{ heroFreeGame.store === 'epic' ? 'Epic Games' : 'Steam' }}</span>
              </div>
            </div>

            <div class="lg:col-span-5 flex flex-col justify-between p-6 sm:p-8 space-y-6">
              <div class="space-y-3">
                <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-emerald-400">
                  <span class="inline-block size-2 rounded-full bg-emerald-400 animate-ping"></span>
                  <span>Disponível por tempo limitado</span>
                </div>
                <h2 class="text-2xl sm:text-3xl font-black text-white tracking-tight leading-tight">{{ heroFreeGame.title }}</h2>
                <p class="text-xs sm:text-sm text-slate-300">Resgate agora para adicionar permanentemente à sua biblioteca sem pagar nada.</p>
                <div class="flex items-baseline gap-3 pt-2">
                  <span class="text-3xl sm:text-4xl font-black text-emerald-400 tracking-tight">GRÁTIS</span>
                  <span v-if="heroFreeGame.original_price > 0" class="text-sm sm:text-base text-slate-400 line-through">{{ formatCurrencyBRL(heroFreeGame.original_price) }}</span>
                  <span class="rounded-lg bg-emerald-500/20 border border-emerald-500/30 px-2 py-0.5 text-xs font-extrabold text-emerald-300">100% OFF</span>
                </div>
              </div>

              <div v-if="countdownUnits.length > 0" class="rounded-2xl border border-slate-800 bg-slate-950/70 p-4 space-y-2">
                <div class="flex items-center justify-between text-xs font-semibold text-slate-300">
                  <span class="flex items-center gap-1.5 text-amber-300">
                    <svg class="size-4 animate-spin text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><circle cx="12" cy="12" r="10" stroke-width="2" stroke="currentColor" stroke-dasharray="16" /></svg>
                    <span>Tempo restante para resgate:</span>
                  </span>
                  <span class="text-slate-400 text-[11px]">Não perca!</span>
                </div>
                <div class="grid grid-cols-4 gap-2 text-center pt-1">
                  <div v-for="u in countdownUnits" :key="u.label" class="rounded-xl border border-slate-800/80 bg-slate-900/90 py-2">
                    <span class="block text-xl sm:text-2xl font-black font-mono" :class="u.highlight ? 'text-emerald-400' : 'text-white'">{{ u.val }}</span>
                    <span class="text-[10px] uppercase font-bold text-slate-400 tracking-wider">{{ u.label }}</span>
                  </div>
                </div>
              </div>

              <a :href="heroFreeGame.store_url" target="_blank" rel="noopener noreferrer" class="flex min-h-12 w-full items-center justify-center gap-2.5 rounded-2xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-sm uppercase tracking-wide transition-all shadow-xl shadow-emerald-500/25 active:scale-[0.99] focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-emerald-400/40">
                <span>Resgatar Jogo de Graça</span>
                <svg class="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M14 5l7 7m0 0l-7 7m7-7H3" /></svg>
              </a>
            </div>
          </div>
        </div>
      </section>

      <section aria-label="Radar de Brindes Epic Games" class="space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-800/80 pb-3">
          <div>
            <div class="flex items-center gap-2">
              <span class="text-xl font-black tracking-tight text-white">Radar de Brindes</span>
              <span class="rounded-full border border-violet-500/30 bg-violet-500/10 px-2 py-0.5 text-[11px] font-bold text-violet-300">Epic Games Freebies</span>
            </div>
            <p class="text-xs text-slate-400">Jogos gratuitos de hoje conectados à rotação da próxima semana</p>
          </div>
          <div class="flex items-center gap-3 text-xs text-slate-400 font-medium">
            <span class="flex items-center gap-1.5"><span class="size-2 rounded-full bg-emerald-400"></span> Resgate Ativo</span>
            <span class="flex items-center gap-1.5"><span class="size-2 rounded-full bg-violet-400"></span> Próxima Semana (P&B)</span>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-12 gap-4 items-stretch">
          <div class="md:col-span-7 space-y-3">
            <div class="flex items-center justify-between text-xs font-bold uppercase tracking-wider text-emerald-400">
              <span class="flex items-center gap-1.5"><span class="size-2 rounded-full bg-emerald-400 animate-ping"></span><span>Disponíveis para resgate agora</span></span>
              <span class="text-slate-400 font-normal">({{ epicFreebiesCurrent.length }} jogos)</span>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <GameCard v-for="deal in epicFreebiesCurrent" :key="deal.id" :deal="deal" />
            </div>
          </div>

          <div class="hidden md:flex md:col-span-1 flex-col items-center justify-center text-slate-600">
            <div class="h-10 border-l border-dashed border-slate-700"></div>
            <div class="my-2 rounded-full border border-slate-700 bg-slate-900 p-2 text-slate-400 shadow-md">
              <svg class="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" /></svg>
            </div>
            <div class="h-10 border-l border-dashed border-slate-700"></div>
          </div>

          <div class="md:col-span-4 space-y-3">
            <div class="flex items-center justify-between text-xs font-bold uppercase tracking-wider text-violet-400">
              <span>⏰ Próxima Semana (Em breve)</span>
              <span class="text-slate-400 font-normal">({{ epicFreebiesUpcoming.length }})</span>
            </div>
            <div v-if="epicFreebiesUpcoming.length > 0" class="space-y-4">
              <GameCard v-for="deal in epicFreebiesUpcoming" :key="deal.id" :deal="deal" />
            </div>
            <div v-else class="rounded-2xl border border-slate-800 bg-slate-900/40 p-8 text-center text-xs text-slate-400">
              Os próximos jogos grátis da semana serão revelados em breve pela Epic Games.
            </div>
          </div>
        </div>
      </section>

      <section aria-label="Carrossel de Ofertas Steam Specials" class="space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800/80 pb-3">
          <div>
            <div class="flex items-center gap-2">
              <h2 class="text-xl font-black tracking-tight text-white">Ofertas em Destaque</h2>
              <span class="rounded-full border border-sky-800/60 bg-sky-950/70 px-2 py-0.5 text-[11px] font-bold text-sky-300">Steam Specials</span>
            </div>
            <p class="text-xs text-slate-400">Navegue pelas melhores promoções ativas na Steam hoje</p>
          </div>
          <div class="flex items-center gap-2">
            <button type="button" @click="scrollCarousel('left')" class="flex size-9 items-center justify-center rounded-xl border border-slate-800 bg-slate-900/90 text-slate-300 hover:border-slate-700 hover:bg-slate-800 hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500" aria-label="Rolar esquerda">◄</button>
            <button type="button" @click="scrollCarousel('right')" class="flex size-9 items-center justify-center rounded-xl border border-slate-800 bg-slate-900/90 text-slate-300 hover:border-slate-700 hover:bg-slate-800 hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500" aria-label="Rolar direita">►</button>
          </div>
        </div>

        <div ref="steamCarouselRef" class="no-scrollbar flex gap-4 overflow-x-auto scroll-smooth snap-x pb-2 pt-1">
          <div v-for="deal in steamSpecials" :key="deal.id" class="w-72 sm:w-80 shrink-0 snap-start">
            <GameCard :deal="deal" />
          </div>
        </div>
      </section>

      <section aria-label="Filtros pílula" class="space-y-4 rounded-2xl border border-slate-800 bg-slate-900/60 p-4 sm:p-5 backdrop-blur-md shadow-sm">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div class="flex flex-wrap items-center gap-2">
            <button v-for="p in pills" :key="p.id" type="button" @click="activePill = p.id; visibleCount = pageSize" class="rounded-full px-4 py-2 text-xs font-bold transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 border border-slate-800 bg-slate-900 text-slate-300" :class="activePill === p.id ? p.active : p.inactive">{{ p.label }}</button>
          </div>

          <div class="flex items-center gap-2">
            <label for="sort-select" class="text-xs text-slate-400 whitespace-nowrap">Ordenar por:</label>
            <select id="sort-select" v-model="sortBy" class="rounded-xl border border-slate-700 bg-slate-950 px-3 py-1.5 text-xs text-slate-200 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500">
              <option value="discount">Maior Desconto (%)</option>
              <option value="price">Menor Preço (R$)</option>
              <option value="name">Ordem Alfabética (A-Z)</option>
            </select>
          </div>
        </div>
      </section>

      <section aria-label="Feed de ofertas" class="space-y-6">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Exibindo <strong>{{ visibleDeals.length }}</strong> de <strong>{{ filteredDeals.length }}</strong> ofertas</span>
          <button v-if="searchQuery || activePill !== 'all' || selectedStore !== 'all'" type="button" @click="resetAllFilters" class="text-emerald-400 hover:underline">Limpar filtros ativos</button>
        </div>

        <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4" aria-busy="true">
          <div v-for="n in 8" :key="n" class="animate-pulse space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
            <div class="aspect-video w-full rounded-xl bg-slate-800"></div>
            <div class="h-4 w-3/4 rounded bg-slate-800"></div>
            <div class="h-3 w-1/2 rounded bg-slate-800/70"></div>
            <div class="h-10 w-full rounded-xl bg-slate-800 pt-2"></div>
          </div>
        </div>

        <div v-else-if="error" class="space-y-3 rounded-2xl border border-rose-900/60 bg-rose-950/30 p-8 text-center">
          <span class="text-3xl" aria-hidden="true">⚠️</span>
          <h2 class="text-base font-bold text-rose-200">Não foi possível carregar o feed</h2>
          <p class="mx-auto max-w-md text-xs text-rose-300 sm:text-sm">{{ error }}</p>
          <div>
            <button type="button" @click="fetchDeals" class="min-h-11 rounded-xl bg-rose-600 px-5 py-2 text-xs font-bold text-white shadow hover:bg-rose-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-400">Tentar novamente</button>
          </div>
        </div>

        <div v-else-if="visibleDeals.length > 0" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <GameCard v-for="deal in visibleDeals" :key="deal.id" :deal="deal" />
        </div>

        <div v-else class="space-y-3 rounded-2xl border border-slate-800 bg-slate-900/60 p-12 text-center backdrop-blur-md">
          <span class="text-4xl" aria-hidden="true">🎮</span>
          <h3 class="text-base font-bold text-slate-100">Nenhum jogo encontrado para estes filtros</h3>
          <p class="mx-auto max-w-md text-xs text-slate-400 sm:text-sm">Tente ajustar o termo de busca ou selecionar outro filtro pílula.</p>
          <div class="pt-2">
            <button type="button" @click="resetAllFilters" class="min-h-10 rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500">Redefinir filtros</button>
          </div>
        </div>

        <div class="flex flex-col items-center justify-center pt-6 space-y-3">
          <div ref="infiniteSentinel" class="h-2 w-full"></div>
          <button v-if="hasMoreDeals" type="button" @click="loadMore" class="min-h-11 rounded-xl border border-slate-700 bg-slate-900/90 hover:bg-slate-800 px-6 py-2.5 text-xs font-bold text-slate-200 transition-colors shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500">
            Carregar mais ofertas ({{ filteredDeals.length - visibleDeals.length }} restantes)
          </button>
          <span v-else-if="filteredDeals.length > 0" class="text-xs text-slate-500">Você visualizou todas as ofertas disponíveis hoje! 🎉</span>
        </div>
      </section>
    </main>

    <footer class="border-t border-slate-800/80 bg-slate-950 py-8 text-center text-xs text-slate-400">
      <div class="mx-auto max-w-7xl px-4 space-y-2">
        <div class="flex items-center justify-center gap-2">
          <span class="text-sm">🎮</span>
          <span class="font-bold text-slate-200">BadukiSales</span>
          <span class="text-slate-600">•</span>
          <span>Radar de Ofertas Gamer</span>
        </div>
        <p class="text-slate-400">Dados e preços sincronizados via APIs públicas oficiais da Epic Games Store e Steam. Todos os direitos reservados.</p>
      </div>
    </footer>
  </div>
</template>
