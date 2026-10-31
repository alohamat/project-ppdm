<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import GlassDiv from '../components/GlassDiv.vue'
import { useAuthStore } from '../auth'
import { api } from '../api'

const router = useRouter()
const authStore = useAuthStore()

// estado
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

// submit
const handleSubmit = async () => {
  errorMessage.value = ''

  if (!formBody.value.email || !formBody.value.password) {
    errorMessage.value = 'Preencha o e-mail e a senha.'
    return
  }

  if (registering.value) {
    if (formBody.value.password.length < 6) {
      errorMessage.value = 'A senha deve ter pelo menos 6 caracteres.'
      return
    }
    if (formBody.value.password !== formBody.value.confirmPassword) {
      errorMessage.value = 'As senhas não coincidem.'
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
      errorMessage.value = 'Ocorreu um erro ao conectar ao servidor.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="h-screen bg-linear-to-b from-blue-700 via-purple-500 to-orange-500 animate-gradient flex flex-col items-center font-sf justify-center p-4">
    <div class="w-full max-w-sm mb-3">
      <button
        @click="router.push('/')"
        class="text-xs text-white/80 hover:text-white flex items-center gap-1 transition cursor-pointer"
      >
        <span>← Voltar para as ofertas de jogos</span>
      </button>
    </div>

    <GlassDiv class="w-full max-w-sm">
      <h2 class="text-xl font-bold text-center mb-4">
        {{ registering ? 'Criar Nova Conta' : 'Acessar Sistema' }}
      </h2>

      <!-- erro -->
      <div v-if="errorMessage" class="bg-red-500/40 border border-red-500/60 p-2.5 rounded-xl text-xs mb-3 text-white">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-3">
        <div class="flex flex-col">
          <label for="email" class="text-xs mb-1 font-medium">E-mail:</label>
          <input
            id="email"
            type="email"
            required
            placeholder="seu@email.com"
            v-model="formBody.email"
            class="rounded-lg py-1.5 px-3"
          />
        </div>

        <div class="flex flex-col">
          <label for="password" class="text-xs mb-1 font-medium">Senha:</label>
          <input
            id="password"
            type="password"
            required
            placeholder="Sua senha (mín. 6 dígitos)"
            v-model="formBody.password"
            class="rounded-lg py-1.5 px-3"
          />
        </div>

        <div v-if="registering" class="flex flex-col">
          <label for="confirmPassword" class="text-xs mb-1 font-medium">Confirmar Senha:</label>
          <input
            id="confirmPassword"
            type="password"
            required
            placeholder="Confirme sua senha"
            v-model="formBody.confirmPassword"
            class="rounded-lg py-1.5 px-3"
          />
        </div>

        <!-- role -->
        <div v-if="registering" class="flex flex-col">
          <label for="role" class="text-xs mb-1 font-medium">Papel do Usuário:</label>
          <select
            id="role"
            v-model="formBody.role"
            class="rounded-lg py-1.5 px-3 text-sm text-black bg-white/90 outline-none"
          >
            <option value="user">Usuário Comum (User)</option>
            <option value="admin">Administrador (Admin)</option>
          </select>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="bg-white/25 hover:bg-white/35 active:scale-95 rounded-full py-2 font-semibold transition-all duration-200 cursor-pointer w-full mt-2 shadow"
        >
          {{ loading ? 'Processando...' : (registering ? 'Cadastrar' : 'Entrar') }}
        </button>
      </form>

      <div class="text-center text-xs mt-4">
        <span>{{ registering ? 'Já tem uma conta?' : 'Não tem uma conta?' }}</span>
        <button
          type="button"
          class="text-blue-200 hover:text-white underline cursor-pointer ml-1 font-medium"
          @click="toggleMode"
        >
          {{ registering ? 'Fazer Login' : 'Cadastre-se' }}
        </button>
      </div>
    </GlassDiv>
  </div>
</template>

<style lang="css" scoped>
input {
  border: 1px solid rgba(255, 255, 255, 0.3);
  background-color: rgba(255, 255, 255, 0.15);
  outline: none;
  color: white;
}
input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}
input:focus {
  border-color: rgba(255, 255, 255, 0.8);
  background-color: rgba(255, 255, 255, 0.25);
}
</style>
