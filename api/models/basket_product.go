package models

type BasketProduct struct {
	ID string		`json:"id"`
	BasketID string `json:"basket_id"`
	ProductID string`json:"product_id"`
	Quantity string	`json:"quantity"`
}

type BasketProductResponse struct {
	BasketProducts []BasketProduct `json:"basket_products"`
	Count int					   `json:"count"`
}

type UpdateBasketProduct struct {
	ID string		`json:"id"`
	BasketID string `json:"basket_id"`
	ProductID string`json:"product_id"`
	Quantity string	`json:"quantity"`
}

type CreateBasketProduct struct {
	BasketID string `json:"basket_id"`
	ProductID string`json:"product_id"`
	Quantity string	`json:"quantity"`
}