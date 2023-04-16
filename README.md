一个基于Gin+MongoDB的单机im系统

消息表
```json
{
    "_id": ObjectId("6438b250f5ed51e8040c1c03"),
    "user_identity": "用户的唯一标识",
    "room_identity": "房间的唯一标识",
    "data": "发送的数据",
    "created_at": 1,
    "updated_at": 1
}
```
房间表
```json
{
    "_id": ObjectId("6438babaf5ed51e8040c1c04"),
    "number": "房间号",
    "name": "房间名称",
    "info": "房间简介",
    "user_identity": "房间创建者的唯一标识",
    "created_at": 1,
    "updated_at": 1,
    "identity": "ad"
}
```
用户表
```json
{
    "_id": ObjectId("6438afc8f5ed51e8040c1c02"),
    "account": "账号",
    "password": "密码",
    "nickname": "昵称",
    "sex": 1,
    "email": "邮箱",
    "avatar": "头像",
    "create_at": 1,
    "updated_at": 1
}
```
用户消息房间关联
```json
{
    "_id": ObjectId("6438bb07f5ed51e8040c1c05"),
    "user_identity": "用户的唯一标识",
    "room_identity": "房间的唯一标识",
    "created_at": 1,
    "updated_at": 1,
    "room_type": "1"
}
```