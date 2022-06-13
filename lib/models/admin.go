package models

type Admin struct {
	Base
	Name        string     `json:"name" gorm:"unique;not null"`
	Username 	string 	   `json:"username" gorm:"unique;not null"`
	Password  	string     `json:"password" gorm:"not null"`
	Notices		[]Notice   `json:"-"`
}

func (adminLogin *Login) GetAdmin()  (Admin, error){
	admin := Admin{}
	err := GetDB().Where("username = ?",adminLogin.Username).First(&admin).Error
	return admin,err
}

func (adminLogin *Login) GetStudent()  (Student, error){
	student := Student{}
	err := GetDB().Where("username = ?",adminLogin.Username).First(&student).Error
	return student,err
}

type Login struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type ResAdminLogin struct {
	Admin
	Token string `json:"token"`
}