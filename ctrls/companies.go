package ctrls

import (
	"bullshape/models"
	"bullshape/utils"
	u "bullshape/utils"
	"bullshape/utils/kafkala"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ctrlServices struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
	Kafka  *kafkala.CompaniesProducer
}

func NewCtrlServices(logger *zap.SugaredLogger, db *gorm.DB, kafka *kafkala.CompaniesProducer) *ctrlServices {
	ctrlSrves := &ctrlServices{
		DB:     db,
		Logger: logger,
		Kafka:  kafka,
	}
	return ctrlSrves
}

func (ctrl *ctrlServices) GetCompany(w http.ResponseWriter, r *http.Request) {

	log := ctrl.Logger.With("REQ ID:", utils.GetReqId(r))
	id, err := u.GetUintParam(r, "id")
	if err != nil {
		log.Error("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	m := models.NewCtrlServices(log, ctrl.DB)
	data, status, err := m.GetCompany(id)
	if err != nil {
		log.Error("Got error: ", err)
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, data)
}

func (ctrl *ctrlServices) CreateCompany(w http.ResponseWriter, r *http.Request) {
	log := ctrl.Logger.With("REQ ID:", utils.GetReqId(r))
	productsProducer := kafkala.NewCompanyProducer()
	productsProducer.Run()
	defer productsProducer.Close()

	company := &models.NewCompany{}
	err := json.NewDecoder(r.Body).Decode(company)
	if err != nil {
		log.Error("Could not unmarshal req body. Error:", err)
		u.HttpError(w, http.StatusBadRequest, err)
	}
	m := models.NewCtrlServices(log, ctrl.DB)
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

func (ctrl *ctrlServices) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	log := ctrl.Logger.With("REQ ID:", utils.GetReqId(r))
	id, err := u.GetUintParam(r, "id")
	if err != nil {
		log.Error("Can not find company ID")
		//err := u.NewError(nil, GET_COMPANY_COMPANY_ID_ERRCODE, err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	m := models.NewCtrlServices(log, ctrl.DB)
	status, err := m.DeleteCompany(id)
	if err != nil {
		u.HttpError(w, status, err)
		log.Error("Delete company returned error: ", err)
	}
	u.HttpRespond(w, status, nil)
}

func (ctrl *ctrlServices) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	log := ctrl.Logger.With("REQ ID:", utils.GetReqId(r))
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
	m := models.NewCtrlServices(log, ctrl.DB)
	compOpts := &models.EditCompanyOpts{}
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
