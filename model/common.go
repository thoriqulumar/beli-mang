package model

type GeneralResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}
