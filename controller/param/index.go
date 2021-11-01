package param

type VerifyCodeRequest struct {
	Mail string `json:"mail" validate:"required"`
}
