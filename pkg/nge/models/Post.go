package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	// Id   uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Text string `json:"text" gorm:"unique"`
}

func AddPost(db *gorm.DB, post Post) (Post, error) {
	err := db.Create(&post).Error
	return post, err
}

func GetPost(db *gorm.DB, id uint) (Post, error) {
	var post Post
	err := db.First(&post, id).Error
	return post, err
}

func UpdatePost(db *gorm.DB, updatedPost Post) (Post, error) {
	err := db.Save(&updatedPost).Error
	return updatedPost, err
}

func DeletePost(db *gorm.DB, id uint) error {
	err := db.Delete(&Post{}, id).Error
	return err
}

func GetAllPosts(db *gorm.DB) ([]Post, error) {
	var posts []Post
	err := db.Find(&posts).Error
	return posts, err
}
