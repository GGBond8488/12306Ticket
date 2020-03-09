package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

//响应请求的结构体
type ResponseBody struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Resp(w http.ResponseWriter, code int, msg string, data interface{}) {
	resp(w, code, msg, data)
}
func resp(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ResponseBody{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	//将结构体转为json字符串返回
	ret, err := json.Marshal(response)
	if err != nil {
		log.Panicln(err.Error())
	}
	w.Write(ret)
}
