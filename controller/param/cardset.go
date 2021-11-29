package param

type NewCardsetRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Access      int    `json:"access"`
}

type UpdateCardsetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Access      int    `json:"access"`
}

type SearchCardsetRequest struct {
	PageRequest
	Keyword string `query:"keyword"`
}

type RandomCardsetsRequest struct {
	Count int `query:"count" validate:"required"`
}

type GetCardsetResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Access      int      `json:"access"`
	Cards       []string `json:"cards"`
}
