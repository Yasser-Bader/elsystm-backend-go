package dto

type CommentCreateDTO struct {
	Content  string   `form:"content"`
	PostId   string   `form:"post_id" binding:"required,number"`
	Hashtags []string `form:"hashtags" binding:"dive,number"`
	Mentions []string `form:"mentions" binding:"dive,number"`
}

type CommentUpdateDTO struct {
	Content string `form:"content"`
}
