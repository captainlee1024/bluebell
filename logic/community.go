package logic

import (
	"github.com/captainlee1024/bluebell/dao/mysql"
	"github.com/captainlee1024/bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库，找到所有的 community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
