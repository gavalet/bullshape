package ctrls

import (
	"bullshape/models"
	u "bullshape/utils"
	l "bullshape/utils/logger"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log := l.NewLogger("")
	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		log.Error("Could not unmarshal req Body. Error: ", err)
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}

	resp, status, err := models.CreateUser(account)
	if err != nil {
		log.Error("Create user returned error. Error: ", err)
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, resp)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	log := l.NewLogger("")
	account := &models.User{}
	body, _ := ioutil.ReadAll(r.Body)
	log.Debug("Auth Body  :\n %s", body)

	err := json.Unmarshal(body, &account)
	if err != nil {
		log.Error("Could not unmarshal req Body. Error: ", err)
		u.HttpError(w, http.StatusInternalServerError, errors.New("Error req body json."))
		return

	}

	resp, cookie, status, err := models.Login(account.Username, account.Password)
	if err != nil {
		log.Error("Login in failed. Error: ", err)
		u.HttpError(w, status, err)
	}
	http.SetCookie(w, cookie)
	u.HttpRespond(w, status, resp)
}
