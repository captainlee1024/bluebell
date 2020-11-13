package mysql

import (
	"database/sql"

	"github.com/captainlee1024/bluebell/models"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			// 注意，这里与登录不同，查询不到不能返回错误，应该返回空的信息
			// 记录一条 warnlevel 的日志信息
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select
		community_id, community_name, introduction
		from community
		where community_id = ?`
	if err = db.Get(communityDetail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDetail, err
}
