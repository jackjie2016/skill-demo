package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespJson(w http.ResponseWriter, data interface{}) {
	header := w.Header()
	header.Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	ret, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(ret)
}

type H struct {
	Code int
	Data interface{}
	Msg  string
}

// 当操作成功返回Ok,
func RespOk2(w http.ResponseWriter, code int, data interface{}) {
	RespJson(w, H{Code: code, Data: data})
}

// 当操作成功返回Ok,
func RespOk(w http.ResponseWriter, code int, data interface{}) {
	RespJson(w, H{Code: http.StatusOK, Data: data})
}

// 当操作失败返回Error,
func RespFail(w http.ResponseWriter, msg string) {
	RespJson(w, H{Code: http.StatusNotFound, Msg: msg})
}
