package models

// 定义请求的参数结构体

// ParamSigUp 注册请求参数
type ParamSigUp struct {
	Username string `json:"username" binding:"required"`               // 用户名
	Password string `json:"password" binding:"required,checkPassword"` // 用户密码
	//RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 确认密码
	RePassword string `json:"confirm_password" binding:"required,eqfield=Password"` // 确认密码
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 用户密码
}

// ParamVoteData 投票请求数据
type ParamVoteData struct {
	//UserID(可以通过GetCurrentID拿到)
	PostID    string `json:"post_id" binding:"required"`              // 贴子 id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票（1）反对票（-1）取消投票（0）
}

// ParamPostList 获取帖子列表2 query string 参数
// communityID 可以为空
// 如果前端没有传社区id，communityID为空，查询所有帖子列表-->GetPostList2
// 如果前端传了社区id，查询对应社区的帖子列表-->
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 社区ID　可以为空
	Page        int64  `json:"page" form:"page" example:"1"`       // 页码
	Size        int64  `json:"size" form:"size" example:"10"`      // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamCommunityPostList 安社区获取帖子列表query string 参数
type ParamCommunityPostList struct {
	ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
