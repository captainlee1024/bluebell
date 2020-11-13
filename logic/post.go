package logic

import (
	"github.com/captainlee1024/bluebell/dao/mysql"
	"github.com/captainlee1024/bluebell/dao/redis"
	"github.com/captainlee1024/bluebell/models"
	"github.com/captainlee1024/bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成 postID
	p.ID = snowflake.GenID()

	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	// 3. 初始化投票数据
	err = redis.CreatePost(p.ID, p.CommunityID)
	// 4. 返回相应
	//return mysql.CreatePost(p)
	return
}

/*
func GetPostById(pid int64) (date *models.Post, err error) {
	return mysql.GetPostById(pid)
}
*/
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	// 查询帖子信息
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	// 根据帖子信息里的作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("post.AuthorID", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("post.CommunityID", post.CommunityID),
			zap.Error(err))
		return
	}
	// 组合成接口需要的数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 获取所有帖子
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	// 首先初始化返回值定义的变量，那里只是声明，并没有申请内存
	data = make([]*models.ApiPostDetail, 0, len(posts))

	// 依次获取所有帖子的作者和社区信息
	for _, post := range posts {
		// 根据帖子信息里的作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("post.AuthorID", post.AuthorID),
				zap.Error(err))
			//return
			continue
		}
		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("post.CommunityID", post.CommunityID),
				zap.Error(err))
			//return
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去 redis 查询 id 列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListByIDs(p) return 0 data")
		return
	}
	// 根据 id 去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询号每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	// 依次获取所有帖子的作者和社区信息
	for idx, post := range posts {
		// 根据帖子信息里的作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("post.AuthorID", post.AuthorID),
				zap.Error(err))
			//return
			continue
		}
		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("post.CommunityID", post.CommunityID),
				zap.Error(err))
			//return
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityPostList
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去 redis 查询 id 列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListByIDs(p) return 0 data")
		return
	}
	// 根据 id 去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询号每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	// 依次获取所有帖子的作者和社区信息
	for idx, post := range posts {
		// 根据帖子信息里的作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("post.AuthorID", post.AuthorID),
				zap.Error(err))
			//return
			continue
		}
		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("post.CommunityID", post.CommunityID),
				zap.Error(err))
			//return
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将两个帖子列表查询逻辑合二为一的函数
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		// 查询所有帖子列表
		data, err = GetPostList2(p)
	} else {
		// 根据社区id查询帖子列表
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
