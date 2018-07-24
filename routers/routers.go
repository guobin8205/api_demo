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
	e.POST("/user/login",controllers.Login)			//用户登录
	e.POST("/user/push",controllers.UserPush)			//直播弹幕
}