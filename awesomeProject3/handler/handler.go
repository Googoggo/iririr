package handler

import (
	"awesomeProject3/db"
	"awesomeProject3/meta"
	"awesomeProject3/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter,r *http.Request){
	if r.Method == "GET"{
		//返回html页面
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w,"innerserver")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == "POST"{
		//接受文件存到本地
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/"+head.Filename,
			UploadAt: time.Now().Format("2006-06-02 13:59:59"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err !=nil{
			fmt.Printf("Failed to create file,err:%s\n",err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize,err = io.Copy(newFile,file)
		if err != nil{
			fmt.Printf("Failed to save data into file,err:%s",err.Error())
			return
		}
		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		_=meta.UpdataFileMetaDB(fileMeta)
		//更新用户文件表
		r.ParseForm()
		username := r.Form.Get("username")
		suc := db.OnUserFileUploadFinished(username, fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
		if suc{
			http.Redirect(w,r,"/file/upload/suc", http.StatusFound)
		}else{
			w.Write([]byte("Upload Failed"))
		}
	}
}
func UploadSucHandler(w http.ResponseWriter, r *http.Request){
	io.WriteString(w,"Upload finished")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMetaDB(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f,err:= os.Open(fm.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err :=ioutil.ReadAll(f)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("content-dispostion","attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

func FileMetaUpdateHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	meta.UpdataFileMeta(curFileMeta)
	data, err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filesha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filesha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(filesha1)
	w.WriteHeader(http.StatusOK)
}

func FileQueryHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	limitCnt,_ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMeta,_ := meta.GetLastFileMetaDB(limitCnt)
	userFiles, err := db.QueryUserFileMetas(username,limitCnt)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize,_:= strconv.Atoi(r.Form.Get("filesize"))

	fileMeta := meta.GetFileMetaDB(filehash)

	if fileMeta == struct {
		FileSha1 string
		FileName string
		FileSize int64
		Location string
		UploadAt string
	}{}{
		resp := util.RespMsg{
			Code: -1,
			Msg: "秒传失败",
		}
		w.Write(resp.JSONBytes())
		return
	}
	suc := db.OnUserFileUploadFinished(username,filehash,filename,int64(filesize))
	if suc{
		resp := util.RespMsg{
			Code: 0,
			Msg: "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	}else{
		resp := util.RespMsg{
			Code: -2,
			Msg: "秒传失败",
		}
		w.Write(resp.JSONBytes())
		return
	}
}