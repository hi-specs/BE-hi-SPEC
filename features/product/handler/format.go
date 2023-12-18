package handler

type PutProductRequest struct {
	Category  string `json:"category" form:"category"`
	Name      string `json:"name" form:"name"`
	CPU       string `json:"cpu" form:"cpu"`
	RAM       string `json:"ram" form:"ram"`
	Display   string `json:"display" form:"display"`
	Storage   string `json:"storage" form:"storage"`
	Thickness string `json:"thickness" form:"thickness"`
	Weight    string `json:"weight" form:"weight"`
	Bluetooth string `json:"bluetooth" form:"bluetooth"`
	HDMI      string `json:"hdmi" form:"hdmi"`
	Price     int    `json:"price" form:"price"`
	Picture   string `json:"picture" form:"picture"`
}

type ProductRequest struct {
	Laptop   string `json:"laptop" form:"laptop"`
	Picture  string `json:"picture" form:"picture"`
	Category string `json:"category" form:"category"`
}

type ProductResponse struct {
	ID        uint   `json:"product_id"`
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
	ID       uint   `json:"product_id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Picture  string `json:"picture"`
}

type AllResponse struct {
	ID       uint   `json:"product_id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Picture  string `json:"picture"`
}
