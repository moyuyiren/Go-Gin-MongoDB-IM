package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type UserRoom struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int    `bson:"room_type"` //1.独聊 2.群聊
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

// GetUserRoomByUserIdentityRoomIdentity 查询用户与消息房间是否绑定
func GetUserRoomByUserIdentityRoomIdentity(useridentity, roomidentity string) (*UserRoom, error) {
	ur := new(UserRoom)
	err := Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(),
			bson.D{{"user_identity", useridentity}, {"room_identity", roomidentity}}).
		Decode(&ur)
	return ur, err
}

// GetUserRoomByRoomIdentity 获取房间是否存在
func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{"room_identity", roomIdentity}})
	if err != nil {
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ub := new(UserRoom)
		err := cursor.Decode(&ub)
		if err != nil {
			return nil, err
		}
		urs = append(urs, ub)
	}
	return urs, nil

}

// GetJudgeUserIsFriend 是否有房间
func GetJudgeUserIsFriend(userIdentity1, userIdentity2 string) (bool, error) {
	//查询 userIdentity1 单聊房间列表
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{
		{"user_identiyt", userIdentity1},
		{"room_type", 1},
	})
	roomIdentitys := make([]string, 0)
	if err != nil {
		log.Printf("[DB ERROR]:%v", err.Error())
		return false, nil
	}
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err := cursor.Decode(&ur)
		if err != nil {
			return false, nil
		}
		roomIdentitys = append(roomIdentitys, ur.RoomIdentity)
	}
	//获取关联
	cnt, err := Mongo.Collection(UserRoom{}.CollectionName()).CountDocuments(context.Background(), bson.D{
		{"user_identiyt", userIdentity2},
		{"root_identity", bson.M{"$in": roomIdentitys}},
	})
	if err != nil {
		log.Printf("[DB ERROR]:%v", err.Error())
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

// InsertOneUserRoom 添加用户好友
func InsertOneUserRoom(ur *UserRoom) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).
		InsertOne(context.Background(), ur)
	return err
}

// GetUserRoomIdentity 获取房间号
func GetUserRoomIdentity(userIdentity1, userIdentity2 string) string {
	//查询 userIdentity1 单聊房间列表
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{
		{"user_identiyt", userIdentity1},
		{"room_type", 1},
	})
	roomIdentitys := make([]string, 0)
	if err != nil {
		log.Printf("[DB ERROR]:%v", err.Error())
		return ""
	}
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err := cursor.Decode(&ur)
		if err != nil {
			return ""
		}
		roomIdentitys = append(roomIdentitys, ur.RoomIdentity)
	}
	//获取关联
	ur := new(UserRoom)
	err = Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{"user_identiyt", userIdentity2},
			{"room_identity", bson.M{"$in": roomIdentitys}},
		}).Decode(ur)
	if err != nil {
		log.Printf("[DB ERROR]:%v\n", err)
		return ""
	}
	return ur.RoomIdentity
}

// DeleteUserRoom 删除房间号
func DeleteUserRoom(roomIdentity string) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).
		DeleteOne(context.Background(), bson.M{"room_identity": roomIdentity})
	if err != nil {
		log.Printf("[DB ERROR]:%v\n", err)
		return err
	}
	return nil
}
