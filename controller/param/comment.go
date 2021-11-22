package param

type NewCommentRequest struct {
	Comment string `json:"comment" validate:"required"`
}

type CommentResponseItem struct {
	ID      string       `json:"id"`
	Owner   UserResponse `json:"owner"`
	Comment string       `json:"comment"`
}

type GetCommentResponse = []CommentResponseItem
