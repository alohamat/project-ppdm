<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useAuthStore } from '../auth'
import GameCard from '../components/GameCard.vue'
import type { GameDeal, DealsResponse } from '../types/deal'

const router = useRouter()
const authStore = useAuthStore()

const dealsData = ref<DealsResponse | null>(null)
const loading = ref(true)
const error = ref('')

const activeTab = ref<'all' | 'free_now' | 'free_upcoming' | 'big_deals'>('all')
const selectedStore = ref<'all' | 'epic' | 'steam'>('all')
const searchQuery = ref('')
const sortBy = ref<'discount' | 'price' | 'name'>('discount')

const fetchDeals = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get<DealsResponse>('/api/deals')
    dealsData.value = res.data
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Não deu pra carregar as ofertas agora. Dá uma checada na conexão.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDeals()
})

const allDeals = computed<GameDeal[]>(() => {
  if (!dealsData.value) return []
  return [
    ...dealsData.value.free_now,
    ...dealsData.value.free_upcoming,
    ...dealsData.value.big_deals
  ]
})

const counts = computed(() => {
  return {
    all: allDeals.value.length,
    free_now: dealsData.value?.free_now.length || 0,
    free_upcoming: dealsData.value?.free_upcoming.length || 0,
    big_deals: dealsData.value?.big_deals.length || 0
  }
})

const resetFilters = () => {
  searchQuery.value = ''
  selectedStore.value = 'all'
  activeTab.value = 'all'
  sortBy.value = 'discount'
}

const filteredDeals = computed(() => {
  let list: GameDeal[] = []

  if (activeTab.value === 'all') {
    list = [...allDeals.value]
  } else if (activeTab.value === 'free_now') {
    list = dealsData.value?.free_now || []
  } else if (activeTab.value === 'free_upcoming') {
    list = dealsData.value?.free_upcoming || []
  } else if (activeTab.value === 'big_deals') {
    list = dealsData.value?.big_deals || []
  }

  if (selectedStore.value !== 'all') {
    list = list.filter((item) => item.store === selectedStore.value)
  }

  if (searchQuery.value.trim() !== '') {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter((item) => item.title.toLowerCase().includes(q))
  }

  return list.sort((a, b) => {
    if (sortBy.value === 'discount') {
      return b.discount_percent - a.discount_percent
    }
    if (sortBy.value === 'price') {
      return a.final_price - b.final_price
    }
    if (sortBy.value === 'name') {
      return a.title.localeCompare(b.title)
    }
    return 0
  })
})

