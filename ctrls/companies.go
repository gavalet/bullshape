package ctrls

import (
	m "bullshape/models"
	u "bullshape/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCompany(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetUintParam(r, "id")
	if err != nil {
		fmt.Println("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	data, err := m.GetCompany(id)
	if err != nil {
		u.HttpError(w, http.StatusInternalServerError, err)
	}
	u.HttpRespond(w, http.StatusOK, data)
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {

	company := &m.NewCompany{}
	err := json.NewDecoder(r.Body).Decode(company)
	if err != nil {
		u.HttpError(w, http.StatusBadRequest, err)
	}
	data, status, err := m.CreateCompany(*company)
	if err != nil {
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, data)
}

func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetUintParam(r, "id")
	if err != nil {
		fmt.Println("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	status, err := m.DeleteCompany(id)
	if err != nil {
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, nil)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {

	id, err := u.GetUintParam(r, "id")
	if err != nil {
		fmt.Println("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	compOpts := &m.EditCompanyOpts{}
	err = json.NewDecoder(r.Body).Decode(compOpts)
	if err != nil {
		u.HttpError(w, http.StatusBadRequest, err)
	}
	company, status, err := m.UpdateCompany(id, *compOpts)
	if err != nil {
		u.HttpError(w, http.StatusInternalServerError, err)
	}
	u.HttpRespond(w, status, company)
}
