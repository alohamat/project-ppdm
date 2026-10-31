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

// estado
const stats = ref<Stats | null>(null)
const usersList = ref<User[]>([])
const loadingUsers = ref(false)
const errorMessage = ref('')

// logout
const handleLogout = () => {
  authStore.logout()
  router.push('/')
}

// admin
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
    errorMessage.value = err.response?.data?.error || 'Erro ao carregar dados administrativos.'
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
  <div class="min-h-screen bg-linear-to-b from-blue-700 via-purple-500 to-orange-500 text-white font-sf p-6">
    <div class="max-w-5xl mx-auto space-y-6">
      
      <!-- header -->
      <header class="bg-white/10 border border-white/20 rounded-2xl p-4 backdrop-blur-md shadow-lg flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center font-bold text-lg">
            {{ authStore.user?.email?.[0]?.toUpperCase() || 'U' }}
          </div>
          <div>
            <h1 class="font-semibold text-lg">{{ authStore.user?.email }}</h1>
            <span
              class="text-xs px-2.5 py-0.5 rounded-full font-bold uppercase tracking-wider inline-block"
              :class="authStore.isAdmin ? 'bg-amber-400 text-amber-950' : 'bg-blue-300 text-blue-950'"
            >
              {{ authStore.user?.role }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="router.push('/')"
            class="flex items-center gap-1.5 rounded-xl border border-white/20 bg-white/15 px-3.5 py-2 text-xs font-semibold hover:bg-white/25 transition cursor-pointer"
          >
            <span>🎮 Ver Ofertas de Jogos</span>
          </button>

          <button
            @click="handleLogout"
            class="bg-red-500/80 hover:bg-red-500 text-white px-4 py-2 rounded-xl text-xs font-semibold transition cursor-pointer shadow"
          >
            Sair da Conta
          </button>
        </div>
      </header>

      <!-- erro -->
      <div v-if="errorMessage" class="bg-red-500/30 border border-red-500/50 p-3 rounded-xl text-sm">
        {{ errorMessage }}
      </div>

      <!-- perfil -->
      <section class="bg-white/10 border border-white/20 rounded-2xl p-6 backdrop-blur-md shadow-lg">
        <h2 class="text-2xl font-bold mb-2">Painel Geral</h2>
        <p class="text-white/80">
          Você está autenticado no sistema via JWT! Seu ID de usuário é:
          <code class="bg-black/30 px-2 py-0.5 rounded text-xs ml-1">{{ authStore.user?.id }}</code>
        </p>
        <p class="text-white/80 mt-1">
          Nível de acesso concedido: <strong class="capitalize">{{ authStore.user?.role }}</strong>.
        </p>
      </section>

      <!-- admin -->
      <section v-if="authStore.isAdmin" class="space-y-6">
        <div class="flex items-center justify-between">
          <h2 class="text-2xl font-bold">Gerenciamento Administrativo</h2>
          <button
            @click="loadAdminData"
            :disabled="loadingUsers"
            class="bg-white/20 hover:bg-white/30 text-xs px-3 py-1.5 rounded-lg cursor-pointer transition"
          >
            {{ loadingUsers ? 'Carregando...' : 'Atualizar Dados' }}
          </button>
        </div>

        <!-- stats -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4" v-if="stats">
          <div class="bg-white/10 border border-white/20 rounded-xl p-4 backdrop-blur-md">
            <span class="text-xs uppercase text-white/70">Total de Usuários</span>
            <p class="text-3xl font-extrabold mt-1">{{ stats.total_users }}</p>
          </div>
          <div class="bg-white/10 border border-white/20 rounded-xl p-4 backdrop-blur-md">
            <span class="text-xs uppercase text-white/70">Administradores</span>
            <p class="text-3xl font-extrabold mt-1 text-amber-300">{{ stats.total_admins }}</p>
          </div>
          <div class="bg-white/10 border border-white/20 rounded-xl p-4 backdrop-blur-md">
            <span class="text-xs uppercase text-white/70">Usuários Padrão</span>
            <p class="text-3xl font-extrabold mt-1 text-blue-200">{{ stats.total_normal }}</p>
          </div>
        </div>

        <!-- tabela -->
        <div class="bg-white/10 border border-white/20 rounded-2xl p-6 backdrop-blur-md shadow-lg overflow-x-auto">
          <h3 class="font-semibold text-lg mb-4">Lista de Usuários no Banco de Dados</h3>
          
          <table class="w-full text-left text-sm">
            <thead class="border-b border-white/20 text-white/70 uppercase text-xs">
              <tr>
                <th class="py-2 px-3">E-mail</th>
                <th class="py-2 px-3">Papel</th>
                <th class="py-2 px-3">ID (Mongo ObjectID)</th>
                <th class="py-2 px-3">Data de Cadastro</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-white/10">
              <tr v-for="user in usersList" :key="user.id" class="hover:bg-white/5 transition">
                <td class="py-3 px-3 font-medium">{{ user.email }}</td>
                <td class="py-3 px-3">
                  <span
                    class="text-xs px-2 py-0.5 rounded font-bold uppercase"
                    :class="user.role === 'admin' ? 'bg-amber-400 text-amber-950' : 'bg-blue-300 text-blue-950'"
                  >
                    {{ user.role }}
                  </span>
                </td>
                <td class="py-3 px-3 font-mono text-xs opacity-80">{{ user.id }}</td>
                <td class="py-3 px-3 text-xs opacity-80">
                  {{ user.created_at ? new Date(user.created_at).toLocaleString() : '-' }}
                </td>
              </tr>
              <tr v-if="usersList.length === 0">
                <td colspan="4" class="py-4 text-center text-white/60">
                  Nenhum usuário encontrado.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- user -->
      <section v-else class="bg-white/10 border border-white/20 rounded-2xl p-6 backdrop-blur-md shadow-lg">
        <h3 class="font-semibold text-lg text-amber-200">Acesso Padrão</h3>
        <p class="text-sm text-white/80 mt-1">
          Sua conta é do tipo <strong>user</strong>. Recursos de gerenciamento e relatórios do sistema são exclusivos para administradores.
        </p>
      </section>

    </div>
  </div>
</template>
