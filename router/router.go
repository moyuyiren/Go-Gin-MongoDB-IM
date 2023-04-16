package router

import (
	"GoIm/middlewars"
	"GoIm/server_core"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	//用户登录
	r.POST("/login", server_core.Login)
	r.POST("/login", server_core.Login)
	//用户注册
	r.POST("/login", server_core.Register)
	//发送验证码
	r.POST("/SendCode", server_core.SendCode)

	//登录后逻辑，需要Token=====================================================================================
	auth := r.Group("/u", middlewars.AuthChechk())
	//用户详情
	auth.GET("/user/detail", server_core.UserDetailIdentidy)
	//查询指定用户的个人信息
	auth.GET("/user/query", server_core.Query)
	//发送，接收消息
	auth.GET("/websocket/message", server_core.WebSocketMessage)
	//聊天记录列表
	auth.GET("/chat/list", server_core.ChatList)
	//添加好友
	auth.POST("/user/adduser", server_core.AddFriend)
	//删除好友
	auth.POST("/user/delete", server_core.DeleteFriend)

	return r
}
