package models

import (
	"bullshape/db"
	"bullshape/utils"
	u "bullshape/utils"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Company struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	UUID        string   `json:"uuid"`
	Description string   `json:"description"`
	NumEmployes uint     `json:"num_of_employes"`
	Registered  bool     `json:"registered"`
	Type        db.CType `json:"type"`
}

type NewCompany struct {
	Name        *string   `json:"name,omitempty"`
	UUID        *string   `json:"uuid"`
	Description string    `json:"description"`
	NumEmployes *uint     `json:"num_of_employes"`
	Registered  *bool     `json:"registered"`
	Type        *db.CType `json:"type"`
}

type EditCompanyOpts struct {
	Description *string   `json:"description"`
	NumEmployes *uint     `json:"num_of_employes"`
	Registered  *bool     `json:"registered"`
	Type        *db.CType `json:"type"`
}

const (
	CORPORATION        = "corporation"
	NONPROFIT          = "non profit"
	COOPERATIVE        = "cooperative"
	SOLEPROPRIETORSHIP = "sole proprietorship"
)

// companiesMX stores locks for each company.
var companiesMX = sync.Map{}

type services struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}

func NewCtrlServices(logger *zap.SugaredLogger, db *gorm.DB) *services {
	ctrlSrves := &services{
		DB:     db,
		Logger: logger,
	}
	return ctrlSrves
}

func (srv *services) GetCompany(id uint) (*Company, int, error) {
	dbCompany, err := db.GetCompanyByID(srv.DB, id)
	if err != nil {
		srv.Logger.Error("Could not get Company. Error: ", err)
		return nil, http.StatusNotFound,
			u.NewError(err, GET_COMPANY_ERRCODE, errors.New(GET_COMPANY_ERR))
	}
	company := serializeCompany(dbCompany)
	return &company, http.StatusOK, nil
}

func (srv *services) CreateCompany(newCompany NewCompany) (*Company, int, error) {

	if newCompany.UUID == nil {
		uuid := utils.NewUUIDV4()
		newCompany.UUID = &uuid
		srv.Logger.Info("Create new Company UUID")
	}
	if newCompany.Name == nil || newCompany.NumEmployes == nil ||
		newCompany.Type == nil || newCompany.Registered == nil {
		srv.Logger.Error(CREATE_COMPANY_EMPTY_PARAMS_ERR)
		return nil, http.StatusBadRequest,
			u.NewError(nil, CREATE_COMPANY_EMPTY_PARAMS_ERRCODE,
				errors.New(CREATE_COMPANY_EMPTY_PARAMS_ERR))
	}
	dbCompany, err := createDBCompanyObj(newCompany)
	if err != nil {
		srv.Logger.Error("Failed to create company object")
		err := u.NewError(nil, CREATE_COMPANY_COMPANY_WRONG_TYPE_ERRCODE,
			errors.New(CREATE_COMPANY_COMPANY_NAME_ERR))
		return nil, http.StatusBadRequest, err
	}
	err = dbCompany.Create(srv.DB)
	if err != nil {
		srv.Logger.Error("Failed to create DB object : Error:", err)
		return nil, http.StatusInternalServerError, err
	}

	cmpn := serializeCompany(dbCompany)
	return &cmpn, http.StatusOK, nil
}

func (srv *services) DeleteCompany(id uint) (int, error) {
	companiesMX, _ := companiesMX.LoadOrStore(fmt.Sprint(id), &sync.Mutex{})
	companiesMX.(*sync.Mutex).Lock()
	defer companiesMX.(*sync.Mutex).Unlock()

	dbCompany, err := db.GetCompanyByID(srv.DB, id)
	if err != nil {
		srv.Logger.Error("Company does not include in DB.")
		return http.StatusNotFound,
			u.NewError(err, DELETE_COMPANY_NO_ENTRY_ERRCODE,
				errors.New(GET_COMPANY_ERR))
	}
	err = dbCompany.Delete(srv.DB)
	if err != nil {
		srv.Logger.Error("Unable to delete company from DB.")
		return http.StatusInternalServerError,
			u.NewError(err, DELETE_COMPANY_COULD_NOT_DELETE_ERRCODE,
				errors.New(GET_COMPANY_ERR))
	}
	return http.StatusNoContent, nil
}

func (srv *services) UpdateCompany(id uint, opt EditCompanyOpts) (*Company, int, error) {
	companiesMX, _ := companiesMX.LoadOrStore(fmt.Sprint(id), &sync.Mutex{})
	companiesMX.(*sync.Mutex).Lock()
	defer companiesMX.(*sync.Mutex).Unlock()

	dbCompany, err := db.GetCompanyByID(srv.DB, id)
	if err != nil {
		return nil, http.StatusNotFound,
			u.NewError(err, UPDATE_COMPANY_NO_ENTRY_ERRCODE,
				errors.New(GET_COMPANY_ERR))
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
	if opt.Type != nil && IsValidType(*opt.Type) {
		dbCompany.Type = *opt.Type
	}
	err = dbCompany.Update(srv.DB)
	if err != nil {
		return nil, http.StatusInternalServerError,
			u.NewError(err, UPDATE_COMPANY_COULD_NOT_UPDATE,
				errors.New(GET_COMPANY_ERR))
	}
	cmpn := serializeCompany(dbCompany)
	return &cmpn, http.StatusOK, nil
}

func createDBCompanyObj(newCompany NewCompany) (*db.Company, error) {
	dbCompany := db.Company{}
	dbCompany.UUID = *newCompany.UUID
	dbCompany.Name = *newCompany.Name
	if !IsValidType(*newCompany.Type) {
		return nil, errors.New("Type is not valid.")
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
	company.Name = dbCompany.Name
	company.UUID = dbCompany.UUID
	company.Description = dbCompany.Description
	company.NumEmployes = dbCompany.NumEmployes
	company.Registered = dbCompany.Registered
	company.Type = dbCompany.Type
	return company
}

func IsValidType(t db.CType) bool {
	switch t {
	case COOPERATIVE, NONPROFIT, SOLEPROPRIETORSHIP, CORPORATION:
		return true
	}
	return false
}
