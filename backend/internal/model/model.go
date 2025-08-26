package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ParseParamJWT struct {
	DataSession map[string]string
}

type ParseParamJWT2 struct {
	CustomerId string `json:"customerId"`
	Name       string `json:"name"`
	FullName   string `json:"fullName"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

type MResponse struct {
	Id         string `json:"id,omitempty"`
	StatusCode string `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}
