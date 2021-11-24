package handler

import (
	"awesomeProject3/db"
	"awesomeProject3/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const pwd_salt = "*#123"
func SignupHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username)<3 || len(password)<5{
		w.Write([]byte("invalid param"))
		return
	}
	enc_password := util.Sha1([]byte(password+pwd_salt))
	suc := db.UserSignup(username, enc_password)
	if suc{
		w.Write([]byte("success"))
	}else{
		w.Write([]byte("failed"))
	}
}

func SignInHandler( w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password+pwd_salt))

	pwdChecked := db.UserSignin(username,encPasswd)
	if !pwdChecked{
		w.Write([]byte("FAILED"))
		return
	}

	token := GetToken(username)
	upRes := db.UpdateToken(username,token)
	if !upRes{
		w.Write([]byte("FAILED"))
		return
	}

	//w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "http://"+r.Host+"/static/view/home.html",
			Username: username,
			Token: token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GetToken(username string) string{
	ts := fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix + ts[:8]
}

func IsTokenValid(token string) bool{
	return true
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")

	isTokenValid := IsTokenValid(token)
	if !isTokenValid{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user, err := db.GetUserInfo(username)
	if err != nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}
