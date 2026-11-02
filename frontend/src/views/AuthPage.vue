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
  confirmPassword: ''
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
    errorMessage.value = 'Preencha o e-mail e a senha para continuar.'
    return
  }

  if (registering.value) {
    if (formBody.value.password.length < 6) {
      errorMessage.value = 'A senha precisa ter pelo menos 6 caracteres.'
      return
    }
    if (formBody.value.password !== formBody.value.confirmPassword) {
      errorMessage.value = 'As senhas não coincidem. Verifique e tente novamente.'
      return
    }
  }

  loading.value = true

  try {
    const endpoint = registering.value ? '/api/register' : '/api/login'
    const payload = {
      email: formBody.value.email,
      password: formBody.value.password
    }

    const response = await api.post(endpoint, payload)
    const { token, user } = response.data

    authStore.setAuth(user, token)
    router.push('/dashboard')
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || 'Não foi possível conectar ao servidor. Tente novamente mais tarde.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-slate-50 text-slate-900 dark:bg-slate-950 dark:text-slate-100 p-4 font-sf transition-colors duration-200">
    <div class="mb-4 w-full max-w-md">
      <button
        type="button"
        @click="router.push('/')"
        class="inline-flex items-center gap-1.5 rounded-md py-1 text-xs text-slate-500 hover:text-emerald-600 dark:text-slate-400 dark:hover:text-emerald-400 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
      >
        <span>&larr; Voltar para as ofertas de jogos</span>
      </button>
    </div>

    <div class="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-6 shadow-xl dark:border-slate-800 dark:bg-slate-900 sm:p-8">
      <div class="mb-6 text-center">
        <div class="mb-3 inline-flex size-12 items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10 text-2xl" aria-hidden="true">
          🎮
        </div>
        <h1 class="text-xl font-black text-slate-900 dark:text-white sm:text-2xl">
          {{ registering ? 'Criar sua conta' : 'Entrar na sua conta' }}
        </h1>
        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400 sm:text-sm">
          {{ registering ? 'Leva menos de um minuto e libera seu espaço no sistema.' : 'Digite suas credenciais para acessar seu painel.' }}
        </p>
      </div>

      <div
        v-if="errorMessage"
        role="alert"
        class="mb-4 rounded-xl border border-rose-300 bg-rose-50 p-3 text-xs text-rose-700 dark:border-rose-900/60 dark:bg-rose-950/40 dark:text-rose-300"
      >
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="auth-email" class="mb-1.5 block text-xs font-semibold text-slate-700 dark:text-slate-300">
            E-mail
          </label>
          <input
            id="auth-email"
            type="email"
            required
            autocomplete="email"
            placeholder="exemplo@email.com"
            v-model="formBody.email"
            class="min-h-11 w-full rounded-xl border border-slate-300 bg-slate-50 px-3.5 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 dark:placeholder:text-slate-500 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          />
        </div>

        <div>
          <label for="auth-password" class="mb-1.5 block text-xs font-semibold text-slate-700 dark:text-slate-300">
            Senha
          </label>
          <input
            id="auth-password"
            type="password"
            required
            autocomplete="current-password"
            placeholder="Mínimo de 6 dígitos"
            v-model="formBody.password"
            class="min-h-11 w-full rounded-xl border border-slate-300 bg-slate-50 px-3.5 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 dark:placeholder:text-slate-500 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          />
        </div>

        <div v-if="registering">
          <label for="auth-confirm" class="mb-1.5 block text-xs font-semibold text-slate-700 dark:text-slate-300">
            Confirmar senha
          </label>
          <input
            id="auth-confirm"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Repita sua senha"
            v-model="formBody.confirmPassword"
            class="min-h-11 w-full rounded-xl border border-slate-300 bg-slate-50 px-3.5 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 dark:placeholder:text-slate-500 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
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

      <div class="mt-6 border-t border-slate-200 dark:border-slate-800 pt-4 text-center text-xs text-slate-500 dark:text-slate-400">
        <span>{{ registering ? 'Já tem cadastro por aqui?' : 'Ainda não tem conta?' }}</span>
        <button
          type="button"
          class="ml-1.5 rounded font-semibold text-emerald-600 dark:text-emerald-400 underline hover:text-emerald-500 dark:hover:text-emerald-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          @click="toggleMode"
        >
          {{ registering ? 'Fazer login' : 'Cadastre-se rapidinho' }}
        </button>
      </div>
    </div>
  </div>
</template>
