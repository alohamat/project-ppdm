/**
 * Utilitários para formatação de moeda, datas e temporizadores de expiração.
 */

export function formatCurrencyBRL(val: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  }).format(val)
}

export interface ExpirationTimeResult {
  text: string
  isUrgent: boolean
  remainingMs: number
}

export function getRemainingTime(dateString?: string): ExpirationTimeResult | null {
  if (!dateString) return null
  const target = new Date(dateString).getTime()
  if (isNaN(target)) return null

  const now = Date.now()
  const diff = target - now
  if (diff <= 0) {
    return { text: 'Expirando!', isUrgent: true, remainingMs: 0 }
  }

  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))

  const isUrgent = days <= 1

  if (days > 1) {
    return { text: `Expira em ${days}d ${hours}h`, isUrgent, remainingMs: diff }
  } else if (days === 1) {
    return { text: `Expira em 1d ${hours}h`, isUrgent, remainingMs: diff }
  } else if (hours > 0) {
    return { text: `Expira em ${hours}h ${minutes}m!`, isUrgent, remainingMs: diff }
  } else {
    return { text: `Expira em ${minutes}m!`, isUrgent, remainingMs: diff }
  }
}

export function formatUpcomingDate(dateString?: string): string {
  if (!dateString) return 'Em breve'
  const d = new Date(dateString)
  if (isNaN(d.getTime())) return 'Em breve'

  const dateFormatted = d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' })
  const timeFormatted = d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
  return `Libera em ${dateFormatted} às ${timeFormatted}`
}
