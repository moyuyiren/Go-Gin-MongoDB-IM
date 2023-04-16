package server_core

import (
	"GoIm/models"
	"GoIm/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func ChatList(c *gin.Context) {
	roomIdentity := c.Query("room_identity")
	if roomIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间号不能为空",
		})
		return
	}
	// 判断用户是否属于该房间
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	_, err := models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "非法访问",
		})
		return
	}
	//分页信息
	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	skip := (pageIndex - 1) * pageSize
	//查询聊天信息
	data, err := models.GetMessageListByRoomIdentity(roomIdentity, &pageSize, &skip)
	if err != nil {
		log.Println("[DB ERROR] 数据查询异常" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询异常，请稍后重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功" + roomIdentity,
		"data": data,
	})
}
