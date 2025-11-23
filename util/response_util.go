package util

type response struct {
	Payload interface{} `json:"payload"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
}

func GenerateResponse(payload interface{}, status bool, message string) response {
	return response{
		Payload: payload,
		Status:  status,
		Message: message,
	}
}
