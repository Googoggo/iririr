package db

import (
	"awesomeProject3/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

func OnUserFileUploadFinished(username,filehash,filename string,filesize int64) bool{
	sql := "insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`,`upload_ad`) values (?,?,?,?,?)"
	stmt, err := mysql.DBConn().Prepare(sql)
	if err != nil{
		fmt.Sprintf(err.Error())
		return false
	}
	_, errstmt := stmt.Exec(username,filehash,filename,filesize,time.Now())
	if errstmt != nil{
		return false
	}
	return true
}

func QueryUserFileMetas(username string, limit int)([]UserFile, error){
	sql := "select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?"
	stmt, err := mysql.DBConn().Prepare(sql)
	if err != nil {
		fmt.Sprintf(err.Error())
		return nil,err
	}
	defer stmt.Close()

	rows, errrows := stmt.Query(username,limit)
	if errrows != nil{
		return nil, errrows
	}
	var userFiles []UserFile
	for rows.Next(){
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash,&ufile.FileName,&ufile.FileSize,&ufile.UploadAt,&ufile.LastUpdated)
		if err!= nil{
			fmt.Printf(err.Error())
			break
		}
		userFiles = append(userFiles,ufile)
	}
	return userFiles,nil
}