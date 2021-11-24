package services

import (
	"awesomeProject1/datamodels"
	"awesomeProject1/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string,pwd string)(user *datamodels.User,isOk bool)
	AddUser(user *datamodels.User)(userId int64 ,err error)
}

type UserService struct {
	userReposiotory repositories.IUserRepository
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}


func ValidatePassword(userPassword string, hashed string)(isOk bool, err error){
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword));err!=nil{
		return false, errors.New("密码比对错误")
	}
	return true, nil
}

func (u *UserService)IsPwdSuccess(userName string, pwd string)(user *datamodels.User, isOk bool){
	var err error
	user, err =u.userReposiotory.Select(userName)
	if err != nil{
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOk{
		return &datamodels.User{}, false
	}
	return
}

func GeneratePassword(userPassword string)([]byte,error){
	return bcrypt.GenerateFromPassword([]byte(userPassword),bcrypt.DefaultCost)
}

func (u *UserService)AddUser(user *datamodels.User)(userId int64, err error){
	pwdByte ,errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil{
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.userReposiotory.Insert(user)
}