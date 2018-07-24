/**
 * Created by IDEA.
 * User: godsoul
 */
package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/guobin8205/api_demo/middleware"
)
//获取参数
func GetReqParas(c *gin.Context) *middleware.Req {
	p, _ := c.Get("req")
	return p.(*middleware.Req)
}
//返回json封装
func jsonReturn(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg, "data": data})
}

