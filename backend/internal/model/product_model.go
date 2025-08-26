package model

type MProductRes struct {
	Id         string `json:"id,omitempty"`
	StatusCode string `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

type MGetProductRes struct {
	Id          string  `json:"id,omitempty"`
	StatusCode  string  `json:"status_code,omitempty"`
	Message     string  `json:"message,omitempty"`
	ProductName string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	ImageUrl    string  `json:"image_url,omitempty"`
}

type MListProductRes struct {
	Message    string     `json:"message,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
	Products   []Product  `json:"products,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	ItemPerPage int `json:"item_per_page"`
	TotalPage   int `json:"total_page_count"`
	TotalItem   int `json:"total_item_count"`
}

type Product struct {
	Id          string  `json:"id,omitempty"`
	StatusCode  string  `json:"status_code,omitempty"`
	Message     string  `json:"message,omitempty"`
	ProductName string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	ImageUrl    string  `json:"image_url,omitempty"`
}
