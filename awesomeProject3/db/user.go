package db

import (
	mydb"awesomeProject3/db/mysql"
	"fmt"
)
func UserSignup(username string, password string) bool{
	sql := "Insert ignore into tbl_user(`user_name`,`user_pwd`) values(?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err!=nil {
		fmt.Println("failed to insert"+err.Error())
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(username,password)
	if err != nil{
		fmt.Printf(err.Error())
		return false
	}
	if rowAffected, err := res.RowsAffected();nil == err && rowAffected>0{
		return true
	}
	return false

}
func UserSignin(username string, encpwd string) bool{
	stmt, err := mydb.DBConn().Prepare("select *from tbl_user where username=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil{
		fmt.Printf(err.Error())
		return false
	}else if rows == nil{
		fmt.Printf("notfound")
		return false
	}
	pRows := mydb.ParseRows(rows)
	if len(pRows)>0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd{
		return true
	}
	return false
}

func UpdateToken(username string, token string)bool{
	sql := "replace into tbl_user_token (`user_name`,`user_token`) values(?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	 defer stmt.Close()

	_, errstmt := stmt.Exec(username,token)
	if errstmt != nil {
		fmt.Printf(err.Error())
		return false
	}
	return true
}

type User struct {
	Username string
	Email string
	Phone string
	SignAt string
	LasActiveAt string
	Status int
}
func GetUserInfo(username string)(User, error){
	user := User{}
	sql := "select user_name,signup_at from tbl_user where user_name=? limit 1"
	stmt ,err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf(err.Error())
		return User{},nil
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username,&user.SignAt)
	if err != nil {
		return user,err
	}
	return user,nil
}
