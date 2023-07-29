package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ErrorJson struct {
	ErrorCode int `json:"code"`
	// A human-readable message, describing the error.
	Message string `json:"message"`
	Err     error  `json:"errors"`
}

func NewError(errs error, errorCode int, message error) error {
	err := ErrorJson{}
	err.ErrorCode = errorCode
	err.Err = errs
	err.Message = message.Error()
	return &err
}

func (e *ErrorJson) Error() string {
	var errStr string
	errStr = fmt.Sprintln(errStr, "")
	if e.Err != nil {
		errStr = fmt.Sprintln(errStr, "Error:", e.Err)
	}

	errStr = fmt.Sprintln(errStr, "GENERAL CODE:", e.ErrorCode)
	errStr = fmt.Sprintln(errStr, "GENERAL MESSAGE:", e.Message)

	return errStr
}

func GetUintParam(r *http.Request, param string) (uint, error) {
	sid, ok := mux.Vars(r)[param]
	if !ok {
		return 0, errors.New("The parameter " + param + " does not exist")
	}
	id, err := strconv.ParseUint(sid, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func HttpError(w http.ResponseWriter, status int, str_errors ...error) {
	// Write the HTTP Status on header.
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(str_errors) > 0 {
		serrs := []string{}
		for _, v := range str_errors {
			if v != nil {
				serrs = append(serrs, v.Error())
			}
		}
		errs := map[string]interface{}{
			"errors": serrs,
			"result": "FAILURE",
		}
		b, _ := json.Marshal(errs)
		w.Write(b)
	}
}

func HttpRespond(w http.ResponseWriter, header int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(header)
	json.NewEncoder(w).Encode(data)
}
