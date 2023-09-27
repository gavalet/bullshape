package db

import (
	"bullshape/confs"
	"bullshape/utils"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// var GormDB *gorm.DB
var AllTables = []interface{}{&User{}, &Company{}}

func NewDatabaseConnection() *gorm.DB {

	dbhost := utils.EnvString("DB_HOST", "127.0.0.1")
	// log.Debug("db host is : ", dbhost)
	cmd := fmt.Sprintf(confs.Conf.MySQLUser + ":" + confs.Conf.MySQLPassword + "@tcp(" + dbhost + ":3306)/" +
		confs.Conf.MySQLDatabase + "?charset=utf8&parseTime=True&loc=Local")

	// log.Debug("db cmd command: ", cmd)
	db, err := gorm.Open("mysql", cmd)
	if err != nil {
		return nil
		//log.Error("Failed with error:", err)
	}
	db.Debug().AutoMigrate(AllTables...)
	return db
}
