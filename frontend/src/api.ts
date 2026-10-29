import axios from 'axios'

export const api = axios.create({
  baseURL: 'https://project-ppdm-backend.onrender.com',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Intercepta todas as requisições para anexar o token JWT no cabeçalho Authorization
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Intercepta respostas para deslogar se o token for inválido/expirado (401)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      if (window.location.pathname !== '/') {
        window.location.href = '/'
      }
    }
    return Promise.reject(error)
  }
)

