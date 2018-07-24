/**
 * Created by IDEA.
 * User: godsoul
 */
//解析request到结构体
package middleware

import (
	"io/ioutil"
	"github.com/guobin8205/api_demo/utils/helper/convert"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"log"
	"strings"
)

type Req struct {
	Accept string
	Ver    string
	Area   string
	OS     int
	UID    int
	Token  string

	ReqBody
}

func ParseRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &Req{
			Ver:    c.Request.Header.Get("Ver"),
			Area:   c.Request.Header.Get("Area"),
			OS:     convert.ToInt(c.Request.Header.Get("OS")),
			UID:    convert.ToInt(c.Request.Header.Get("UID")),
			Token:  c.Request.Header.Get("Token"),
		}
		log.Println(c.Request.Header)
		body, _ := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		//sbody, _ := base64.StdEncoding.DecodeString(string(body))
		log.Println(c.Request.RequestURI)
		if strings.Contains(c.Request.RequestURI,"debug/pprof"){
		//if c.Request.RequestURI == "/debug/pprof/" || c.Request.RequestURI == "/debug/pprof/heap" || c.Request.RequestURI == "/debug/pprof/goroutine" ||
		//	c.Request.RequestURI == "/debug/pprof/block" || c.Request.RequestURI == "/debug/pprof/threadcreate"{
			req.b = string(body)
			c.Set("req", req)
		}else{
			//dst := make([]byte, len(body))
			//dst,_ = crypto.NewAesCrypto().Decrypt(body)
			//req.b = string(dst)
			//log.Println(string(dst))
			////log.Println("======================")
			//c.Set("req", req)
			req.b = string(body)
			c.Set("req", req)
		}

	}
}

type ReqBody struct {
	b string
}

func (req *ReqBody) String(key string) string {
	s := gjson.Get(req.b, key).String()
	if s == "null" {
		return ""
	}
	return s
}

func (req *ReqBody) DefaultString(key string, def string) string {
	v := gjson.Get(req.b, key).String()
	if v == "null" {
		return def
	}
	return v
}

func (req *ReqBody) Int(key string) int {
	return int(gjson.Get(req.b, key).Int())
}

func (req *ReqBody) Int64(key string) int64 {
	return gjson.Get(req.b, key).Int()
}

func (req *ReqBody) DefaultInt(key string, def int) int {
	v := gjson.Get(req.b, key).Int()
	if v == 0 {
		return def
	}
	return int(v)
}

func (req *ReqBody) GetBody() string {
	return req.b
}
