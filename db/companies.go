package db

import (
	"bullshape/utils"

	"github.com/jinzhu/gorm"
)

const (
	CORPORATION        = "corporation"
	NONPROFIT          = "non profit"
	COOPERATIVE        = "cooperative"
	SOLEPROPRIETORSHIP = "sole proprietorship"
)

type Company struct {
	gorm.Model
	UUID        string `gorm:"unique_index"`
	Description string
	NumEmployes uint
	Registered  bool
	Type        string
}

func (company *Company) Create(tx *gorm.DB) error {
	if len(company.UUID) == 0 {
		company.UUID = utils.NewUUIDV4()
	}
	//todo check if it is unique
	return tx.Create(&company).Error
}

func (company *Company) Update(tx *gorm.DB) error {
	return tx.Update(&company).Error
}

func (company *Company) Delete(tx *gorm.DB) error {
	return tx.Delete(&company).Error
}

func GetCompanyByID(tx *gorm.DB, id uint) []Company {
	company := []Company{}
	tx = tx.Table("company").Where("id = ? ", id).Find(&company)
	return company
}

func GetCompanyByUUID(tx *gorm.DB, uuid string) []Company {
	company := []Company{}
	tx = tx.Table("company").Where("image = ? ", uuid).Find(&company)
	return company
}
