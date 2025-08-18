package response

type DefaultResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
