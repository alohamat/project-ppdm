<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore, type User } from '../auth'
import { useTheme } from '../composables/useTheme'
import { api } from '../api'

const router = useRouter()
const authStore = useAuthStore()
const { isDark, toggleTheme } = useTheme()

interface Stats {
  total_users: number
  total_admins: number
  total_normal: number
}

const stats = ref<Stats | null>(null)
const usersList = ref<User[]>([])
const loadingUsers = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const actionLoadingUserId = ref<string | null>(null)
const userToDelete = ref<User | null>(null)
const isDeleteModalOpen = ref(false)

const statCards = computed(() => {
  if (!stats.value) return []
  return [
    {
      label: 'Total de contas',
      value: stats.value.total_users,
      labelClass: 'text-slate-500 dark:text-slate-400',
      valueClass: 'text-slate-900 dark:text-white'
    },
    {
      label: 'Administradores',
      value: stats.value.total_admins,
      labelClass: 'text-amber-600 dark:text-amber-400/80',
      valueClass: 'text-amber-600 dark:text-amber-300'
    },
    {
      label: 'Jogadores',
      value: stats.value.total_normal,
      labelClass: 'text-emerald-600 dark:text-emerald-400/80',
      valueClass: 'text-emerald-600 dark:text-emerald-400'
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
    errorMessage.value = err.response?.data?.error || 'Não foi possível carregar os dados administrativos.'
  } finally {
    loadingUsers.value = false
  }
}

const handleUpdateRole = async (targetUser: User, newRole: 'admin' | 'user') => {
  if (targetUser.role === newRole) return
  actionLoadingUserId.value = targetUser.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await api.put<{ message: string; user: User }>(`/api/admin/users/${targetUser.id}/role`, { role: newRole })
    targetUser.role = newRole
    successMessage.value = `Cargo de ${targetUser.email} atualizado para ${newRole === 'admin' ? 'Administrador' : 'Jogador'}!`

    if (authStore.user?.id === targetUser.id) {
      authStore.user.role = newRole
      localStorage.setItem('user', JSON.stringify(authStore.user))
    }
    await loadAdminData()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao atualizar cargo do usuário.'
  } finally {
    actionLoadingUserId.value = null
  }
}

const openDeleteModal = (targetUser: User) => {
  userToDelete.value = targetUser
  isDeleteModalOpen.value = true
}

const closeDeleteModal = () => {
  userToDelete.value = null
  isDeleteModalOpen.value = false
}

const confirmDeleteUser = async () => {
  if (!userToDelete.value) return
  const target = userToDelete.value
  actionLoadingUserId.value = target.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await api.delete(`/api/admin/users/${target.id}`)
    if (authStore.user?.id === target.id) {
      authStore.logout()
      router.push('/login')
      return
    }
    usersList.value = usersList.value.filter((u) => u.id !== target.id)
    successMessage.value = `Conta de ${target.email} deletada com sucesso.`
    closeDeleteModal()
    await loadAdminData()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao deletar usuário.'
    closeDeleteModal()
  } finally {
    actionLoadingUserId.value = null
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
  <div class="min-h-screen bg-slate-50 text-slate-900 dark:bg-slate-950 dark:text-slate-100 p-4 font-sf sm:p-6 lg:p-8 transition-colors duration-200">
    <div class="mx-auto max-w-5xl space-y-6">

      <header class="flex flex-col justify-between gap-4 rounded-2xl border border-slate-200 bg-white/90 shadow-sm dark:border-slate-800 dark:bg-slate-900 p-4 sm:flex-row sm:items-center sm:p-5">
        <div class="flex items-center gap-3.5">
          <div class="flex size-11 items-center justify-center rounded-xl border border-slate-200 bg-slate-100 dark:border-slate-700 dark:bg-slate-800 text-lg font-bold text-emerald-600 dark:text-emerald-400 shadow-inner" aria-hidden="true">
            {{ authStore.user?.email?.[0]?.toUpperCase() || 'U' }}
          </div>
          <div>
            <h1 class="text-base font-bold text-slate-900 dark:text-white sm:text-lg">{{ authStore.user?.email }}</h1>
            <span
              class="mt-0.5 inline-block rounded-md border px-2.5 py-0.5 text-xs font-bold uppercase tracking-wider"
              :class="authStore.isAdmin ? 'border-amber-400/40 bg-amber-50 text-amber-800 dark:border-amber-500/30 dark:bg-amber-500/20 dark:text-amber-300' : 'border-emerald-400/40 bg-emerald-50 text-emerald-800 dark:border-emerald-500/30 dark:bg-emerald-500/20 dark:text-emerald-300'"
            >
              {{ authStore.isAdmin ? 'Administrador' : 'Jogador' }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2.5">
          <button
            type="button"
            @click="toggleTheme"
            class="flex size-11 items-center justify-center rounded-xl border border-slate-200 bg-slate-100 text-slate-700 hover:bg-slate-200 dark:border-slate-700 dark:bg-slate-800 dark:text-amber-300 dark:hover:bg-slate-700 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
            :title="isDark ? 'Mudar para tema claro' : 'Mudar para tema escuro'"
            :aria-label="isDark ? 'Ativar modo claro' : 'Ativar modo escuro'"
          >
            <svg v-if="isDark" class="size-4.5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
              <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
            </svg>
            <svg v-else class="size-4.5 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </button>

          <button
            type="button"
            @click="router.push('/')"
            class="flex min-h-11 items-center gap-1.5 rounded-xl border border-slate-200 bg-slate-100 dark:border-slate-700 dark:bg-slate-800 px-4 py-2 text-xs font-semibold text-slate-700 dark:text-slate-200 transition-colors hover:bg-slate-200 dark:hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            <span>🎮 Ofertas</span>
          </button>

          <button
            type="button"
            @click="handleLogout"
            class="min-h-11 rounded-xl border border-rose-300 bg-rose-50 text-rose-700 dark:border-rose-900/50 dark:bg-rose-950/40 dark:text-rose-300 px-4 py-2 text-xs font-semibold transition-colors hover:bg-rose-100 dark:hover:bg-rose-900/60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-500"
          >
            Sair
          </button>
        </div>
      </header>

      <div v-if="errorMessage" role="alert" class="rounded-xl border border-rose-300 bg-rose-50 text-rose-700 dark:border-rose-900/60 dark:bg-rose-950/40 p-4 text-xs dark:text-rose-300 sm:text-sm">
        {{ errorMessage }}
      </div>

      <div v-if="successMessage" role="status" class="rounded-xl border border-emerald-300 bg-emerald-50 text-emerald-800 dark:border-emerald-900/60 dark:bg-emerald-950/40 p-4 text-xs dark:text-emerald-300 sm:text-sm">
        {{ successMessage }}
      </div>

      <section class="space-y-3 rounded-2xl border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900 p-6">
        <h2 class="text-xl font-bold text-slate-900 dark:text-white">Status da Sessão</h2>
        <p class="text-sm text-slate-600 dark:text-slate-300">
          Você está conectado e autenticado pelo token JWT. Seu ID:
          <code class="ml-1 rounded border border-slate-200 bg-slate-100 dark:border-slate-800 dark:bg-slate-950 px-2 py-0.5 font-mono text-xs text-emerald-600 dark:text-emerald-400">{{ authStore.user?.id }}</code>
        </p>
        <p class="text-sm text-slate-500 dark:text-slate-400">
          Nível de permissão atual: <strong class="text-slate-800 dark:text-slate-200 capitalize">{{ authStore.user?.role === 'admin' ? 'Administrador' : 'Jogador' }}</strong>.
        </p>
      </section>

      <section v-if="authStore.isAdmin" class="space-y-6">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-xl font-bold text-slate-900 dark:text-white">Painel do Administrador</h2>
            <p class="text-xs text-slate-500 dark:text-slate-400 sm:text-sm">Acompanhe as contas e métricas do sistema.</p>
          </div>
          <button
            type="button"
            @click="loadAdminData"
            :disabled="loadingUsers"
            class="min-h-11 rounded-xl border border-slate-200 bg-slate-100 text-slate-700 hover:bg-slate-200 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-slate-700 px-4 py-2 text-xs font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            {{ loadingUsers ? 'Carregando...' : 'Atualizar métricas' }}
          </button>
        </div>

        <div v-if="stats" class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div
            v-for="card in statCards"
            :key="card.label"
            class="rounded-2xl border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900 p-5 shadow-sm"
          >
            <span class="text-xs font-semibold uppercase tracking-wider" :class="card.labelClass">{{ card.label }}</span>
            <p class="mt-1 text-3xl font-extrabold" :class="card.valueClass">{{ card.value }}</p>
          </div>
        </div>

        <div class="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900 p-5 sm:p-6">
          <h3 class="mb-4 text-base font-bold text-slate-900 dark:text-white sm:text-lg">Gerenciamento de Contas</h3>

          <div class="-mx-5 overflow-x-auto px-5 sm:mx-0 sm:px-0">
            <table class="w-full text-left text-sm">
              <thead class="border-b border-slate-200 dark:border-slate-800 text-xs font-semibold uppercase tracking-wider text-slate-500 dark:text-slate-400">
                <tr>
                  <th class="px-3 py-3">E-mail</th>
                  <th class="px-3 py-3">Cargo Atual</th>
                  <th class="px-3 py-3">Alterar Cargo</th>
                  <th class="px-3 py-3">Data de cadastro</th>
                  <th class="px-3 py-3 text-right">Ações</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-slate-800/80">
                <tr v-for="user in usersList" :key="user.id" class="transition-colors hover:bg-slate-50 dark:hover:bg-slate-850/60">
                  <td class="px-3 py-3">
                    <span class="font-medium text-slate-800 dark:text-slate-200">{{ user.email }}</span>
                    <span v-if="user.id === authStore.user?.id" class="ml-2 rounded border border-emerald-300 bg-emerald-50 text-emerald-800 dark:border-slate-700 dark:bg-slate-800 dark:text-emerald-400 px-1.5 py-0.5 text-[10px] font-semibold">
                      Você
                    </span>
                  </td>
                  <td class="px-3 py-3">
                    <span
                      class="inline-block rounded-md border px-2 py-0.5 text-xs font-bold uppercase tracking-wider"
                      :class="user.role === 'admin' ? 'border-amber-400/40 bg-amber-50 text-amber-800 dark:border-amber-500/30 dark:bg-amber-500/20 dark:text-amber-300' : 'border-emerald-400/40 bg-emerald-50 text-emerald-800 dark:border-emerald-500/30 dark:bg-emerald-500/20 dark:text-emerald-300'"
                    >
                      {{ user.role === 'admin' ? 'Admin' : 'Jogador' }}
                    </span>
                  </td>
                  <td class="px-3 py-3">
                    <select
                      :value="user.role"
                      :disabled="actionLoadingUserId === user.id"
                      @change="handleUpdateRole(user, ($event.target as HTMLSelectElement).value as 'admin' | 'user')"
                      class="rounded-lg border border-slate-300 bg-white text-slate-800 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200 px-2.5 py-1 text-xs transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 disabled:opacity-50"
                      :aria-label="`Alterar cargo de ${user.email}`"
                    >
                      <option value="user">Jogador</option>
                      <option value="admin">Administrador</option>
                    </select>
                  </td>
                  <td class="px-3 py-3 text-xs text-slate-500 dark:text-slate-400">
                    {{ user.created_at ? new Date(user.created_at).toLocaleString('pt-BR') : 'Indisponível' }}
                  </td>
                  <td class="px-3 py-3 text-right">
                    <button
                      type="button"
                      @click="openDeleteModal(user)"
                      :disabled="actionLoadingUserId === user.id"
                      class="inline-flex items-center gap-1 rounded-lg border border-rose-300 bg-rose-50 text-rose-700 hover:bg-rose-100 dark:border-rose-900/60 dark:bg-rose-950/40 dark:text-rose-300 dark:hover:bg-rose-900/70 dark:hover:text-white px-2.5 py-1 text-xs font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-500 disabled:opacity-50"
                      :title="`Deletar conta de ${user.email}`"
                    >
                      <svg class="size-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      <span>Excluir</span>
                    </button>
                  </td>
                </tr>
                <tr v-if="usersList.length === 0">
                  <td colspan="5" class="py-6 text-center text-slate-400">Nenhuma conta encontrada.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section v-else class="space-y-2 rounded-2xl border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900 p-6">
        <h3 class="text-lg font-bold text-emerald-600 dark:text-emerald-300">Espaço do Jogador</h3>
        <p class="text-sm text-slate-600 dark:text-slate-300">
          Sua conta está ativa e você pode navegar à vontade pelas ofertas! As ferramentas de moderação e listagem de usuários são exclusivas dos administradores.
        </p>
      </section>

    </div>

    <div
      v-if="isDeleteModalOpen && userToDelete"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/80 p-4 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <div class="w-full max-w-md space-y-4 rounded-2xl border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900 p-6 shadow-2xl">
        <div class="flex items-center gap-3">
          <div class="flex size-10 shrink-0 items-center justify-center rounded-xl border border-rose-500/20 bg-rose-500/10 text-xl text-rose-500 dark:text-rose-400" aria-hidden="true">
            ⚠️
          </div>
          <div>
            <h3 id="modal-title" class="text-base font-bold text-slate-900 dark:text-white">Confirmar exclusão de conta</h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Esta ação é irreversível e removerá todos os dados do usuário.</p>
          </div>
        </div>

        <div class="rounded-xl border border-slate-200 bg-slate-50 dark:border-slate-800 dark:bg-slate-950 p-3 text-xs text-slate-700 dark:text-slate-300">
          <p>Você está prestes a deletar a conta de:</p>
          <p class="mt-1 font-semibold text-rose-600 dark:text-rose-300">{{ userToDelete.email }}</p>
          <p v-if="userToDelete.id === authStore.user?.id" class="mt-2 text-amber-600 dark:text-amber-400 font-medium">
            Atenção: Esta é a sua própria conta! Ao confirmar, você será deslogado do sistema imediatamente.
          </p>
        </div>

        <div class="flex items-center justify-end gap-2.5 pt-2">
          <button
            type="button"
            @click="closeDeleteModal"
            class="min-h-10 rounded-xl border border-slate-300 bg-slate-100 hover:bg-slate-200 text-slate-700 dark:border-slate-700 dark:bg-slate-800 dark:hover:bg-slate-700 dark:text-slate-200 px-4 py-2 text-xs font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            Cancelar
          </button>
          <button
            type="button"
            @click="confirmDeleteUser"
            :disabled="actionLoadingUserId === userToDelete.id"
            class="min-h-10 rounded-xl bg-rose-600 px-4 py-2 text-xs font-bold text-white shadow-md transition-colors hover:bg-rose-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-400 disabled:opacity-50"
          >
            {{ actionLoadingUserId === userToDelete.id ? 'Deletando...' : 'Sim, deletar conta' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
