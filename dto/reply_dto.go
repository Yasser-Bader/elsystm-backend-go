package dto

type ReplyCreateDTO struct {
	Content   string   `form:"content"`
	PostId    string   `form:"post_id" binding:"required,number"`
	CommentId string   `form:"comment_id" binding:"required,number"`
	Hashtags  []string `form:"hashtags" binding:"dive,number"`
	Mentions  []string `form:"mentions" binding:"dive,number"`
}

type ReplyUpdateDTO struct {
	Content string `form:"content"`
}
