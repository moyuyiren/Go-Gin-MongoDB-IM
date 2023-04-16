package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserBasic struct {
	Identity  string `bson:"identity"`
	Account   string `bson:"account"`
	Password  string `bson:"password"`
	Nickname  string `bson:"nickname"`
	Sex       int    `bson:"sex"`
	Email     string `bson:"email"`
	Avatar    string `bson:"avatar"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

func (UserBasic) Collection() string {
	return "user_basic"
}

// GetUserBasicByAccountPassword 查询登录密码
func GetUserBasicByAccountPassword(account, password string) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.Collection()).
		FindOne(context.Background(), bson.D{{"account", account}, {"password", password}}).
		Decode(&ub)
	return ub, err
}

// GetUserBasicByIdentity 查询登录用户信息
func GetUserBasicByIdentity(identity string) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.Collection()).
		FindOne(context.Background(), bson.D{{"identity", identity}}).
		Decode(&ub)
	return ub, err
}

// GetUserBasicCountByEmail 查询用户email是否重复注册
func GetUserBasicCountByEmail(email string) (int64, error) {
	return Mongo.Collection(UserBasic{}.Collection()).
		CountDocuments(context.Background(), bson.D{{"email", email}})
}

// GetUserBasicCountByAccount 查询用户email是否重复注册
func GetUserBasicCountByAccount(account string) (int64, error) {
	return Mongo.Collection(UserBasic{}.Collection()).
		CountDocuments(context.Background(), bson.D{{"account", account}})
}

// InsertOneUserBasic 保存用户
func InsertOneUserBasic(ub *UserBasic) error {
	_, err := Mongo.Collection(UserBasic{}.Collection()).InsertOne(context.Background(), ub)
	return err
}

// QueryUserByAccount 查询用户信息
func QueryUserByAccount(account string) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.Collection()).
		FindOne(context.TODO(),
			bson.D{
				{"account", account},
			}).Decode(&ub)
	if err != nil {
		return nil, err
	}
	return ub, err

}
