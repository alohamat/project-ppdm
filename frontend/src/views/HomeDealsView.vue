<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { useAuthStore } from '../auth'
import GameCard from '../components/GameCard.vue'
import type { GameDeal, DealsResponse } from '../types/deal'

const router = useRouter()
const authStore = useAuthStore()

// estado
const dealsData = ref<DealsResponse | null>(null)
const loading = ref(true)
const error = ref('')

// filtros
const activeTab = ref<'all' | 'free_now' | 'free_upcoming' | 'big_deals'>('all')
const selectedStore = ref<'all' | 'epic' | 'steam'>('all')
const searchQuery = ref('')
const sortBy = ref<'discount' | 'price' | 'name'>('discount')

// carregar
const fetchDeals = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get<DealsResponse>('/api/deals')
    dealsData.value = res.data
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Não foi possível carregar as ofertas.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDeals()
})

// lista unificada
const allDeals = computed<GameDeal[]>(() => {
  if (!dealsData.value) return []
  return [
    ...dealsData.value.free_now,
    ...dealsData.value.free_upcoming,
    ...dealsData.value.big_deals
  ]
})

// contagens
const counts = computed(() => {
  return {
    all: allDeals.value.length,
    free_now: dealsData.value?.free_now.length || 0,
    free_upcoming: dealsData.value?.free_upcoming.length || 0,
    big_deals: dealsData.value?.big_deals.length || 0
  }
})

