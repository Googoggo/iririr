package mysql

import 				(
	"database/sql"
	"fmt"
	"os"
	_"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init(){
	db,err := sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/database3?charset=utf8")
	if err != nil{
		fmt.Println(err.Error())
	}
	db.SetMaxOpenConns(1000)
	errping := db.Ping()
	if errping != nil{
		fmt.Println(errping.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB{
	return db
}