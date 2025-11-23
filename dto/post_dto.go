package dto

type PostCreateDTO struct {
	Content      string   `form:"content"`
	Hashtags     []string `form:"hashtags" binding:"dive,number"`
	Mentions     []string `form:"mentions" binding:"dive,number"`
	SharedPostId uint64   `form:"shared_post_id"`
}

type PostUpdateDTO struct {
	Content string `form:"content"`
}
