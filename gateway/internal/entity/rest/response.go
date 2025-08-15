package rest

// RestResult for proto-api response
type RestResult struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data,omitempty"`
}

// RestResultSingle for proto-api response
type RestResultSingle struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ReqService struct {
	SvcName string
	SvcFunc string
	Headers map[string]string
	Params  map[string]string
}
