package models

import (
	"bullshape/confs"
	"bullshape/db"
	"bullshape/utils"
	u "bullshape/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}

func serialiseUser(dbUser *db.User, token string) User {
	user := User{}
	user.Username = dbUser.Username
	user.Token = token
	return user
}

//validateUser validates incoming user details...
func validateUser(user *User) error {

	if user.Username == "" {
		return errors.New("Username should not be empty")
	}
	if len(user.Password) < 6 {
		return errors.New("Password should not at least 7 characters")
	}
	err, _ := db.GetUserByUsername(db.GormDB, user.Username)
	if err == nil {
		return errors.New("Username already in use by another user")
	}
	return nil
}

func CreateUser(user *User) (*User, int, error) {

	err := validateUser(user)
	if err != nil {
		return nil, http.StatusBadRequest,
			u.NewError(err, CREATE_USER_WRONG_PARAMS_ERRCODE,
				errors.New(CREATE_USER_WRONG_PARAMS_ERR))
	}

	hashedPassword := utils.EncryptPass(user.Password)
	userDB := db.User{Username: user.Username, Password: hashedPassword, Token: user.Token}

	if err := userDB.Create(db.GormDB); err != nil {
		//log.Error("Failed to create user.")
		return nil, CREATE_USER_DB_WRITE_FAIL_ERRCODE,
			u.NewError(err, CREATE_USER_WRONG_PARAMS_ERRCODE,
				errors.New("Failed to create use."))
	}

	expirationTime := time.Now().Add(time.Duration(confs.Conf.ExpirationCookie) * time.Minute)

	tk := &Token{UserId: userDB.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	userNew := serialiseUser(&userDB, tokenString)
	return &userNew, http.StatusOK, nil
}

func Login(username, password string) (*User, *http.Cookie, int, error) {

	err, user := db.GetUserByUsername(db.GormDB, username)
	if err != nil {
		//log.Error("Can not find user")
		return nil, nil, http.StatusBadRequest, errors.New("Can not find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, nil, http.StatusUnauthorized, errors.New("Invalid //login credentials")
	}
	user.Password = ""
	expirationTime := time.Now().Add(time.Duration(confs.Conf.ExpirationCookie) * time.Minute)

	tk := &Token{UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(confs.TokenPass))
	user.Token = tokenString
	usr := serialiseUser(user, tokenString)

	return &usr, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	}, http.StatusOK, nil
}

func GetUser(id uint) (*User, int, error) {

	err, userDB := db.GetUserByID(db.GormDB, id)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("Could not found user with ID:" + fmt.Sprint(id))
	}
	user := serialiseUser(userDB, "")
	return &user, http.StatusOK, nil
}
