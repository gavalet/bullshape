package confs

import (
	"bullshape/utils"
	"log"
	"os"

	toml "github.com/pelletier/go-toml"
)

type Configuration struct {
	ExpirationCookie int // In seconds

	ServerPort    int64
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string

	DefaultUsername string
	DefaultPassword string
}

var Conf Configuration

func init() {
	pwd, _ := utils.GetPWD()
	confFile := pwd + "/bullshape-api.conf"
	config, err := toml.LoadFile(confFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var ok bool
	expiration, ok := config.Get("expiration_cookie").(int64)
	if ok {
		Conf.ExpirationCookie = int(expiration)
	} else {
		Conf.ExpirationCookie = 3600
	}

	port, ok := config.Get("port").(int64)
	if ok {
		Conf.ServerPort = port
	} else {
		Conf.ServerPort = 8999
	}

	mySQLUser, ok := config.Get("mysql_user").(string)
	if ok {
		Conf.MySQLUser = mySQLUser
	} else {
		Conf.MySQLUser = "root"
	}

	mySQLPassword, ok := config.Get("mysql_password").(string)
	if ok {
		Conf.MySQLPassword = mySQLPassword
	} else {
		Conf.MySQLPassword = ""
	}
	mySQLDatabase, ok := config.Get("mysql_database").(string)
	if ok {
		Conf.MySQLDatabase = mySQLDatabase
	} else {
		Conf.MySQLDatabase = ""
	}

	DefaultUsername, ok := config.Get("default_user").(string)
	if ok {
		Conf.DefaultUsername = DefaultUsername
	} else {
		Conf.DefaultUsername = "admin"
	}

	DefaultPassword, ok := config.Get("default_pwd").(string)
	if ok {
		Conf.DefaultPassword = DefaultPassword
	} else {
		Conf.DefaultPassword = "admin"
	}
}
