package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	// Id       uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Nickname string `json:"nickname"`
	//SubscriberProfilemap []Profile `"gorm:"many2many:profile_subscriber_map;"`
	// Followers       []Profile `json:"followers"`
	// ProfilePostsMap []Post `gorm:"many2many:profile_user_map;"`
}

func AddProfile(db *gorm.DB, profile *Profile) error {
	result := db.Create(profile)
	return result.Error
}

func GetProfileById(db *gorm.DB, id uint) (*Profile, error) {
	var profile Profile
	result := db.First(&profile, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

func UpdateProfileById(db *gorm.DB, id uint, updatedProfile *Profile) error {
	var profile Profile
	result := db.First(&profile, id)
	if result.Error != nil {
		return result.Error
	}
	profile.Nickname = updatedProfile.Nickname

	result = db.Save(&profile)
	return result.Error
}

func DeleteProfileById(db *gorm.DB, id uint) error {
	result := db.Delete(&Profile{}, id)
	return result.Error
}
