<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore, type User } from '../auth'
import { api } from '../api'

const router = useRouter()
const authStore = useAuthStore()

interface Stats {
  total_users: number
  total_admins: number
  total_normal: number
}

const stats = ref<Stats | null>(null)
const usersList = ref<User[]>([])
const loadingUsers = ref(false)
const errorMessage = ref('')

const statCards = computed(() => {
  if (!stats.value) return []
  return [
    {
      label: 'Total de contas',
      value: stats.value.total_users,
      labelClass: 'text-slate-400',
      valueClass: 'text-white'
    },
    {
      label: 'Administradores',
      value: stats.value.total_admins,
      labelClass: 'text-amber-400/80',
      valueClass: 'text-amber-300'
    },
    {
      label: 'Jogadores',
      value: stats.value.total_normal,
      labelClass: 'text-emerald-400/80',
      valueClass: 'text-emerald-400'
    }
  ]
})

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}

const loadAdminData = async () => {
  if (!authStore.isAdmin) return

  loadingUsers.value = true
  errorMessage.value = ''
  try {
    const [statsRes, usersRes] = await Promise.all([
      api.get<Stats>('/api/admin/stats'),
      api.get<User[]>('/api/admin/users')
    ])
    stats.value = statsRes.data
    usersList.value = usersRes.data
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Não deu pra puxar os dados de admin agora.'
  } finally {
    loadingUsers.value = false
  }
}

onMounted(async () => {
  try {
    await authStore.fetchProfile()
    if (authStore.isAdmin) {
      await loadAdminData()
    }
  } catch (err) {
    console.error(err)
  }
})
</script>

