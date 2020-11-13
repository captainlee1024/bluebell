package models

import "time"

// Post 即需要与用户传递的参数对应，也需要与数据库表对应
// 注意内存对齐
type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`                            // 帖子ID
	AuthorID    int64     `json:"author_id" db:"author_id"`                          // 用户ID
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"` // 社区ID　可以为空
	Status      int32     `json:"status" db:"status"`                                // 投票状态
	Title       string    `json:"title" db:"title" binding:"required"`               // 标题
	Conent      string    `json:"content" db:"content" binding:"required"`           // 内容
	CreateTime  time.Time `json:"create_time" db:"create_time"`                      // 创建时间
}

// ApiPostDetail 帖子详情接口结构体
// 因为在帖子详情中我们返回给前端的信息中是author_id，前段希望拿到的是作者的名字
// 这样才可以进行渲染
type ApiPostDetail struct {
	AuthorName       string             `json:"author_id"` // 用户名
	VoteNum          int64              `json:"vote_num"`  // 票数
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息结构体
}
