package httpresponse

type RestError struct {
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"causes,omitempty"`
}

type RestSuccess struct {
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
