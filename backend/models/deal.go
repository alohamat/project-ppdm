package models

import "time"

type GameDeal struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Store           string     `json:"store"`  // "epic" ou "steam"
	Status          string     `json:"status"` // "free_now", "free_upcoming", "deal"
	OriginalPrice   float64    `json:"original_price"`
	FinalPrice      float64    `json:"final_price"`
	DiscountPercent int        `json:"discount_percent"`
	Currency        string     `json:"currency"`
	ImageUrl        string     `json:"image_url"`
	StoreUrl        string     `json:"store_url"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
}

type DealsResponse struct {
	FreeNow      []GameDeal `json:"free_now"`
	FreeUpcoming []GameDeal `json:"free_upcoming"`
	BigDeals     []GameDeal `json:"big_deals"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
