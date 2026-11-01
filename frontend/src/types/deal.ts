export interface GameDeal {
  id: string
  title: string
  store: 'epic' | 'steam'
  status: 'free_now' | 'free_upcoming' | 'deal'
  original_price: number
  final_price: number
  discount_percent: number
  currency: string
  image_url: string
  store_url: string
  start_date?: string
  end_date?: string
}

export interface DealsResponse {
  free_now: GameDeal[]
  free_upcoming: GameDeal[]
  big_deals: GameDeal[]
  updated_at: string
}

