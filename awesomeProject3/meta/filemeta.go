package meta

import (
	mydb "awesomeProject3/db"
	"fmt"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}

//新增文件元信息
func UpdataFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1] = fmeta
}

func UpdataFileMetaDB(fmeta FileMeta) bool{
	return mydb.OnFileUploadFinished(fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.Location)
}

//获得文件元信息
func GetFileMeta(filesha1 string) FileMeta{
	return fileMetas[filesha1]
}

func GetFileMetaDB(filesha1 string) FileMeta{
	tfile, err:= mydb.GetFileMeta(filesha1)
	if err != nil{
		fmt.Println(err.Error())
		return FileMeta{}
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta
}
func RemoveFileMeta(filesha1 string)  {
	delete(fileMetas,filesha1)
}