// Package logic provides ...
package logic

import (
	"fmt"

	"github.com/captainlee1024/bluebell/dao/mysql"
	"github.com/captainlee1024/bluebell/dao/redis"
	"github.com/captainlee1024/bluebell/models"
	"github.com/captainlee1024/bluebell/pkg/jwt"
	"github.com/captainlee1024/bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

// SignUp 注册
func SignUp(p *models.ParamSigUp) (err error) {
	// 1.判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成UID
	userID := snowflake.GenID()
	// 构造一个 User 实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.密码加密
	// 4.保存进数据库
	if err = mysql.InsertUser(user); err != nil {
		// 记录日志
		return
	}
	// 5.redis.xxx ...
	// ...
	return
}

// Login 登录
//func Login(p *models.ParamLogin) (token string, err error) {
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	//return mysql.Login(user)
	// 注意这里传递的是指针，这样改变才会影响到外面的值，才能拿到 user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成 JWT
	//return jwt.GenToken(user.UserID, user.Username)
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	// token 添加到 redis 中
	redis.SetAToken(fmt.Sprint(user.UserID), token)
	user.Token = token
	return
}
