<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
  <div class="min-h-screen bg-slate-950 text-slate-100 font-sf p-4 sm:p-6 lg:p-8">
    <div class="max-w-5xl mx-auto space-y-6">
      
      <header class="bg-slate-900 border border-slate-800 rounded-2xl p-4 sm:p-5 shadow-lg flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-xl bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-lg text-emerald-400 shadow-inner" aria-hidden="true">
            {{ authStore.user?.email?.[0]?.toUpperCase() || 'U' }}
          </div>
          <div>
            <h1 class="font-bold text-base sm:text-lg text-white">{{ authStore.user?.email }}</h1>
            <span
              class="text-[11px] px-2.5 py-0.5 rounded-md font-bold uppercase tracking-wider inline-block mt-0.5"
              :class="authStore.isAdmin ? 'bg-amber-500/20 text-amber-300 border border-amber-500/30' : 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'"
            >
              {{ authStore.isAdmin ? 'Administrador' : 'Jogador' }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2.5">
          <button
            type="button"
            @click="router.push('/')"
            class="flex min-h-11 items-center gap-1.5 rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors hover:border-slate-600 hover:bg-slate-750 hover:text-white focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
          >
            <span>🎮 Ver ofertas de jogos</span>
          </button>

          <button
            type="button"
            @click="handleLogout"
            class="min-h-11 rounded-xl border border-rose-900/50 bg-rose-950/40 px-4 py-2 text-xs font-semibold text-rose-300 transition-colors hover:bg-rose-900/60 hover:text-white focus-visible:ring-2 focus-visible:ring-rose-500 focus-visible:outline-none"
          >
            Sair da conta
          </button>
        </div>
      </header>

      <div
        v-if="errorMessage"
        role="alert"
        class="bg-rose-950/40 border border-rose-900/60 p-4 rounded-xl text-xs sm:text-sm text-rose-300"
      >
        {{ errorMessage }}
      </div>

      <section class="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-sm space-y-3">
        <h2 class="text-xl font-bold text-white">Status da Sessão</h2>
        <p class="text-sm text-slate-300">
          Você tá conectado e autenticado pelo token JWT. Seu ID no banco é:
          <code class="bg-slate-950 border border-slate-800 px-2 py-0.5 rounded text-xs text-emerald-400 font-mono ml-1">{{ authStore.user?.id }}</code>
        </p>
        <p class="text-sm text-slate-400">
          Nível de permissão atual: <strong class="text-slate-200 capitalize">{{ authStore.user?.role === 'admin' ? 'Administrador' : 'Jogador' }}</strong>.
        </p>
      </section>

      <section v-if="authStore.isAdmin" class="space-y-6">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-xl font-bold text-white">Painel do Administrador</h2>
            <p class="text-xs sm:text-sm text-slate-400">Acompanhe as contas e métricas do sistema.</p>
          </div>
          <button
            type="button"
            @click="loadAdminData"
            :disabled="loadingUsers"
            class="min-h-11 rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-200 transition-colors hover:border-slate-600 hover:bg-slate-750 hover:text-white focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none"
          >
            {{ loadingUsers ? 'Carregando dados...' : 'Atualizar métricas' }}
          </button>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4" v-if="stats">
          <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-sm">
            <span class="text-xs uppercase tracking-wider text-slate-400 font-semibold">Total de contas</span>
            <p class="text-3xl font-extrabold mt-1 text-white">{{ stats.total_users }}</p>
          </div>
          <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-sm">
            <span class="text-xs uppercase tracking-wider text-amber-400/80 font-semibold">Administradores</span>
            <p class="text-3xl font-extrabold mt-1 text-amber-300">{{ stats.total_admins }}</p>
          </div>
          <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-sm">
            <span class="text-xs uppercase tracking-wider text-emerald-400/80 font-semibold">Jogadores</span>
            <p class="text-3xl font-extrabold mt-1 text-emerald-400">{{ stats.total_normal }}</p>
          </div>
        </div>

        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 sm:p-6 shadow-sm overflow-hidden">
          <h3 class="font-bold text-base sm:text-lg text-white mb-4">Contas cadastradas no sistema</h3>
          
          <div class="overflow-x-auto -mx-5 px-5 sm:mx-0 sm:px-0">
            <table class="w-full text-left text-sm">
              <thead class="border-b border-slate-800 text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
                <tr>
                  <th class="py-3 px-3">E-mail</th>
                  <th class="py-3 px-3">Papel</th>
                  <th class="py-3 px-3">Identificador (ID)</th>
                  <th class="py-3 px-3">Data de cadastro</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-800/80">
                <tr v-for="user in usersList" :key="user.id" class="hover:bg-slate-850/60 transition-colors">
                  <td class="py-3 px-3 font-medium text-slate-200">{{ user.email }}</td>
                  <td class="py-3 px-3">
                    <span
                      class="text-[11px] px-2 py-0.5 rounded-md font-bold uppercase tracking-wider"
                      :class="user.role === 'admin' ? 'bg-amber-500/20 text-amber-300 border border-amber-500/30' : 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'"
                    >
                      {{ user.role === 'admin' ? 'Admin' : 'Jogador' }}
                    </span>
                  </td>
                  <td class="py-3 px-3 font-mono text-xs text-slate-400">{{ user.id }}</td>
                  <td class="py-3 px-3 text-xs text-slate-400">
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

      <section v-else class="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-sm space-y-2">
        <h3 class="font-bold text-lg text-emerald-300">Espaço do Jogador</h3>
        <p class="text-sm text-slate-300">
          Sua conta tá ativa e você pode navegar à vontade pelas ofertas! As ferramentas de moderação e listagem de usuários são exclusivas dos administradores.
        </p>
      </section>

    </div>
  </div>
</template>
