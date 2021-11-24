package common

import (
	"net/http"
	"strings"
)

type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

//拦截器结构体
type Filter struct {
	//用来存储需要拦截的URI
	filtermap map[string]FilterHandle
}
//filter初始化函数
func NewFilter ()*Filter{
	return &Filter{filtermap: make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(uri string ,handler FilterHandle){
	f.filtermap[uri] = handler
}

//根据uri获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle{
	return f.filtermap[uri]
}

//声明新的函数类型
type WebHandle func(rw http.ResponseWriter,req *http.Request)

//执行拦截器，返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request){
	return func(rw http.ResponseWriter, r *http.Request){
		for path, handle := range f.filtermap{
			if strings.Contains(r.RequestURI,path){
				//执行拦截业务逻辑
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		//执行正常注册的函数
		webHandle(rw,r)
	}
}
