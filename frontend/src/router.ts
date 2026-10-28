import { createRouter, createWebHistory } from 'vue-router'
import AuthPage from './views/AuthPage.vue'
import DashboardPage from './views/DashboardPage.vue'

const routes = [
  {
    path: '/',
    name: 'auth',
    component: AuthPage,
    meta: { guestOnly: true }
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

// navigation guard
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    return next({ path: '/' })
  }

  if (to.meta.guestOnly && token) {
    return next({ path: '/dashboard' })
  }

  next()
})

export default router
