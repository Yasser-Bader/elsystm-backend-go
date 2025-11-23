package models

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   Media  `json:"avatar" gorm:"foreignKey:ModelId"`
}

type UserNameAndUsername struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserWithStories struct {
	ID       uint64                             `json:"id" gorm:"primaryKey"`
	Name     string                             `json:"name"`
	Username string                             `json:"username"`
	Email    string                             `json:"email"`
	SeenAll  bool                               `json:"seen_all"`
	Avatar   Media                              `json:"avatar" gorm:"foreignKey:ModelId"`
	Stories  []SocialStoryWithoutSeenSerializer `json:"stories" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}

func (UserWithStories) TableName() string {
	return "users"
}

func (UserNameAndUsername) TableName() string {
	return "users"
}
