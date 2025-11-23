package dto

type StoryCreateDTO struct {
	Data []struct {
		Content  string `form:"content" json:"content"`
		Duration string `form:"duration" json:"duration"`
	} `json:"data" form:"data" binding:"dive"`
}

type StoryUserSeen struct {
	Ids []uint64 `json:"ids" binding:"required"`
}
