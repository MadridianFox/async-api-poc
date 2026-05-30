package api

import "time"

type Config struct {
	ListenAddr    string
	ProductHost   string
	BasketHost    string
	RequestTimeout time.Duration
}

type ProductFilter struct {
	MinQty    *int
	MaxQty    *int
	MinPrice  *int
	MaxPrice  *int
	NameLike  *string
	Page      *int
}

type ProductSearchResult struct {
	Products    []Product
	Total       int
	PerPage     int
	CurrentPage int
}

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     int     `json:"price"`
	Qty       int     `json:"qty"`
	CreatedAt ISOTime `json:"created_at"`
	UpdatedAt ISOTime `json:"updated_at"`
}

type Basket struct {
	ID     *int        `json:"id"`
	UserID int         `json:"user_id"`
	Price  int         `json:"price"`
	Status int         `json:"status"`
	Items  []BasketItem `json:"items"`
}

type BasketItem struct {
	ID         int `json:"id"`
	ProductID  int `json:"product_id"`
	PriceByOne int `json:"price_by_one"`
	Qty        int `json:"qty"`
}

type searchResponse struct {
	Data []Product  `json:"data"`
	Meta searchMeta `json:"meta"`
}

type searchMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
}

type basketResponse struct {
	Data Basket `json:"data"`
}