// filtragem
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

  // loja
  if (selectedStore.value !== 'all') {
    list = list.filter((item) => item.store === selectedStore.value)
  }

  // busca
  if (searchQuery.value.trim() !== '') {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter((item) => item.title.toLowerCase().includes(q))
  }

  // ordenacao
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
  <div class="min-h-screen bg-linear-to-b from-blue-700 via-purple-500 to-orange-500 text-white font-sf p-4 sm:p-6 lg:p-8">
    <div class="max-w-7xl mx-auto space-y-6">

      <!-- navbar -->
      <header class="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-white/20 bg-white/10 p-4 sm:p-5 backdrop-blur-md shadow-lg">
        <div class="flex items-center gap-3">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-white/20 text-xl font-bold shadow-inner">
            🎮
          </div>
          <div>
            <h1 class="text-xl sm:text-2xl font-black tracking-tight">GameDeals</h1>
            <p class="text-xs text-white/70">Jogos gratuitos e super descontos da Epic Games & Steam</p>
          </div>
        </div>

        <!-- acoes -->
        <div class="flex items-center gap-2.5">
          <button
            @click="fetchDeals"
            :disabled="loading"
            class="flex items-center gap-1.5 rounded-xl border border-white/20 bg-white/10 px-3 py-2 text-xs font-semibold hover:bg-white/20 transition cursor-pointer active:scale-95"
            title="Atualizar lista de jogos"
          >
            <svg class="h-4 w-4" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span class="hidden sm:inline">Atualizar</span>
          </button>

          <!-- user autenticado -->
          <template v-if="authStore.isAuthenticated">
            <button
              @click="router.push('/dashboard')"
              class="flex items-center gap-1.5 rounded-xl border border-white/20 bg-white/15 px-3 py-2 text-xs font-semibold hover:bg-white/25 transition cursor-pointer"
            >
              <span>{{ authStore.isAdmin ? 'Painel Admin' : 'Meu Painel' }}</span>
            </button>
            <button
              @click="handleLogout"
              class="rounded-xl bg-red-500/70 hover:bg-red-500 px-3 py-2 text-xs font-semibold transition cursor-pointer"
            >
              Sair
            </button>
          </template>

          <!-- guest -->
          <template v-else>
            <button
              @click="router.push('/login')"
              class="rounded-xl border border-white/30 bg-white/20 hover:bg-white/30 px-4 py-2 text-xs font-bold transition cursor-pointer shadow"
            >
              Entrar / Cadastro
            </button>
          </template>
        </div>
      </header>

      <!-- abas e filtros -->
      <section class="space-y-4 rounded-2xl border border-white/15 bg-white/10 p-4 sm:p-5 backdrop-blur-md shadow-md">
        <!-- abas -->
        <div class="flex flex-wrap items-center gap-2 border-b border-white/15 pb-4">
          <button
            @click="activeTab = 'all'"
            class="flex items-center gap-1.5 rounded-xl px-3.5 py-1.5 text-xs font-bold transition cursor-pointer"
            :class="activeTab === 'all' ? 'bg-white text-zinc-900 shadow-md' : 'bg-white/10 text-white/80 hover:bg-white/20'"
          >
            <span>Todos</span>
            <span class="rounded-full px-1.5 py-0.2 text-[10px]" :class="activeTab === 'all' ? 'bg-zinc-200 text-zinc-800' : 'bg-white/20 text-white'">
              {{ counts.all }}
            </span>
          </button>

          <button
            @click="activeTab = 'free_now'"
            class="flex items-center gap-1.5 rounded-xl px-3.5 py-1.5 text-xs font-bold transition cursor-pointer"
            :class="activeTab === 'free_now' ? 'bg-emerald-500 text-white shadow-md' : 'bg-white/10 text-white/80 hover:bg-white/20'"
          >
            <span>Grátis Agora</span>
            <span class="rounded-full px-1.5 py-0.2 text-[10px]" :class="activeTab === 'free_now' ? 'bg-emerald-700 text-white' : 'bg-white/20 text-white'">
              {{ counts.free_now }}
            </span>
          </button>

          <button
            @click="activeTab = 'free_upcoming'"
            class="flex items-center gap-1.5 rounded-xl px-3.5 py-1.5 text-xs font-bold transition cursor-pointer"
            :class="activeTab === 'free_upcoming' ? 'bg-purple-500 text-white shadow-md' : 'bg-white/10 text-white/80 hover:bg-white/20'"
          >
            <span>Grátis em Breve</span>
            <span class="rounded-full px-1.5 py-0.2 text-[10px]" :class="activeTab === 'free_upcoming' ? 'bg-purple-700 text-white' : 'bg-white/20 text-white'">
              {{ counts.free_upcoming }}
            </span>
          </button>

          <button
            @click="activeTab = 'big_deals'"
            class="flex items-center gap-1.5 rounded-xl px-3.5 py-1.5 text-xs font-bold transition cursor-pointer"
            :class="activeTab === 'big_deals' ? 'bg-amber-400 text-amber-950 shadow-md' : 'bg-white/10 text-white/80 hover:bg-white/20'"
          >
            <span>Grandes Ofertas</span>
            <span class="rounded-full px-1.5 py-0.2 text-[10px]" :class="activeTab === 'big_deals' ? 'bg-amber-600 text-white' : 'bg-white/20 text-white'">
              {{ counts.big_deals }}
            </span>
          </button>
        </div>

        <!-- controles -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-12 gap-3 pt-1">
          <!-- busca -->
          <div class="lg:col-span-6 relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Buscar pelo título do jogo..."
              class="w-full rounded-xl border border-white/20 bg-white/10 py-2 pl-9 pr-3 text-xs text-white placeholder-white/50 outline-none focus:border-white/50 focus:bg-white/20 transition"
            />
            <svg class="absolute left-3 top-2.5 h-3.5 w-3.5 text-white/50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>

          <!-- loja -->
          <div class="lg:col-span-3">
            <select
              v-model="selectedStore"
              class="w-full rounded-xl border border-white/20 bg-white/10 py-2 px-3 text-xs text-white outline-none focus:border-white/50 focus:bg-white/20 transition cursor-pointer"
            >
              <option value="all" class="text-zinc-900">Todas as Lojas</option>
              <option value="epic" class="text-zinc-900">Epic Games</option>
              <option value="steam" class="text-zinc-900">Steam</option>
            </select>
          </div>

          <!-- ordenacao -->
          <div class="lg:col-span-3">
            <select
              v-model="sortBy"
              class="w-full rounded-xl border border-white/20 bg-white/10 py-2 px-3 text-xs text-white outline-none focus:border-white/50 focus:bg-white/20 transition cursor-pointer"
            >
              <option value="discount" class="text-zinc-900">Maior Desconto</option>
              <option value="price" class="text-zinc-900">Menor Preço</option>
              <option value="name" class="text-zinc-900">Nome (A - Z)</option>
            </select>
          </div>
        </div>
      </section>

      <!-- loading -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
        <div v-for="n in 8" :key="n" class="animate-pulse rounded-2xl border border-white/15 bg-white/10 p-4 space-y-3">
          <div class="aspect-16/9 w-full rounded-xl bg-white/10"></div>
          <div class="h-4 w-3/4 rounded bg-white/15"></div>
          <div class="h-3 w-1/2 rounded bg-white/10"></div>
        </div>
      </div>

      <!-- erro -->
      <div v-else-if="error" class="rounded-2xl border border-red-500/40 bg-red-500/20 p-6 text-center backdrop-blur-md space-y-3">
        <p class="text-sm font-semibold">{{ error }}</p>
        <button
          @click="fetchDeals"
          class="rounded-xl bg-white/20 hover:bg-white/30 px-4 py-2 text-xs font-bold transition cursor-pointer"
        >
          Tentar Novamente
        </button>
      </div>

      <!-- lista de cards -->
      <div v-else-if="filteredDeals.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
        <GameCard
          v-for="deal in filteredDeals"
          :key="deal.id"
          :deal="deal"
        />
      </div>

      <!-- vazio -->
      <div v-else class="rounded-2xl border border-white/15 bg-white/10 p-12 text-center backdrop-blur-md space-y-2">
        <span class="text-3xl">🔍</span>
        <h3 class="text-base font-bold">Nenhum jogo encontrado</h3>
        <p class="text-xs text-white/70">Tente ajustar seus termos de busca ou mudar os filtros de loja.</p>
      </div>

    </div>
  </div>
</template>
