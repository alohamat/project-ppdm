import { createRouter, createWebHistory } from 'vue-router'
import AuthPage from './views/AuthPage.vue'

const routes = [
  {
    path: '/auth',
    name: 'auth',
    component: AuthPage
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router