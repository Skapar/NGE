package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	// ID   uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Text string `json:"text" gorm:"unique"`
}

type DBModel struct {
	DB *gorm.DB
}

func (dbm *DBModel) AddPost(post Post) (Post, error) {
	err := dbm.DB.Create(&post).Error
	return post, err
}

func (db *DBModel) GetPost(id uint) (Post, error) {
	var post Post
	err := db.DB.First(&post, id).Error
	return post, err
}

func (db *DBModel) UpdatePost(updatedPost Post) (Post, error) {
	err := db.DB.Save(&updatedPost).Error
	return updatedPost, err
}

func (db *DBModel) DeletePost(id uint) error {
	var post Post
	err := db.DB.Delete(&post, id).Error
	return err
}
