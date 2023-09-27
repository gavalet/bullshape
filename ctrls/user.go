package ctrls

import (
	"bullshape/models"
	"bullshape/utils"
	u "bullshape/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func (ctrl *ctrlServices) CreateUser(w http.ResponseWriter, r *http.Request) {
	log := ctrl.Logger.With("REQ ID:", utils.GetReqId(r))
	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		log.Error("Could not unmarshal req Body. Error: ", err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}
	m := models.NewCtrlServices(log, ctrl.DB)
	resp, status, err := m.CreateUser(account)
	if err != nil {
		log.Error("Create user returned error. Error: ", err)
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, resp)
}

func (ctrl *ctrlServices) Authenticate(w http.ResponseWriter, r *http.Request) {
	log := ctrl.Logger.With("REQ ID:", utils.GetReqId(r))
	account := &models.User{}
	body, _ := ioutil.ReadAll(r.Body)
	log.Debug("Auth Body  :\n %s", body)

	err := json.Unmarshal(body, &account)
	if err != nil {
		log.Error("Could not unmarshal req Body. Error: ", err)
		u.HttpError(w, http.StatusInternalServerError, errors.New("Error req body json."))
		return

	}
	m := models.NewCtrlServices(log, ctrl.DB)
	resp, cookie, status, err := m.Login(account.Username, account.Password)
	if err != nil {
		log.Error("Login in failed. Error: ", err)
		u.HttpError(w, status, err)
	}
	http.SetCookie(w, cookie)
	u.HttpRespond(w, status, resp)
}
