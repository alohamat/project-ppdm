<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../auth'
import { api } from '../api'

const router = useRouter()
const authStore = useAuthStore()

const registering = ref(false)
const errorMessage = ref('')
const loading = ref(false)

const formBody = ref({
  email: '',
  password: '',
  confirmPassword: '',
})

const resetForm = () => {
  errorMessage.value = ''
  formBody.value.email = ''
  formBody.value.password = ''
  formBody.value.confirmPassword = ''
}

const toggleMode = () => {
  registering.value = !registering.value
  resetForm()
}

const handleSubmit = async () => {
  errorMessage.value = ''

  if (!formBody.value.email || !formBody.value.password) {
    errorMessage.value = 'Preencha o e-mail e a senha pra continuar.'
    return
  }

  if (registering.value) {
    if (formBody.value.password.length < 6) {
      errorMessage.value = 'A senha precisa ter pelo menos 6 caracteres.'
      return
    }
    if (formBody.value.password !== formBody.value.confirmPassword) {
      errorMessage.value = 'As duas senhas não estão iguais. Confere aí.'
      return
    }
  }

  loading.value = true

  try {
    const endpoint = registering.value ? '/api/register' : '/api/login'
    const payload = registering.value
      ? {
          email: formBody.value.email,
          password: formBody.value.password,
        }
      : {
          email: formBody.value.email,
          password: formBody.value.password
        }

    const response = await api.post(endpoint, payload)
    const { token, user } = response.data

    authStore.setAuth(user, token)
    router.push('/dashboard')
  } catch (error: any) {
    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error
    } else {
      errorMessage.value = 'Não deu pra conectar com o servidor. Tenta de novo daqui a pouco.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-slate-950 p-4 font-sf text-slate-100">
    <!-- Navigation Back Link -->
    <div class="mb-4 w-full max-w-md">
      <button
        type="button"
        @click="router.push('/')"
        class="inline-flex items-center gap-1.5 rounded-md py-1 text-xs text-slate-400 transition-colors hover:text-emerald-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
      >
        <span>&larr; Voltar pras ofertas de jogos</span>
      </button>
    </div>

    <!-- Auth Card -->
    <div class="w-full max-w-md rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-xl sm:p-8">
      <div class="mb-6 text-center">
        <div class="mb-3 inline-flex size-12 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10 text-2xl" aria-hidden="true">
          🎮
        </div>
        <h1 class="text-xl font-black text-white sm:text-2xl">
          {{ registering ? 'Criar sua conta' : 'Bora entrar' }}
        </h1>
        <p class="mt-1 text-xs text-slate-400 sm:text-sm">
          {{ registering ? 'Leva menos de um minuto e você libera o acesso ao sistema.' : 'Digite suas credenciais pra acessar seu painel.' }}
        </p>
      </div>

      <!-- Error State -->
      <div
        v-if="errorMessage"
        role="alert"
        class="mb-4 rounded-xl border border-rose-900/60 bg-rose-950/40 p-3 text-xs text-rose-300"
      >
        {{ errorMessage }}
      </div>

      <!-- Auth Form -->
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="auth-email" class="mb-1.5 block text-xs font-semibold text-slate-300">
            E-mail
          </label>
          <input
            id="auth-email"
            type="email"
            required
            autocomplete="email"
            placeholder="exemplo@email.com"
            v-model="formBody.email"
            class="min-h-11 w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-100 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40"
          />
        </div>

        <div>
          <label for="auth-password" class="mb-1.5 block text-xs font-semibold text-slate-300">
            Senha
          </label>
          <input
            id="auth-password"
            type="password"
            required
            autocomplete="current-password"
            placeholder="Mínimo de 6 dígitos"
            v-model="formBody.password"
            class="min-h-11 w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-100 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40"
          />
        </div>

        <div v-if="registering">
          <label for="auth-confirm" class="mb-1.5 block text-xs font-semibold text-slate-300">
            Confirmar senha
          </label>
          <input
            id="auth-confirm"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Repita sua senha"
            v-model="formBody.confirmPassword"
            class="min-h-11 w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-100 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40"
          />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="mt-2 min-h-11 w-full rounded-xl bg-emerald-500 px-4 py-2.5 text-sm font-bold text-slate-950 shadow-md shadow-emerald-500/20 transition-colors hover:bg-emerald-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-300 active:scale-[0.98]"
        >
          {{ loading ? 'Processando...' : (registering ? 'Criar minha conta' : 'Entrar na conta') }}
        </button>
      </form>

      <!-- Toggle Mode Link -->
      <div class="mt-6 border-t border-slate-800 pt-4 text-center text-xs text-slate-400">
        <span>{{ registering ? 'Já tem cadastro por aqui?' : 'Ainda não tem conta?' }}</span>
        <button
          type="button"
          class="ml-1.5 rounded font-semibold text-emerald-400 underline hover:text-emerald-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          @click="toggleMode"
        >
          {{ registering ? 'Fazer login' : 'Cadastre-se rapidinho' }}
        </button>
      </div>
    </div>
  </div>
</template>