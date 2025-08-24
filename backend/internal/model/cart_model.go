package model

type HandlerResponseCart struct {
	Id      string      `json:"id,omitempty"`
	Message string      `json:"message,omitempty"`
	Items   []*ListCart `json:"items,omitempty"`
}

type ListCart struct {
	CartId         string  `json:"cartId"`
	ProductId      string  `json:"productId"`
	ProductName    string  `json:"productName"`
	ProductImagUrl string  `json:"productImagUrl"`
	Price          float64 `json:"price"`
	Quantity       int64   `json:"quantity"`
}
