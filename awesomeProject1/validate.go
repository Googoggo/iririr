package main

import (
	"awesomeProject1/common"
	"awesomeProject1/datamodels"
	"awesomeProject1/encrypt"
	"awesomeProject1/rabbitmq"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)
//设置集群地址，最好内外ip
var hostArray = []string{"127.0.0.1","127.0.0.1"}
//
var localHost = ""

var port = "8081"

//数量控制接口服务器内网ip，或者slb内网ip
var GetOneIp = "127.0.0.1"

//
var GetOnePort = "8084"

var hashConsistent *common.Consistent

var rabbitMqValidate *rabbitmq.RabbitMQ
//用来存放控制信息
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourceArray map[int]interface{}
	*sync.RWMutex
}
//存储数据
var accessControl = &AccessControl{sourceArray: make(map[int]interface{})}

func (m *AccessControl)GetNewRecord(uid int) interface{}{
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int){
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = "hello lpx"
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool{
	uid ,err := req.Cookie("uid")
	if err != nil {
		return false
	}
	//采取一致性hash算法，根据用户id，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil{
		return false
	}
	//判断是否是本机
	if hostRequest == localHost{
		//执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	}else {
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}

//获取其他节点处理结果
func (m *AccessControl) GetDataFromOtherMap(host string ,request *http.Request) bool{
	hostUrl := "http://"+host+":"+port+"/"
	response,body,err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true"{
			return true
		}else {
			return false
		}
	}
	return false
}

//获取本机map， 并且处理业务逻辑，返回结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string)(isOk bool){
	uidInt ,err := strconv.Atoi(uid)
	if err != nil{
		return false
	}
	data := m.GetNewRecord(uidInt)
	if data != nil {
		return true
	}
	return false
}

//统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error{
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil{
		return err
	}
	return nil
}

//身份校验
func CheckUserInfo(r *http.Request) error{
	//获取uid，cookie
	uidCookie, err:= r.Cookie("uid")
	if err != nil{
		return errors.New("用户uid的cookie获取失败")
	}
	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil{
		return errors.New("用户加密串cookie获取失败")
	}
	//对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil{
		return errors.New("加密串已被篡改")
	}
	fmt.Println("结果比对")
	fmt.Println(uidCookie.Value)
	fmt.Println(string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)){
		return nil
	}
	return errors.New("身份校验失败")
}

func checkInfo(checkStr string ,signStr string) bool{
	if checkStr == signStr{
		return true
	}
	return false
}

//模拟请求
func GetCurl(hostUrl string,request *http.Request)(response *http.Response, body []byte, err error){
	//获取uid
	uidPre ,err := request.Cookie("uid")
	if err != nil{
		return
	}
	//获取sign
	uidSign ,err :=request.Cookie("sign")
	if err != nil {
		return
	}
	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET",hostUrl,nil)
	if err != nil {
		return
	}
	//手动指定，排查多余cookies
	cookieUid := &http.Cookie{Name: "uid",Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign",Value: uidSign.Value,Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	//获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil{
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request){
	right := accessControl.GetDistributedRight(r)
	if !right{
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

func Check(w http.ResponseWriter, r *http.Request){
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"])<0  {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)
	//获取uidCookie
	userCookie, err := r.Cookie("uid")
	if err != nil{
		w.Write([]byte("false"))
		return
	}
	//1.分布式验证权限
	right := accessControl.GetDistributedRight(r)
	if right == false{
		w.Write([]byte("false"))
	}

	//2.获取数量控制权限，防止秒杀超卖现象
	hostUrl := "http://"+GetOneIp+":"+GetOnePort+"/getOne"
	responserValidate, validateBody ,err := GetCurl(hostUrl,r)
	if err != nil{
		w.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态
	if responserValidate.StatusCode == 200{
		if string(validateBody) == "true"{
			//整合下单
			productID, err := strconv.ParseInt(productString,10,64)
			if err != nil{
				w.Write([]byte("false"))
				return
			}
			userID, err := strconv.ParseInt(userCookie.Value,20,64)
			if err != nil{
				w.Write([]byte("false"))
				return
			}
			//创建消息体
			message := &datamodels.Message{userID, productID}
			//类型转化
			byteMessage , err := json.Marshal(message)
			if err != nil{
				w.Write([]byte("false"))
				return
			}
			//生产消息
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}

	}
	w.Write([]byte("false"))
	return
}

func main(){
	localIp ,err := common.GetIntranceIp()
	if err != nil{
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("lpxlpx")
	defer  rabbitMqValidate.Destory()
	//负载均衡过滤器
	hashConsistent = common.NewConsistent()
	//采用一致性hash算法添加节点
	for _, v := range hostArray{
		hashConsistent.Add(v)
	}
	//1.过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check",Auth)
	filter.RegisterFilterUri("/checkRight",Auth)
	//启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight",filter.Handle(CheckRight))
	http.ListenAndServe(":8003",nil)
}
