package handlers

type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type authBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}