const handleLogout = () => {
  authStore.logout()
}
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 font-sf flex flex-col">
    <div class="max-w-7xl mx-auto w-full p-4 sm:p-6 lg:p-8 space-y-6 flex-1">

      <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-2xl border border-slate-800 bg-slate-900/90 p-4 sm:p-6 shadow-lg">
        <div class="flex items-center gap-3.5">
          <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-2xl shadow-inner" aria-hidden="true">
            🎮
          </div>
          <div>
            <h1 class="text-xl sm:text-2xl font-black tracking-tight text-white">BadukiSales</h1>
            <p class="text-xs sm:text-sm text-slate-400">Garimpe jogos grátis e os maiores descontos da Epic Games e Steam</p>
          </div>
        </div>

        <nav aria-label="Navegação rápida" class="flex flex-wrap items-center gap-2.5">
          <button
            type="button"
            @click="fetchDeals"
            :disabled="loading"
            class="flex min-h-11 items-center gap-2 rounded-xl border border-slate-700 bg-slate-800 px-3.5 py-2 text-xs font-semibold text-slate-200 transition-colors hover:border-slate-600 hover:bg-slate-750 hover:text-white focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
            aria-label="Atualizar lista de jogos"
          >
            <svg class="h-4 w-4 text-emerald-400" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span class="inline">Atualizar</span>
          </button>

          <template v-if="authStore.isAuthenticated">
            <button
              type="button"
              @click="router.push('/dashboard')"
              class="flex min-h-11 items-center gap-1.5 rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors hover:border-slate-600 hover:bg-slate-750 hover:text-white focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
            >
              <span>{{ authStore.isAdmin ? 'Painel do Admin' : 'Meu Espaço' }}</span>
            </button>
            <button
              type="button"
              @click="handleLogout"
              class="min-h-11 rounded-xl border border-rose-900/50 bg-rose-950/40 px-3.5 py-2 text-xs font-semibold text-rose-300 transition-colors hover:bg-rose-900/60 hover:text-white focus-visible:ring-2 focus-visible:ring-rose-500 focus-visible:outline-none"
            >
              Sair
            </button>
          </template>

          <template v-else>
            <button
              type="button"
              @click="router.push('/login')"
              class="min-h-11 rounded-xl bg-emerald-500 px-4 py-2 text-xs font-bold text-slate-950 transition-colors hover:bg-emerald-400 focus-visible:ring-2 focus-visible:ring-emerald-300 focus-visible:outline-none shadow-sm shadow-emerald-500/20"
            >
              Entrar ou Criar conta
            </button>
          </template>
        </nav>
      </header>

      <main class="space-y-6">
        <section aria-label="Filtros e busca" class="space-y-4 rounded-2xl border border-slate-800 bg-slate-900/80 p-4 sm:p-5 shadow-sm">
          <div role="tablist" aria-label="Categorias de ofertas" class="flex flex-wrap items-center gap-2 border-b border-slate-800 pb-4">
            <button
              type="button"
              role="tab"
              :aria-selected="activeTab === 'all'"
              @click="activeTab = 'all'"
              class="flex min-h-10 items-center gap-2 rounded-xl px-4 py-2 text-xs font-semibold transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
              :class="activeTab === 'all' ? 'bg-slate-100 text-slate-950 shadow-sm' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <span>Tudo</span>
              <span
                class="rounded-full px-2 py-0.5 text-[11px] font-bold"
                :class="activeTab === 'all' ? 'bg-slate-300 text-slate-900' : 'bg-slate-700 text-slate-200'"
              >
                {{ counts.all }}
              </span>
            </button>

            <button
              type="button"
              role="tab"
              :aria-selected="activeTab === 'free_now'"
              @click="activeTab = 'free_now'"
              class="flex min-h-10 items-center gap-2 rounded-xl px-4 py-2 text-xs font-semibold transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
              :class="activeTab === 'free_now' ? 'bg-emerald-600 text-white shadow-sm' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <span>De graça agora</span>
              <span
                class="rounded-full px-2 py-0.5 text-[11px] font-bold"
                :class="activeTab === 'free_now' ? 'bg-emerald-800 text-white' : 'bg-slate-700 text-slate-200'"
              >
                {{ counts.free_now }}
              </span>
            </button>

            <button
              type="button"
              role="tab"
              :aria-selected="activeTab === 'free_upcoming'"
              @click="activeTab = 'free_upcoming'"
              class="flex min-h-10 items-center gap-2 rounded-xl px-4 py-2 text-xs font-semibold transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
              :class="activeTab === 'free_upcoming' ? 'bg-violet-600 text-white shadow-sm' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <span>De graça em breve</span>
              <span
                class="rounded-full px-2 py-0.5 text-[11px] font-bold"
                :class="activeTab === 'free_upcoming' ? 'bg-violet-800 text-white' : 'bg-slate-700 text-slate-200'"
              >
                {{ counts.free_upcoming }}
              </span>
            </button>

            <button
              type="button"
              role="tab"
              :aria-selected="activeTab === 'big_deals'"
              @click="activeTab = 'big_deals'"
              class="flex min-h-10 items-center gap-2 rounded-xl px-4 py-2 text-xs font-semibold transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
              :class="activeTab === 'big_deals' ? 'bg-amber-500 text-slate-950 shadow-sm' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <span>Super descontos</span>
              <span
                class="rounded-full px-2 py-0.5 text-[11px] font-bold"
                :class="activeTab === 'big_deals' ? 'bg-amber-700 text-white' : 'bg-slate-700 text-slate-200'"
              >
                {{ counts.big_deals }}
              </span>
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-12 gap-3 pt-1">
            <div class="lg:col-span-6 relative flex items-center">
              <label for="search-deals" class="sr-only">Buscar pelo nome do jogo</label>
              <svg class="absolute left-3.5 h-4 w-4 text-slate-400 pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input
                id="search-deals"
                v-model="searchQuery"
                type="search"
                placeholder="Procurar jogo por título..."
                class="min-h-11 w-full rounded-xl border border-slate-700 bg-slate-950 py-2.5 pl-10 pr-9 text-xs sm:text-sm text-slate-100 transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
              />
              <button
                v-if="searchQuery"
                type="button"
                @click="searchQuery = ''"
                class="absolute right-3 text-slate-400 hover:text-white text-xs p-1"
                aria-label="Limpar campo de busca"
              >
                ✕
              </button>
            </div>

            <div class="lg:col-span-3">
              <label for="select-store" class="sr-only">Filtrar por loja</label>
              <select
                id="select-store"
                v-model="selectedStore"
                class="min-h-11 w-full rounded-xl border border-slate-700 bg-slate-950 py-2 px-3 text-xs sm:text-sm text-slate-200 transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
              >
                <option value="all">Todas as Lojas</option>
                <option value="epic">Epic Games</option>
                <option value="steam">Steam</option>
              </select>
            </div>

            <div class="lg:col-span-3">
              <label for="select-sort" class="sr-only">Ordenar resultados</label>
              <select
                id="select-sort"
                v-model="sortBy"
                class="min-h-11 w-full rounded-xl border border-slate-700 bg-slate-950 py-2 px-3 text-xs sm:text-sm text-slate-200 transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
              >
                <option value="discount">Maior Desconto</option>
                <option value="price">Menor Preço</option>
                <option value="name">Ordem Alfabética (A a Z)</option>
              </select>
            </div>
          </div>
        </section>

        <section aria-label="Catálogo de jogos">
          <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5" aria-busy="true">
            <div v-for="n in 8" :key="n" class="rounded-2xl border border-slate-800 bg-slate-900 p-4 space-y-3">
              <div class="aspect-video w-full rounded-xl bg-slate-800 animate-pulse"></div>
              <div class="h-4 w-3/4 rounded bg-slate-800 animate-pulse"></div>
              <div class="h-3 w-1/2 rounded bg-slate-800/70 animate-pulse"></div>
              <div class="h-10 w-full rounded-xl bg-slate-800 animate-pulse pt-2"></div>
            </div>
          </div>

          <div v-else-if="error" class="rounded-2xl border border-rose-900/60 bg-rose-950/30 p-8 text-center space-y-3">
            <span class="text-3xl" aria-hidden="true">⚠️</span>
            <h2 class="text-base font-bold text-rose-200">Eita, deu ruim ao buscar as ofertas!</h2>
            <p class="text-xs sm:text-sm text-rose-300 max-w-md mx-auto">{{ error }}</p>
            <div>
              <button
                type="button"
                @click="fetchDeals"
                class="min-h-11 rounded-xl bg-rose-600 hover:bg-rose-500 text-white px-5 py-2 text-xs font-bold transition-colors focus-visible:ring-2 focus-visible:ring-rose-400 focus-visible:outline-none shadow"
              >
                Tentar carregar de novo
              </button>
            </div>
          </div>

          <div v-else-if="filteredDeals.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
            <GameCard
              v-for="deal in filteredDeals"
              :key="deal.id"
              :deal="deal"
            />
          </div>

          <div v-else class="rounded-2xl border border-slate-800 bg-slate-900/70 p-12 text-center space-y-3">
            <span class="text-4xl" aria-hidden="true">🎮</span>
            <h2 class="text-base font-bold text-slate-100">Nenhum jogo encontrado com esses filtros</h2>
            <p class="text-xs sm:text-sm text-slate-400 max-w-md mx-auto">
              Dá uma olhada no termo pesquisado ou tenta mudar a loja e a categoria pra ver outras opções.
            </p>
            <div class="pt-2">
              <button
                type="button"
                @click="resetFilters"
                class="min-h-11 rounded-xl border border-slate-700 bg-slate-800 hover:bg-slate-700 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
              >
                Limpar filtros de busca
              </button>
            </div>
          </div>
        </section>
      </main>
    </div>

    <footer class="border-t border-slate-800/80 bg-slate-950 py-6 text-center text-xs text-slate-400">
      <div class="max-w-7xl mx-auto px-4">
        <p class="mt-1 text-slate-400">Preços e ofertas sincronizados via APIs públicas da Epic Games Store e Steam.</p>
      </div>
    </footer>
  </div>
</template>
