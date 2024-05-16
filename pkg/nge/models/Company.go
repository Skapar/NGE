package models

import (
	"time"

	"gorm.io/gorm"
)

// Company struct represents a company with its details
type Company struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Owners      []*User   `gorm:"many2many:company_owners;" json:"owners"`
}

// CreateCompany creates a new company
func CreateCompany(db *gorm.DB, name string, description string, startDate time.Time, endDate time.Time, owners []*User) error {
	company := Company{
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Owners:      owners,
	}
	return db.Create(&company).Error
}

// DeleteCompany deletes a company by ID
func DeleteCompany(db *gorm.DB, companyID uint) error {
	result := db.Delete(&Company{}, companyID)
	return result.Error
}

// UpdateCompany updates a company's details
func UpdateCompany(db *gorm.DB, companyID uint, name string, description string, startDate time.Time, endDate time.Time, owners []*User) error {
	result := db.Model(&Company{}).Where("id = ?", companyID).Updates(Company{Name: name, Description: description, StartDate: startDate, EndDate: endDate, Owners: owners})
	return result.Error
}

// GetCompanyByID retrieves a company by ID
func GetCompanyByID(db *gorm.DB, companyID uint) (Company, error) {
	var company Company
	result := db.Preload("Owners").First(&company, companyID)
	return company, result.Error
}

// GetAllCompanies retrieves all companies
func GetAllCompanies(db *gorm.DB) ([]Company, error) {
	var companies []Company
	result := db.Preload("Owners").Find(&companies)
	return companies, result.Error
}

// GetCompaniesByUserID retrieves companies owned by a user
func GetCompaniesByUserID(db *gorm.DB, userID uint) ([]Company, error) {
	var companies []Company
	result := db.Preload("Owners").Joins("JOIN company_owners ON companies.id = company_owners.company_id").Where("company_owners.user_id = ?", userID).Find(&companies)
	return companies, result.Error
}
