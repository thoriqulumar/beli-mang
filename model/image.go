package model

type PostImageResponse struct {
	Message string `json:"message"`
	Data    struct {
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
}
