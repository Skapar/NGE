package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Id              uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Subscribers     []uint `json:"subscribers" gorm:"-"`
	Followers       []uint `json:"followers"`
	ProfilePostsMap string `json:"profilePostsMap"`
}

type ProfileSubscriberRequest struct {
	SubscriberID uint64 `json:"subscriber_id"`
}

func AddProfile(db *gorm.DB, profile *Profile) error {
	if err := db.Create(profile).Error; err != nil {
		return err
	}
	return nil
}

func GetProfileById(db *gorm.DB, id uint) (*Profile, error) {
	var profile Profile
	if err := db.First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func UpdateProfileById(db *gorm.DB, id uint, updatedProfile *Profile) error {
	return db.Model(Profile{}).Where("id = ?", id).Updates(updatedProfile).Error
}

func DeleteProfileById(db *gorm.DB, id uint) error {
	return db.Delete(&Profile{}, id).Error
}

func AddSubcriberById(db *gorm.DB, id uint, subscriberID uint) error {
	var profile Profile
	if err := db.First(&profile, id).Error; err != nil {
		return err
	}
	profile.Subscribers = append(profile.Subscribers, subscriberID)
	if err := db.Save(&profile).Error; err != nil {
		return err
	}

	return nil
}

func GetSubcriberById(db *gorm.DB, id uint) ([]uint, error) {
	var profile Profile
	if err := db.First(&profile, id).Error; err != nil {
		return nil, err
	}
	return profile.Subscribers, nil
}

func GetFollowersById(db *gorm.DB, id uint) ([]uint, error) {
	var profile Profile
	if err := db.First(&profile, id).Error; err != nil {
		return nil, err
	}
	return profile.Followers, nil
}

func DeleteFollowersById(db *gorm.DB, id uint, followerID uint) error {
	var profile Profile
	if err := db.First(&profile, id).Error; err != nil {
		return err
	}
	index := -1
	for i, follower := range profile.Followers {
		if follower == followerID {
			index = i
			break
		}
	}
	if index != -1 {
		profile.Followers = append(profile.Followers[:index], profile.Followers[index+1:]...)
	}
	if err := db.Save(&profile).Error; err != nil {
		return err
	}

	return nil
}

func GetUsersPosts(db *gorm.DB, userID uint) ([]Post, error) {
	var posts []Post
	if err := db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func AddSubscriberById(db *gorm.DB, id uint, subscriberID uint) error {
	var profile Profile
	if err := db.First(&profile, id).Error; err != nil {
		return err
	}

	for _, existingSubscriberID := range profile.Subscribers {
		if existingSubscriberID == subscriberID {
			return nil
		}
	}

	profile.Subscribers = append(profile.Subscribers, subscriberID)

	if err := db.Save(&profile).Error; err != nil {
		return err
	}

	return nil
}
