package ctrls

import (
	"bullshape/models"
	m "bullshape/models"
	u "bullshape/utils"
	"bullshape/utils/kafkala"
	l "bullshape/utils/logger"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

func GetCompany(w http.ResponseWriter, r *http.Request) {
	log := l.NewLogger("")
	id, err := u.GetUintParam(r, "id")
	if err != nil {
		log.Error("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	data, status, err := m.GetCompany(id)
	if err != nil {
		log.Error("Got error: ", err)
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, data)
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {
	log := l.NewLogger("")
	productsProducer := kafkala.NewCompanyProducer()
	productsProducer.Run()
	defer productsProducer.Close()

	company := &m.NewCompany{}
	err := json.NewDecoder(r.Body).Decode(company)
	if err != nil {
		log.Error("Could not unmarshal req body. Error:", err)
		u.HttpError(w, http.StatusBadRequest, err)
	}
	data, status, err := m.CreateCompany(*company)
	if err != nil {
		log.Error("Create Company returned error. Error:", err)
		u.HttpError(w, status, err)
	}
	message, err := publishMessage(data)
	if err == nil {
		log.Debug("Send create event")
		productsProducer.PublishCreate(context.Background(), message)
	}
	u.HttpRespond(w, status, data)
}

func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	log := l.NewLogger("")
	id, err := u.GetUintParam(r, "id")
	if err != nil {
		log.Error("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	status, err := m.DeleteCompany(id)
	if err != nil {
		u.HttpError(w, status, err)
		log.Error("Delete company returned error: ", err)
	}
	u.HttpRespond(w, status, nil)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	log := l.NewLogger("")
	productsProducer := kafkala.NewCompanyProducer()
	productsProducer.Run()
	defer productsProducer.Close()

	id, err := u.GetUintParam(r, "id")
	if err != nil {
		log.Error("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	compOpts := &m.EditCompanyOpts{}
	err = json.NewDecoder(r.Body).Decode(compOpts)
	if err != nil {
		log.Error("Could not unmarshal body")
		u.HttpError(w, http.StatusBadRequest, err)
	}
	company, status, err := m.UpdateCompany(id, *compOpts)
	if err != nil {
		u.HttpError(w, status, err)
	}
	message, err := publishMessage(company)
	if err == nil {
		log.Debug("Send update event")
		productsProducer.PublishUpdate(context.Background(), message)
	}
	u.HttpRespond(w, status, company)
}

// PublishCreate create new product
func publishMessage(company *models.Company) (kafka.Message, error) {

	prodBytes, err := json.Marshal(&company)
	if err != nil {
		return kafka.Message{}, errors.New("Can not marshal")
	}
	return kafka.Message{
		Value: prodBytes,
		Time:  time.Now().UTC(),
	}, nil
}
