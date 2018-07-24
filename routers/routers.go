/**
 * Created by IDEA.
 * User: godsoul
 */
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/guobin8205/api_demo/controllers"
)

func SetRouters(e *gin.Engine) {
	e.GET("/user/login",controllers.Login)			//用户登录
	e.GET("/user/push",controllers.UserPush)			//直播弹幕
}