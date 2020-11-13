package models

// User 用户对象，定义了用户的基础信息
type User struct {
	UserID   int64  `db:"user_id"`  // 用户 ID
	Username string `db:"username"` // 用户名
	Password string `db:"password"` // 邮箱
	Token    string
}
