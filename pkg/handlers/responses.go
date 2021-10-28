package handlers

type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type authBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type deleteBody struct {
	FileName string `json:"file_name"`
	Version  string `json:"version"`
}

type authResponse struct {
	Token string `json:"token"`
}
