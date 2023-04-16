package define

import "time"

var MailPassword = "KXUHFYPJHOFVHVRE"

type MessageStruct struct {
	Message      string `json:"message"`
	RoomIdentity string `json:"room_identity"`
}

var RegisterPrefix = "Token_"
var ExpireTime = time.Second * 300
