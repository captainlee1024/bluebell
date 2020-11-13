// Package mysql provides ...
package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"github.com/captainlee1024/bluebell/models"
)

// 把每一步数据库操作封装成函数
// 待 logic 层根据业务需求调用

const secret = "我是盐！"

// InsertUser 插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	// 在数据库中存储密码一定不能使用明文
	user.Password = encryptPassword(user.Password)
	// 执行 SQL 语句入库
	sqlStr := `insert into user (user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// CheckUserExist 检查制定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	// 执行查询 SQL 语句
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		// 数据库查询出错
		return err
	}
	if count > 0 {
		// 把返回的错误定义成常量，不要在代码中突然出现字符串，要让别人在上面能找到
		//return errors.New("用户已存在")
		return ErrorUserExist
	}
	return
}

// 加盐加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	// 先写入 secret
	h.Write([]byte(secret))
	// 再写入用户的真实密码，然后吧得到的字节转换成16进制的字符串，并返回
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(u *models.User) (err error) {
	oPassword := u.Password // 用户登录的密码
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(u, sqlStr, u.Username)
	if err == sql.ErrNoRows {
		//return errors.New("用户不存在")
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != u.Password {
		//return errors.New("密码错误")
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
