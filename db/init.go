package db

import (
	"bullshape/confs"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var GormDB *gorm.DB
var AllTables = []interface{}{&User{}, &Company{}}

func init() {

	cmd := fmt.Sprintf(confs.Conf.MySQLUser + ":" + confs.Conf.MySQLPassword + "@/" +
		confs.Conf.MySQLDatabase + "?charset=utf8&parseTime=True&loc=Local")

	fmt.Println("db cmd command: ", cmd)
	db, err := gorm.Open("mysql", cmd)
	if err != nil {
		fmt.Println("Failed with error:", err)
	}
	db.Debug().AutoMigrate(AllTables...)
	GormDB = db
}
