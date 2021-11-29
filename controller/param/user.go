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

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type LoginRequest struct {
	Mail     string `json:"mail" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Mail     string `json:"mail" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Avatar   string   `json:"avatar"`
	Favorite []string `json:"favorite"`
}

type AddCollectionRequest struct {
	CardsetID string `json:"cardset_id" validate:"required"`
	Liked     bool   `json:"liked"`
}

type CardsetInfoResponse struct {
	ID          string `json:"id"`
	OwnerID     string `json:"owner_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Access      int    `json:"access,omitempty"`
}
