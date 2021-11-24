package db

import (
	"awesomeProject3/db/mysql"
	"database/sql"
	"fmt"
)

func OnFileUploadFinished(filehash string,filename string,
	filesize int64,fileaddr string)bool{
	sql := "Insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`)"+"`file_addr`,`status` values(?,?,?,?,1)"
	stmt, err := mysql.DBConn().Prepare(sql)
	if err != nil{
		fmt.Printf("preparewrong")
		return false
	}
	defer stmt.Close()

	res ,errstmt := stmt.Exec(filehash,filename,filesize,fileaddr)
	if errstmt != nil{
		fmt.Println(errstmt.Error())
		return false
	}
	if rf, err := res.RowsAffected();nil == nil{
		if err != nil{
			fmt.Printf(err.Error())
		}
		if rf<0{
			fmt.Printf("hash:%s",filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string)(*TableFile,error){
	sql := "select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1? and status=1"
	stmt, err := mysql.DBConn().Prepare(sql)
	if err != nil{
		fmt.Printf(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash,&tfile.FileAddr,&tfile.FileName,&tfile.FileSize)
	if err != nil{
		fmt.Printf(err.Error())
		return nil,err
	}
	return &tfile,err
}