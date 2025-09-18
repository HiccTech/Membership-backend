package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Pet struct {
	BaseModel
	Phone                 string         `json:"phone"`
	PetAvatarUrl          string         `json:"petAvatarUrl"`
	PetName               string         `json:"petName"`
	PetType               string         `json:"petType"`
	Breed                 string         `json:"breed"`
	Birthday              string         `json:"birthday"`
	Weight                string         `json:"weight"`
	Gender                string         `json:"gender"`
	VaccinationRecords    *int           `json:"vaccinationRecords"`
	SterilizationStatus   *int           `json:"sterilizationStatus"`
	HasMedicalCondition   *bool          `json:"hasMedicalCondition"`
	MedicalConditionMap   datatypes.JSON `json:"medicalConditionMap"`
	MedicalConditionOther string         `json:"medicalConditionOther"`
	CoatType              string         `json:"coatType"`
	GroomingFrequency     string         `json:"groomingFrequency"`
	ShopifyCustomerId     string         `gorm:"size:64;not null;index" json:"shopifyCustomerId"`
}

// 创建表
func MigratePet(db *gorm.DB) error {
	return db.AutoMigrate(&Pet{})
}
