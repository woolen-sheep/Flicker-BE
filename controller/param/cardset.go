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
	ID            string   `json:"id"`
	OwnerID       string   `json:"owner_id"`
	OwnerName     string   `json:"owner_name"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Access        int      `json:"access"`
	FavoriteCount int      `json:"favorite_count"`
	VisitCount    int      `json:"visit_count"`
	CreateTime    int64    `json:"create_time"`
	Cards         []string `json:"cards"`
	IsFavorite    bool     `json:"is_favorite"`
	LastStudy     int64    `json:"last_study"`
}
