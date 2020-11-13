// @Title mysql
// @Description 定义在多个模块用到的错误信息
// @Author CaptainLee1024 2020-10-07
// @Update CaptainLee1024 2020-10-07
package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
)
