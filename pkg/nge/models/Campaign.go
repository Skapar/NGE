package models

import (
	"time"

	"gorm.io/gorm"
)

// Campaign struct represents a campaign with its details
type Campaign struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Owners      []*User   `gorm:"many2many:campaign_owners;" json:"owners"`
}

// CreateCampaign creates a new campaign
func CreateCampaign(db *gorm.DB, name string, description string, startDate time.Time, endDate time.Time, owners []*User) error {
	campaign := Campaign{
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Owners:      owners,
	}
	return db.Create(&campaign).Error
}

// DeleteCampaign deletes a campaign by ID
func DeleteCampaign(db *gorm.DB, campaignID uint) error {
	result := db.Delete(&Campaign{}, campaignID)
	return result.Error
}

// UpdateCampaign updates a campaign's details
func UpdateCampaign(db *gorm.DB, campaignID uint, name string, description string, startDate time.Time, endDate time.Time, owners []*User) error {
	result := db.Model(&Campaign{}).Where("id = ?", campaignID).Updates(Campaign{Name: name, Description: description, StartDate: startDate, EndDate: endDate, Owners: owners})
	return result.Error
}

// GetCampaignByID retrieves a campaign by ID
func GetCampaignByID(db *gorm.DB, campaignID uint) (Campaign, error) {
	var campaign Campaign
	result := db.Preload("Owners").First(&campaign, campaignID)
	return campaign, result.Error
}

// GetAllCampaigns retrieves all campaigns
func GetAllCampaigns(db *gorm.DB) ([]Campaign, error) {
	var campaigns []Campaign
	result := db.Preload("Owners").Find(&campaigns)
	return campaigns, result.Error
}

// GetCampaignsByUserID retrieves campaigns owned by a user
func GetCampaignsByUserID(db *gorm.DB, userID uint) ([]Campaign, error) {
	var campaigns []Campaign
	result := db.Preload("Owners").Joins("JOIN campaign_owners ON campaigns.id = campaign_owners.campaign_id").Where("campaign_owners.user_id = ?", userID).Find(&campaigns)
	return campaigns, result.Error
}
