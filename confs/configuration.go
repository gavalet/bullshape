package confs

import (
	"bullshape/utils"
	"fmt"

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

var TokenPass = "token_pass"

var Conf Configuration

func init() {
	pwd, _ := utils.GetPWD()
	confFile := pwd + "/bullshape-api.conf"
	config, err := toml.LoadFile(confFile)
	if err != nil {
		fmt.Println(err)
		return
		// os.Exit(1)
	}

	if expiration, ok := config.Get("expiration_cookie").(int64); ok {
		Conf.ExpirationCookie = int(expiration)
	} else {
		Conf.ExpirationCookie = 3600
	}

	if port, ok := config.Get("port").(int64); ok {
		Conf.ServerPort = port
	} else {
		Conf.ServerPort = 8080
	}

	if mySQLUser, ok := config.Get("mysql_user").(string); ok {
		Conf.MySQLUser = mySQLUser
	} else {
		Conf.MySQLUser = "root"
	}

	if mySQLPassword, ok := config.Get("mysql_password").(string); ok {
		Conf.MySQLPassword = mySQLPassword
	} else {
		Conf.MySQLPassword = ""
	}

	if mySQLDatabase, ok := config.Get("mysql_database").(string); ok {
		Conf.MySQLDatabase = mySQLDatabase
	} else {
		Conf.MySQLDatabase = "bullshape"
	}

	if defaultUsername, ok := config.Get("default_user").(string); ok {
		Conf.DefaultUsername = defaultUsername
	} else {
		Conf.DefaultUsername = "admin"
	}

	if defaultPassword, ok := config.Get("default_pwd").(string); ok {
		Conf.DefaultPassword = defaultPassword
	} else {
		Conf.DefaultPassword = "admin"
	}
}
