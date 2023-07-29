package db

import "github.com/jinzhu/gorm"

//User is struct for bullshape's system users
type User struct {
	gorm.Model
	Username string `gorm:"index"`
	Password string
	Token    string
}

//Create creates a new user in DB
func (user *User) Create(tx *gorm.DB) error {
	return tx.Create(&user).Error
}

//Update updates the user in DB
func (user *User) Update(tx *gorm.DB) error {
	return tx.Update(&user).Error
}

//Delete removes a user from the DB
func (user *User) Delete(tx *gorm.DB) error {
	return tx.Delete(&user).Error
}

//GetUserByID finds and returns a user from DB based on its ID
func GetUserByID(tx *gorm.DB, id uint) (error, *User) {
	user := new(User)
	err := tx.Where("ID = ? ", id).Find(&user).Error
	return err, user
}

//GetUserByID finds and returns a user from DB based on its Username
func GetUserByUsername(tx *gorm.DB, username string) (error, *User) {
	user := new(User)
	err := tx.Where("username = ? ", username).Find(&user).Error
	return err, user
}
