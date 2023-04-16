package server_core

import (
	"GoIm/define"
	"GoIm/models"
	"GoIm/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// UserQueryResult 返回结构体
type UserQueryResult struct {
	Nickname string `bson:"nickname"`
	Sex      int    `bson:"sex"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar"`
	IsFriend bool   `bson:"is_Friend"` //
}

// Login 登录
func Login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	fmt.Println(account, password)
	fmt.Println(utils.GetMd5(password))
	ub, err := models.GetUserBasicByAccountPassword(account, utils.GetMd5(password))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg":  "用户名或密码错误",
		})
		return
	}
	token, err := utils.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -3,
			"msg":  "系统内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token0": token,
		},
	})
}

// UserDetailIdentidy 用户详情
func UserDetailIdentidy(c *gin.Context) {
	user, _ := c.Get("user_Claims")
	uc := user.(*utils.UserClaims)
	userBasic, err := models.GetUserBasicByIdentity(uc.Identity)
	if err != nil {
		log.Println("DB error", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据查询成功",
		"data": userBasic,
	})
}

// SendCode 发送验证码
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
	}
	count, err := models.GetUserBasicCountByEmail(email)
	if err != nil {
		log.Printf("[DB ERROR]===> %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg":  "服务器正忙,请稍后重试",
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -3,
			"msg":  "当前邮箱已经注册过,请找回密码或者更换邮箱注册",
		})
		return
	}
	code := utils.GetCode()
	err = utils.SendCode(email, code)
	if err != nil {
		log.Printf("[SendCode ERROR]===> %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "服务器正忙,请稍后重试",
		})
		return
	}
	err = models.RDB.Set(context.Background(), define.RegisterPrefix+email, code, define.ExpireTime).Err()
	if err != nil {
		log.Printf("[Redis ERROR]===> %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "服务器正忙,请稍后重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	code := c.PostForm("code")
	email := c.PostForm("email")
	account := c.PostForm("account")
	password := c.PostForm("password")
	if code == "" || email == "" || account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册参数错误",
		})
		return
	}
	//判断账号是否唯一
	cnt, err := models.GetUserBasicCountByAccount(account)
	if err != nil {
		log.Fatal("[DB Error] 用户account查询失败" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败,请重试",
		})
		return
	}
	if cnt > 0 {
		log.Fatal("[Register Error] 用户account已经被注册")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败,请重试",
		})
		return
	}
	//验证码
	recode, err := models.RDB.Get(context.Background(), define.RegisterPrefix+email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败,请重试",
		})
		return
	}
	if recode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "验证码输入错误,请重新获取验证码",
		})
	}

	ub := new(models.UserBasic)
	ub.Identity = utils.GetUUID()
	ub.Account = account
	ub.Email = email
	ub.Password = utils.GetMd5(password)
	ub.CreatedAt = time.Now().Unix()
	ub.UpdatedAt = time.Now().Unix()
	err = models.InsertOneUserBasic(ub)
	if err != nil {
		log.Fatal("[DB Error] 用户account查询失败" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败,请重试",
		})
		return
	}
	token, err := utils.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -3,
			"msg":  "系统内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token0": token,
		},
	})

}

// Query 查询指定用户的个人信息
func Query(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "用户参数错误",
		})
		return
	}
	ub, err := models.QueryUserByAccount(account)
	if err != nil {
		log.Fatal("[DB ERROR]" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "请稍后重试",
		})
		return
	}
	uc := c.MustGet("user_claims").(utils.UserClaims)
	if b, err := models.GetJudgeUserIsFriend(ub.Identity, uc.Identity); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "请稍后重试",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "查询成功",
			"data": UserQueryResult{
				Nickname: ub.Nickname,
				Sex:      ub.Sex,
				Email:    ub.Email,
				Avatar:   ub.Avatar,
				IsFriend: b,
			},
		})
	}

}

// AddFriend 添加好友
func AddFriend(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "用户参数错误",
		})
		return
	}
	ub, err := models.QueryUserByAccount(account)
	if err != nil {
		log.Fatal("[DB Error] 用户account查询失败" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "添加失败,请重试",
		})
		return
	}
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	if bo, err := models.GetJudgeUserIsFriend(uc.Identity, ub.Identity); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "请稍后重试",
		})
		return
	} else if bo == true {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "不可重复添加",
		})
		return
	} else {
		//保存房间记录
		rb := &models.RoomBasic{
			Identity:     utils.GetUUID(),
			Number:       "",
			Name:         "",
			Info:         "",
			UserIdentity: uc.Identity,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		if err = models.InsertOneRoomBasic(rb); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -4,
				"msg":  "房间创建失败",
			})
			return
		}
		//保存用户好友信息
		ur := models.UserRoom{
			UserIdentity: uc.Identity,
			RoomIdentity: rb.Identity,
			RoomType:     1,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		if err = models.InsertOneUserRoom(&ur); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -4,
				"msg":  "好友添加失败，稍后重试",
			})
			return
		}
		ur1 := models.UserRoom{
			UserIdentity: ub.Identity,
			RoomIdentity: rb.Identity,
			RoomType:     1,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		if err = models.InsertOneUserRoom(&ur1); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -4,
				"msg":  "好友添加失败，稍后重试",
			})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "好友添加成功",
	})

}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "用户参数错误",
		})
		return
	}
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	//获取房间identity
	roomIdentity := models.GetUserRoomIdentity(identity, uc.Identity)
	if roomIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "不是好友，不用删除",
		})
	}
	//删除userroom
	if err := models.DeleteUserRoom(roomIdentity); err != nil {
		log.Fatal("[DB ERROR]" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "删除失败",
		})
	}

	//删除roombasic
	if err := models.DeleteRoomBasic(roomIdentity); err != nil {
		log.Fatal("[DB ERROR]" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "删除失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "好友删除成功",
	})

}
