/**
 * Created by IDEA.
 * User: godsoul
 */
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/guobin8205/api_demo/models"
	"github.com/guobin8205/api_demo/utils/helper/custerror"
	"time"
)

type call struct {
	data interface{}
	err  *custerror.Error
}
//测试下goroutine数量
//用户登录
func Login(c *gin.Context) {
	var (
		req      = GetReqParas(c)
		uuid      = req.String("uuid")
		account = req.String("account")
		ip = c.ClientIP()
	)
	user_model := &models.User{}
	user_model.Account = account
	user_model.IP = ip
	ch := make(chan call, 1)
	go func() {
		data, err := user_model.Login(uuid)
		ch <- call{data: data, err: err}
	}()
	select {
	case resp := <-ch:
		if resp.err != nil {
			jsonReturn(c, resp.err.Code, resp.data.(string), "")
		} else {
			jsonReturn(c, 0, "ok", resp.data)
		}
	case <-time.After(2 * time.Second):
		jsonReturn(c, 4003, "登录超时", "")
	}
}
//弹幕推送
func UserPush(c *gin.Context) {
	var (
		req      = GetReqParas(c)
		uuid      = req.String("uuid")
		account = req.String("account")
		content  = req.String("content")
		nickname  = req.String("nickname")
		roomid   = req.Int("roomid")
		ip = c.ClientIP()
	)

	if models.CheckUser(account, uuid) {
		u := models.GetUserByAccount(account)
		if u == nil{
			err := custerror.New(4005)
			jsonReturn(c,err.Code,err.Msg,"")
			return
		}
		data, err := u.UserPush(roomid, nickname,content,ip)
		if err != nil {
			jsonReturn(c, err.Code, data.(string), "")
		} else {
			jsonReturn(c, 0, "ok", "")
		}
	}else{
		err := custerror.New(4002)
		jsonReturn(c,err.Code,err.Msg,"")
	}
}
