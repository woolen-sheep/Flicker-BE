package param

type VerifyCodeRequest struct {
	Mail string `json:"mail" validate:"required"`
}

type PageRequest struct {
	Skip  int `query:"skip"`
	Limit int `query:"limit"`
}

var (
	DefaultPageRequest = PageRequest{
		Skip:  0,
		Limit: 10,
	}
)
