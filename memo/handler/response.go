package handler

// Response represents http response.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
