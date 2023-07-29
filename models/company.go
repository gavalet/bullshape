package models

import (
	"bullshape/db"
	"bullshape/utils"
	u "bullshape/utils"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Company struct {
	ID          uint   `json:"id"`
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	NumEmployes uint   `json:"num_of_employes"`
	Registered  bool   `json:"registered"`
	Type        string `json:"type"`
}

type NewCompany struct {
	Name        *string `json:"name"`
	UUID        *string `json:"uuid"`
	Description string  `json:"description"`
	NumEmployes *uint   `json:"num_of_employes"`
	Registered  *bool   `json:"registered"`
	Type        *string `json:"type"`
}

type EditCompanyOpts struct {
	Description *string `json:"description"`
	NumEmployes *uint   `json:"num_of_employes"`
	Registered  *bool   `json:"registered"`
	Type        *string `json:"type"`
}

const (
	CORPORATION        = "corporation"
	NONPROFIT          = "non profit"
	COOPERATIVE        = "cooperative"
	SOLEPROPRIETORSHIP = "sole proprietorship"
)

// companiesMX stores locks for each company.
var companiesMX = sync.Map{}

func GetCompany(id uint) (*Company, int, error) {
	dbCompany, err := db.GetCompanyByID(db.GormDB, id)
	if err != nil {
		return nil, http.StatusNotFound, err //u.NewError(err, GET_COMPANY_ERRCODE, errors.New(GET_COMPANY_ERR))
	}
	company := serializeCompany(dbCompany)
	return &company, http.StatusOK, nil
}

func CreateCompany(newCompany NewCompany) (*Company, int, error) {

	if newCompany.UUID == nil {
		uuid := utils.NewUUIDV4()
		newCompany.UUID = &uuid
		fmt.Println("Create new Company UUID")
	}
	if newCompany.Name == nil || newCompany.NumEmployes == nil ||
		newCompany.Type == nil || newCompany.Registered == nil {
		return nil, http.StatusBadRequest,
			u.NewError(nil, CREATE_COMPANY_EMPTY_PARAMS_ERRCODE, errors.New(CREATE_COMPANY_EMPTY_PARAMS_ERR))
	}
	_, err := db.GetCompanyByName(db.GormDB, *newCompany.Name)
	if err == nil {
		return nil, http.StatusBadRequest,
			u.NewError(nil, CREATE_COMPANY_COMPANY_NAME_ERRCODE, errors.New(CREATE_COMPANY_COMPANY_NAME_ERR))
	}
	fmt.Println("Create company")

	dbCompany, err := createDBCompanyObj(newCompany)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = dbCompany.Create(db.GormDB)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	cmpn := serializeCompany(dbCompany)
	return &cmpn, http.StatusOK, nil
}

func DeleteCompany(id uint) (int, error) {
	companiesMX, _ := companiesMX.LoadOrStore(fmt.Sprint(id), &sync.Mutex{})
	companiesMX.(*sync.Mutex).Lock()
	defer companiesMX.(*sync.Mutex).Unlock()

	dbCompany, err := db.GetCompanyByID(db.GormDB, id)
	if err != nil {
		return http.StatusNotFound, err //u.NewError(err, GET_COMPANY_ERRCODE, errors.New(GET_COMPANY_ERR))
	}
	err = dbCompany.Delete(db.GormDB)
	if err != nil {
		return http.StatusInternalServerError, err //u.NewError(err, GET_COMPANY_ERRCODE, errors.New(GET_COMPANY_ERR))
	}
	return http.StatusNoContent, nil
}

func UpdateCompany(id uint, opt EditCompanyOpts) (*Company, int, error) {
	companiesMX, _ := companiesMX.LoadOrStore(fmt.Sprint(id), &sync.Mutex{})
	companiesMX.(*sync.Mutex).Lock()
	defer companiesMX.(*sync.Mutex).Unlock()

	dbCompany, err := db.GetCompanyByID(db.GormDB, id)
	if err != nil {
		return nil, http.StatusNotFound, err //u.NewError(err, GET_COMPANY_ERRCODE, errors.New(GET_COMPANY_ERR))
	}

	if opt.Description != nil {
		dbCompany.Description = *opt.Description
	}
	if opt.NumEmployes != nil {
		dbCompany.NumEmployes = *opt.NumEmployes
	}
	if opt.Registered != nil {
		dbCompany.Registered = *opt.Registered
	}
	if opt.Type != nil && isValidType(*opt.Type) {
		dbCompany.Type = *opt.Type
	}
	err = dbCompany.Update(db.GormDB)
	if err != nil {
		return nil, http.StatusInternalServerError, err //u.NewError(err, GET_COMPANY_ERRCODE, errors.New(GET_COMPANY_ERR))
	}
	cmpn := serializeCompany(dbCompany)
	return &cmpn, http.StatusOK, nil
}

func createDBCompanyObj(newCompany NewCompany) (*db.Company, error) {
	dbCompany := db.Company{}
	dbCompany.UUID = *newCompany.UUID
	dbCompany.Name = *newCompany.Name
	if !isValidType(*newCompany.Type) {
		return nil, errors.New("Type: " + *newCompany.Type + "is not valid.")
	}
	dbCompany.Type = *newCompany.Type
	dbCompany.Description = newCompany.Description
	dbCompany.NumEmployes = *newCompany.NumEmployes
	dbCompany.Registered = *newCompany.Registered
	return &dbCompany, nil
}

func serializeCompany(dbCompany *db.Company) Company {
	company := Company{}
	company.ID = dbCompany.ID
	company.UUID = dbCompany.UUID
	company.Description = dbCompany.Description
	company.NumEmployes = dbCompany.NumEmployes
	company.Registered = dbCompany.Registered
	company.Type = dbCompany.Type
	return company
}

func isValidType(t string) bool {
	switch t {
	case COOPERATIVE, NONPROFIT, SOLEPROPRIETORSHIP, CORPORATION:
		return true
	}
	return false
}
