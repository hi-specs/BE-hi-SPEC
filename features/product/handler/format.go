package handler

type ProductRequest struct {
	Laptop   string `json:"laptop" form:"laptop"`
	Picture  string `json:"picture" form:"picture"`
	Category string `json:"category" form:"category"`
}

type ProductResponse struct {
	ID        uint   `json:"id"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	CPU       string `json:"cpu"`
	RAM       string `json:"ram"`
	Display   string `json:"display"`
	Storage   string `json:"storage"`
	Thickness string `json:"thickness"`
	Weight    string `json:"weight"`
	Bluetooth string `json:"bluetooth"`
	HDMI      string `json:"hdmi"`
	Price     int    `json:"price"`
	Picture   string `json:"picture"`
}

type SearchResponse struct {
	ID      uint   `json:"product_id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Picture string `json:"picture"`
}

type AllResponse struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Picture  string `json:"picture"`
}
