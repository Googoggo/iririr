package handler

import (
	rPool "awesomeProject3/cache/redis"
	"awesomeProject3/db"
	"awesomeProject3/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//MultipartUploadInfo:初始化
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadId string
	ChunkSize int
	ChunkCount int
}
//初始化分块上传
func InitalMultipartUploadHandler(w http.ResponseWriter, r *http.Request){
	//解析用户参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil{
		w.Write(util.NewRespMsg(-1,"params invalid",nil).JSONBytes())
		return
	}
	//获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash: filehash,
		FileSize: filesize,
		UploadId: username + fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize: 5*1024*1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize)/(5*1024*1024))),
	}

	//将初始化信息写入redis缓存
	rConn.Do("HSET","MP_"+upInfo.UploadId,"chunkcount",upInfo.ChunkCount)
	rConn.Do("HSET","MP_"+upInfo.UploadId,"filehash",upInfo.FileHash)
	rConn.Do("HSET","MP_"+upInfo.UploadId,"filesize",upInfo.FileSize)

	//将响应初始化数据返回到客户端
	w.Write(util.NewRespMsg(0,"OK",upInfo).JSONBytes())
}

//上传文件分块
func UploadPartHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//获得的文件句柄，用于存储分块内容
	fpath := "/data/"+uploadID+"/"+chunkIndex
	os.MkdirAll((fpath),0744)
	fd, err := os.Create(fpath)  
	if err != nil {
		w.Write(util.NewRespMsg(-1,"Upload part failed",nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte,1024*1024)
	for{
		n,err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}
	//更新redis缓存状态
	rConn.Do("HSET","MP_"+uploadID,"chkidx_"+chunkIndex,1)

	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}

//通知合并上传
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request){
	//解析参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	//获得redis连接池的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL","MP_"+upid))
	if err != nil{
		w.Write(util.NewRespMsg(-1,"complete upload failed",nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	//模糊
	for i:=0;i<len(data);i+=2{
		k:=string(data[i].([]byte))
		v:=string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount,_ = strconv.Atoi(v)
		}else if strings.HasPrefix(k,"chidx_")&&v == "1"{
			chunkCount += 1
		}
	}
	if totalCount != chunkCount{
		w.Write(util.NewRespMsg(-2,"invalid request",nil).JSONBytes())
		return
	}

	//合并分块

	//更新唯一文件表及用户文件表
	fsize,_ := strconv.Atoi(filesize)
	db.OnUserFileUploadFinished(filename,filehash,filename,int64(fsize))
	db.OnUserFileUploadFinished(username,filehash,filename,int64(fsize))

	//响应处理结果
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}
