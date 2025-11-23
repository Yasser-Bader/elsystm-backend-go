package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Elsystm-Inc/systm-go-social/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var connectionInstance *gorm.DB

func initConnection() *gorm.DB {
	godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.SocialStoryMigration{}, &models.SocialStoryUser{},
		&models.SocialPostMigration{}, &models.SocialPostShare{}, &models.SocialMedia{}, &models.SocialCommentMigration{},
		&models.SocialReplyMigration{}, &models.SocialHistory{}, &models.SocialLike{}, &models.SocialMention{},
		&models.SocialHashtag{}, &models.SocialHashtagPCR{})

	return db
}

func ConnectDB() *gorm.DB {
	if connectionInstance == nil {
		connectionInstance = initConnection()
	} else {
		connectionInstance = initConnection()
	}

	return connectionInstance
}
