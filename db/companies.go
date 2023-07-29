package db

import (
	"github.com/jinzhu/gorm"
)

type Company struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	UUID        string `gorm:"unique"`
	Description string
	NumEmployes uint   `gorm:"not null"`
	Registered  bool   `gorm:"not null"`
	Type        string `gorm:"not null"`
}

func (company *Company) Create(tx *gorm.DB) error {
	return tx.Create(&company).Error
}

func (company *Company) Update(tx *gorm.DB) error {
	return tx.Save(&company).Error
}

func (company *Company) Delete(tx *gorm.DB) error {
	return tx.Delete(&company).Error
}

func GetCompanies(tx *gorm.DB, id uint) (*Company, error) {
	company := Company{}
	err := tx.Where("id = ? ", id).Find(&company).Error
	return &company, err
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
