package param

type NewCommentRequest struct {
	Comment string `json:"comment" validate:"required"`
}

type CommentResponseItem struct {
	ID         string       `json:"id"`
	Owner      UserResponse `json:"owner"`
	Comment    string       `json:"comment"`
	LastUpdate string       `json:"lastupdate"`
	Liked      bool         `json:"liked"`
	Likes      int          `json:"likes"`
}

type GetCommentResponse = []CommentResponseItem

type LikeCommentRequest struct {
	Liked bool `json:"liked"`
}
