package repositories

import (
	"awesomeProject1/common"
	"awesomeProject1/datamodels"
	"database/sql"
	"errors"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string)(user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
	SelectById(userId int64)(user *datamodels.User,err error)
}

type UserManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func NewUserRepostiory(table string, db *sql.DB) IUserRepository{
	return &UserManagerRepository{table, db}
}

func (u *UserManagerRepository) Conn() (err error){
	if u.mysqlConn == nil{
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil{
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == ""{
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string)(user *datamodels.User, err error){
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空")
	}
	if err = u.Conn(); err!= nil {
		return &datamodels.User{}, err
	}

	sql1 := "Select * from "+u.table+" where userName=?"
	rows, errRows := u.mysqlConn.Query(sql1, userName)
	defer rows.Close()
	if errRows != nil {
		return &datamodels.User{}, errRows
	}
	result := common.GetResultRow(rows)
	if len(result) == 0{
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User)(userId int64, err error){
	if err = u.Conn(); err != nil{
		return
	}
	sql1 := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql1)
	if errStmt != nil {
		print("zheli")
		return userId, errStmt
	}
	result, errResult := stmt.Exec(user.NickName,user.UserName,user.HashPassword)
	if errResult != nil {
		print("nali")
		return userId ,errResult
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectById(userId int64)(user *datamodels.User, err error){
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql1 := "select *From "+u.table+" where ID="+strconv.FormatInt(userId,10)
	row, errRow := u.mysqlConn.Query(sql1)
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("不对劲啊")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result ,user)
	return
}