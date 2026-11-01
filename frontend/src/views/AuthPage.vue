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
  role: 'user'
})

const resetForm = () => {
  errorMessage.value = ''
  formBody.value.email = ''
  formBody.value.password = ''
  formBody.value.confirmPassword = ''
  formBody.value.role = 'user'
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
          role: formBody.value.role
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
  <div class="min-h-screen bg-slate-950 text-slate-100 font-sf flex flex-col items-center justify-center p-4">
    <div class="w-full max-w-md mb-4">
      <button
        type="button"
        @click="router.push('/')"
        class="inline-flex items-center gap-1.5 text-xs text-slate-400 hover:text-emerald-400 transition-colors py-1 focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none rounded-md"
      >
        <span>&larr; Voltar pras ofertas de jogos</span>
      </button>
    </div>

    <div class="w-full max-w-md rounded-2xl border border-slate-800 bg-slate-900 p-6 sm:p-8 shadow-xl">
      <div class="text-center mb-6">
        <div class="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-2xl mb-3" aria-hidden="true">
          🎮
        </div>
        <h1 class="text-xl sm:text-2xl font-black text-white">
          {{ registering ? 'Criar sua conta' : 'Bora entrar' }}
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          {{ registering ? 'Leva menos de um minuto e você libera o acesso ao sistema.' : 'Digite suas credenciais pra acessar seu painel.' }}
        </p>
      </div>

      <div
        v-if="errorMessage"
        role="alert"
        class="mb-4 rounded-xl border border-rose-900/60 bg-rose-950/40 p-3 text-xs text-rose-300"
      >
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="auth-email" class="block text-xs font-semibold text-slate-300 mb-1.5">
            E-mail
          </label>
          <input
            id="auth-email"
            type="email"
            required
            autocomplete="email"
            placeholder="exemplo@email.com"
            v-model="formBody.email"
            class="min-h-[44px] w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 transition-colors focus:border-emerald-500 focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
          />
        </div>

        <div>
          <label for="auth-password" class="block text-xs font-semibold text-slate-300 mb-1.5">
            Senha
          </label>
          <input
            id="auth-password"
            type="password"
            required
            autocomplete="current-password"
            placeholder="Mínimo de 6 dígitos"
            v-model="formBody.password"
            class="min-h-[44px] w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 transition-colors focus:border-emerald-500 focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
          />
        </div>

        <div v-if="registering">
          <label for="auth-confirm" class="block text-xs font-semibold text-slate-300 mb-1.5">
            Confirmar senha
          </label>
          <input
            id="auth-confirm"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Repita sua senha"
            v-model="formBody.confirmPassword"
            class="min-h-[44px] w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 transition-colors focus:border-emerald-500 focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
          />
        </div>

        <div v-if="registering">
          <label for="auth-role" class="block text-xs font-semibold text-slate-300 mb-1.5">
            Tipo de conta
          </label>
          <select
            id="auth-role"
            v-model="formBody.role"
            class="min-h-[44px] w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2.5 text-sm text-slate-200 transition-colors focus:border-emerald-500 focus-visible:ring-2 focus-visible:ring-emerald-500/40 focus-visible:outline-none"
          >
            <option value="user">Jogador (Conta comum)</option>
            <option value="admin">Administrador (Controle do sistema)</option>
          </select>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="min-h-[44px] w-full rounded-xl bg-emerald-500 hover:bg-emerald-400 font-bold text-slate-950 text-sm py-2.5 px-4 transition-colors focus-visible:ring-2 focus-visible:ring-emerald-300 focus-visible:outline-none shadow-md shadow-emerald-500/20 active:scale-[0.98] mt-2"
        >
          {{ loading ? 'Processando...' : (registering ? 'Criar minha conta' : 'Entrar na conta') }}
        </button>
      </form>

      <div class="text-center text-xs text-slate-400 mt-6 pt-4 border-t border-slate-800">
        <span>{{ registering ? 'Já tem cadastro por aqui?' : 'Ainda não tem conta?' }}</span>
        <button
          type="button"
          class="text-emerald-400 hover:text-emerald-300 font-semibold ml-1.5 underline focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:outline-none rounded"
          @click="toggleMode"
        >
          {{ registering ? 'Fazer login' : 'Cadastre-se rapidinho' }}
        </button>
      </div>
    </div>
  </div>
</template>
