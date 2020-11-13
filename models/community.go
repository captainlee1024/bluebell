package models

import "time"

// Community 社区类别
type Community struct {
	ID   int64  `json:"id" db:"community_id"`     // 社区ID
	Name string `json:"name" db:"community_name"` // 社区名
}

// CommunityDetail 社区详情
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"` // 这里的 omitempty 是当该字段为空的时候，就不展示该项内容
	CreateTime   time.Time `json:"create_time" db:"create_time"`             // 这里使用 time.Time 类型，在连接数据库的时候就要加上 parseTime=true&loc=Local
}