<template>
  <div class="min-h-screen bg-slate-950 p-4 font-sf text-slate-100 sm:p-6 lg:p-8">
    <div class="mx-auto max-w-5xl space-y-6">

      <!-- User Session Header -->
      <header class="flex flex-col justify-between gap-4 rounded-2xl border border-slate-800 bg-slate-900 p-4 shadow-lg sm:flex-row sm:items-center sm:p-5">
        <div class="flex items-center gap-3.5">
          <div class="flex size-11 items-center justify-center rounded-xl border border-slate-700 bg-slate-800 text-lg font-bold text-emerald-400 shadow-inner" aria-hidden="true">
            {{ authStore.user?.email?.[0]?.toUpperCase() || 'U' }}
          </div>
          <div>
            <h1 class="text-base font-bold text-white sm:text-lg">{{ authStore.user?.email }}</h1>
            <span
              class="mt-0.5 inline-block rounded-md border px-2.5 py-0.5 text-xs font-bold uppercase tracking-wider"
              :class="authStore.isAdmin ? 'border-amber-500/30 bg-amber-500/20 text-amber-300' : 'border-emerald-500/30 bg-emerald-500/20 text-emerald-300'"
            >
              {{ authStore.isAdmin ? 'Administrador' : 'Jogador' }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2.5">
          <button
            type="button"
            @click="router.push('/')"
            class="flex min-h-11 items-center gap-1.5 rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors hover:border-slate-600 hover:bg-slate-700 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            <span>🎮 Ver ofertas de jogos</span>
          </button>

          <button
            type="button"
            @click="handleLogout"
            class="min-h-11 rounded-xl border border-rose-900/50 bg-rose-950/40 px-4 py-2 text-xs font-semibold text-rose-300 transition-colors hover:bg-rose-900/60 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-500"
          >
            Sair da conta
          </button>
        </div>
      </header>

      <!-- Error Notification -->
      <div
        v-if="errorMessage"
        role="alert"
        class="rounded-xl border border-rose-900/60 bg-rose-950/40 p-4 text-xs text-rose-300 sm:text-sm"
      >
        {{ errorMessage }}
      </div>

      <!-- Session Info -->
      <section class="space-y-3 rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-sm">
        <h2 class="text-xl font-bold text-white">Status da Sessão</h2>
        <p class="text-sm text-slate-300">
          Você tá conectado e autenticado pelo token JWT. Seu ID no banco é:
          <code class="ml-1 rounded border border-slate-800 bg-slate-950 px-2 py-0.5 font-mono text-xs text-emerald-400">{{ authStore.user?.id }}</code>
        </p>
        <p class="text-sm text-slate-400">
          Nível de permissão atual: <strong class="text-slate-200 capitalize">{{ authStore.user?.role === 'admin' ? 'Administrador' : 'Jogador' }}</strong>.
        </p>
      </section>

      <!-- Admin Panel Section -->
      <section v-if="authStore.isAdmin" class="space-y-6">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-xl font-bold text-white">Painel do Administrador</h2>
            <p class="text-xs text-slate-400 sm:text-sm">Acompanhe as contas e métricas do sistema.</p>
          </div>
          <button
            type="button"
            @click="loadAdminData"
            :disabled="loadingUsers"
            class="min-h-11 rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors hover:border-slate-600 hover:bg-slate-700 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            {{ loadingUsers ? 'Carregando dados...' : 'Atualizar métricas' }}
          </button>
        </div>

        <!-- Metric Cards -->
        <div v-if="stats" class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div
            v-for="card in statCards"
            :key="card.label"
            class="rounded-2xl border border-slate-800 bg-slate-900 p-5 shadow-sm"
          >
            <span class="text-xs font-semibold uppercase tracking-wider" :class="card.labelClass">{{ card.label }}</span>
            <p class="mt-1 text-3xl font-extrabold" :class="card.valueClass">{{ card.value }}</p>
          </div>
        </div>

        <!-- Registered Users Table -->
        <div class="overflow-hidden rounded-2xl border border-slate-800 bg-slate-900 p-5 shadow-sm sm:p-6">
          <h3 class="mb-4 text-base font-bold text-white sm:text-lg">Contas cadastradas no sistema</h3>

          <div class="-mx-5 overflow-x-auto px-5 sm:mx-0 sm:px-0">
            <table class="w-full text-left text-sm">
              <thead class="border-b border-slate-800 text-xs font-semibold uppercase tracking-wider text-slate-400">
                <tr>
                  <th class="px-3 py-3">E-mail</th>
                  <th class="px-3 py-3">Papel</th>
                  <th class="px-3 py-3">Identificador (ID)</th>
                  <th class="px-3 py-3">Data de cadastro</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-800/80">
                <tr v-for="user in usersList" :key="user.id" class="transition-colors hover:bg-slate-850/60">
                  <td class="px-3 py-3 font-medium text-slate-200">{{ user.email }}</td>
                  <td class="px-3 py-3">
                    <span
                      class="inline-block rounded-md border px-2 py-0.5 text-xs font-bold uppercase tracking-wider"
                      :class="user.role === 'admin' ? 'border-amber-500/30 bg-amber-500/20 text-amber-300' : 'border-emerald-500/30 bg-emerald-500/20 text-emerald-300'"
                    >
                      {{ user.role === 'admin' ? 'Admin' : 'Jogador' }}
                    </span>
                  </td>
                  <td class="px-3 py-3 font-mono text-xs text-slate-400">{{ user.id }}</td>
                  <td class="px-3 py-3 text-xs text-slate-400">
                    {{ user.created_at ? new Date(user.created_at).toLocaleString('pt-BR') : 'Indisponível' }}
                  </td>
                </tr>
                <tr v-if="usersList.length === 0">
                  <td colspan="4" class="py-6 text-center text-slate-400">
                    Nenhuma conta encontrada até o momento.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- Player Space Section -->
      <section v-else class="space-y-2 rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-sm">
        <h3 class="text-lg font-bold text-emerald-300">Espaço do Jogador</h3>
        <p class="text-sm text-slate-300">
          Sua conta tá ativa e você pode navegar à vontade pelas ofertas! As ferramentas de moderação e listagem de usuários são exclusivas dos administradores.
        </p>
      </section>

    </div>
  </div>
</template>