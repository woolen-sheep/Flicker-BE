package param

type NewCardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Image    string `json:"image"`
	Audio    string `json:"audio"`
}

type UpdateCardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Image    string `json:"image"`
	Audio    string `json:"audio"`
}

type GetCardResponse struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Image    string `json:"image"`
	Audio    string `json:"audio"`
}

type CardCommentRequest struct {
	Comment string `json:"comment"`
}
