package ctrls

import (
	"bullshape/models"
	u "bullshape/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.HttpError(w, http.StatusInternalServerError, err)
		return
	}

	resp, status, err := models.CreateUser(account) //Create account
	if err != nil {
		u.HttpError(w, status, err)
	}
	u.HttpRespond(w, status, resp)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("Auth Body  :\n %s", body)

	err := json.Unmarshal(body, &account)
	if err != nil {
		fmt.Println(err)
		u.HttpError(w, http.StatusInternalServerError, errors.New("Error req body json."))
		return

	}

	fmt.Println("lalal Auth:" + account.Username + " " + account.Password + "......")
	fmt.Println("")
	var cookie *http.Cookie
	resp, cookie, status, err := models.Login(account.Username, account.Password)
	if err != nil {
		u.HttpError(w, status, err)
	}
	fmt.Println("To cookie einai!")
	fmt.Print(cookie)
	http.SetCookie(w, cookie)
	// u.HttpRespond(w, status, resp)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
