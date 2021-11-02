package param

type User struct {
	Mail     string `json:"mail"`
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type SignUpRequest struct {
	Mail     string `json:"mail" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type LoginRequest struct {
	Mail     string `json:"mail" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Mail     string `json:"mail" validate:"required"`
	Password string `json:"password" validate:"required"`
}
