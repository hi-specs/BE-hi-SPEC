package handler

type ProductRequest struct {
	Laptop  string `json:"laptop" form:"laptop"`
	Picture string `json:"picture" form:"picture"`
}

type ProductResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CPU       string `json:"cpu"`
	RAM       string `json:"ram"`
	Display   string `json:"display"`
	Storage   string `json:"storage"`
	Thickness string `json:"thickness"`
	Weight    string `json:"weight"`
	Bluetooth string `json:"bluetooth"`
	HDMI      string `json:"hdmi"`
	Price     string `json:"price"`
	Picture   string `json:"picture"`
}
