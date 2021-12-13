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
	ID            string   `json:"id,omitempty"`
	OwnerID       string   `json:"owner_id,omitempty"`
	OwnerName     string   `json:"owner_name,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Access        int      `json:"access,omitempty"`
	FavoriteCount int      `json:"favorite_count,omitempty"`
	VisitCount    int      `json:"visit_count,omitempty"`
	CreateTime    int64    `json:"create_time,omitempty"`
	Cards         []string `json:"cards,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
}
