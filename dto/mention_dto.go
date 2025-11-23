package dto

type MentionCreateDTO struct {
	UserMentionId uint64 `json:"user_mention_id"`
	ModelId       uint64 `json:"model_id"`
	ModelType     string `json:"model_type"`
}
