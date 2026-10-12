<script setup lang="ts">
import GlassDiv from '../components/GlassDiv.vue';
import { ref } from 'vue'
import axios from 'axios';

const registering = ref(false);
const formBody = ref({
  email: "",
  password: ""
})

const register = async () => {
  await axios.post("localhost:8080/register", {email: formBody.value.email, password: formBody.value.password})
  .then((response) => {
    console.log(response.data)
  })
  .catch((error) => {
    console.error(error)
  })
}
</script>

<template>
  <div class="h-screen bg-linear-to-b from-blue-700 via-purple-500 to-orange-500 animate-gradient flex items-center font-sf justify-center">
    <GlassDiv>
      <form class="flex flex-col gap-2">
        <div class="flex flex-col">
          <label for="email">E-mail:</label>
          <input id="email" type="email" placeholder="seu@email.com" v-model="formBody.email" />
        </div>

        <div class="flex flex-col">
          <label for="password">Senha:</label>
          <input id="password" type="password" placeholder="Sua senha" v-model="formBody.password"/>
        </div>

        <div v-if="registering" class="flex flex-col">
          <label for="password">Confirmar Senha:</label>
          <input id="password" type="password" placeholder="Confirme a senha" />
        </div>

        <button type="submit"
          class="bg-white/25 rounded-full hover:scale-110 transition-all duration-200 ease-in-out cursor-pointer w-full my-2">
          {{ registering ? "Cadastrar" : "Login" }}
        </button>
      </form>
      <span class="text-sm">
       {{ registering ? "Já tem uma conta?" : "Não tem uma conta?" }}<span class="text-blue-600 cursor-pointer"
        @click="registering = !registering">{{ registering ? " Fazer Login" : " Cadastre-se" }}</span>
      </span>
    </GlassDiv>
  </div>
</template>

<style lang="css" scoped>
input {
  border: 2px rgb(90, 89, 89) solid;
  background-color: rgba(128, 128, 128, 0.397);
  outline: none;
  padding: 0 4px 0 4px;
}
</style>