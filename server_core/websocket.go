package server_core

import (
	"GoIm/define"
	"GoIm/models"
	"GoIm/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{}

// var wc = make(map[string]*websocket.Conn)
var wcc sync.Map

// WebSocketMessage 发送消息
func WebSocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统内部异常" + err.Error(),
		})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_Claims").(*utils.UserClaims)
	//wc[uc.Identity] = conn
	//并发
	wcc.Store(uc.Identity, conn)
	for {
		ms := new(define.MessageStruct)
		err := conn.ReadJSON(ms)
		if err != nil {
			log.Printf("Read Error:%v\n", err)
			return
		}
		//TODO: 判断用户是否属于消息体的房间
		_, err = models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, ms.RoomIdentity)
		if err != nil {
			log.Printf("UserIdentity:%v RoomIdentity:%v Not Exits\n", uc.Identity, ms.RoomIdentity)
			return
		}
		//TODO: 保存消息
		mb := &models.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		err = models.InsertOneMessageBasic(mb)
		if err != nil {
			log.Printf("[DB Error]:%v\n", err)
			return
		}

		//TODO: 获取在特定房间的在线用户
		userRooms, err := models.GetUserRoomByRoomIdentity(ms.RoomIdentity)
		if err != nil {
			log.Printf("[DB Error]:%v\n", err)
			return
		}
		//
		//for _, room := range userRooms {
		//	if cc, ok := wc[room.UserIdentity]; ok {
		//		err := cc.WriteMessage(websocket.TextMessage, []byte(ms.Message))
		//		if err != nil {
		//			log.Printf("Write Message Error:%v\n", err)
		//			return
		//		}
		//	}
		//}
		//并发
		for _, room := range userRooms {
			if cc, ok := wcc.Load(room.UserIdentity); ok {
				err := cc.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					log.Printf("Write Message Error:%v\n", err)
					return
				}
			}
		}

	}
}
