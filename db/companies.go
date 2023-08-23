package db

import (
	"github.com/jinzhu/gorm"
)

const (
	NONPROFIT          = CType("NonProfit")
	COOPERATIVE        = CType("Cooperative")
	CORPORATIONS       = CType("(Corporations")
	SOLEPROPRIETORSHIP = CType("Solo Proprietorship")
)

type CType string

type Company struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	UUID        string `gorm:"unique"`
	Description string `gorm:"type:varchar(3000)"`
	NumEmployes uint   `gorm:"not null"`
	Registered  bool   `gorm:"not null"`
	Type        CType  `gorm:"not null"`
}

func (company *Company) Create(tx *gorm.DB) error {
	return tx.Create(&company).Error
}

func (company *Company) Update(tx *gorm.DB) error {
	return tx.Save(&company).Error
}

func (company *Company) Delete(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&company).Error
}

func GetCompanies(tx *gorm.DB) ([]Company, error) {
	companies := []Company{}
	err := tx.Find(&companies).Error
	return companies, err
}

func GetCompanyByID(tx *gorm.DB, id uint) (*Company, error) {
	company := Company{}
	err := tx.Where("id = ? ", id).Find(&company).Error
	return &company, err
}

func GetCompanyByUUID(tx *gorm.DB, uuid string) (*Company, error) {
	company := Company{}
	err := tx.Where("uuid = ? ", uuid).Find(&company).Error
	return &company, err
}

func GetCompanyByName(tx *gorm.DB, name string) (*Company, error) {
	company := Company{}
	err := tx.Where("name = ? ", name).Find(&company).Error
	return &company, err
}
