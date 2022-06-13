package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"lib/utils"
	"net/http"
	"strconv"
)

var notAuth = []string{"/", "/student/login", "/admin/login","/token"}

var Auth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message := new(utils.Message)
		message.SetLanguage(r)
		requestPath := r.URL.Path
		w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")             //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		token := r.Header.Get("token")
		if token == ""{
			msg := message.GetMsg(utils.MsgToken, "no_token")
			http.Error(w, fmt.Sprintf(msg), http.StatusUnauthorized)
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil{
			msg := message.GetMsg(utils.MsgToken, "token")
			http.Error(w, fmt.Sprintf(msg,err), http.StatusUnauthorized)
			return
		}
		r.Header.Add("id", strconv.Itoa(int(claims.(jwt.MapClaims)["id"].(float64))))
		//r.Header.Add("admin", string(claims.(jwt.MapClaims)["admin"].(bool)))
		next.ServeHTTP(w, r)
	})
}
