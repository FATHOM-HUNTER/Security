package handlers

type successResponse struct {
	Success bool `json:"success"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type authBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
