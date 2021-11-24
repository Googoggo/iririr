package main

import (
	"awesomeProject3/handler"
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	http.HandleFunc("/user/signup",handler.SignupHandler)
	http.HandleFunc("/user/signin",handler.SignInHandler)
	http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))
	http.HandleFunc("/file/fastupload",handler.HTTPInterceptor(handler.TryFastUploadHandler))

	http.HandleFunc("/file/mpupload/init",handler.HTTPInterceptor(handler.InitalMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",handler.HTTPInterceptor(handler.CompleteUploadHandler))
	err := http.ListenAndServe(":8080",nil)
	if err != nil{
		fmt.Printf("Failed to start err:%s", err.Error())
	}
}
