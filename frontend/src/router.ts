import { createRouter, createWebHistory } from 'vue-router'
import HomeDealsView from './views/HomeDealsView.vue'
import AuthPage from './views/AuthPage.vue'
import DashboardPage from './views/DashboardPage.vue'

// rotas
const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeDealsView
  },
  {
    path: '/login',
    name: 'login',
    component: AuthPage,
    meta: { guestOnly: true }
  },
  {
    path: '/auth',
    redirect: '/login'
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// guards
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    return next({ path: '/login' })
  }

  if (to.meta.guestOnly && token) {
    return next({ path: '/dashboard' })
  }

  next()
})

export default router
