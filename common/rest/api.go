package rest

import "encoding/json"

type APIREQ struct {
	Headers map[string]string `json:"headers"`
	Params  map[string]string `json:"params"`
	Body    json.RawMessage   `json:"body"` // raw JSON, bisa di-unmarshal sesuai kebutuhan di target services
}

type APIRES struct {
	ResponseCode string          `json:"responseCode"`
	Message      string          `json:"message"`
	Data         json.RawMessage `json:"data"` // bisa kirim object, array, dsb.
}

type APIREQDTO struct {
	TxType  string            `json:"txType"`
	Headers map[string]string `json:"headers"`
	Params  map[string]string `json:"params"`
}
