// Package logic provides ...
package logic

import (
	"strconv"

	"github.com/captainlee1024/bluebell/dao/redis"
	"github.com/captainlee1024/bluebell/models"
	"go.uber.org/zap"
)

// 推荐阅读
// 基于用户投票相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 这里使用简化的投票分数
// 一个赞成票加432分
// 根据分数排序，新的一天的帖子时间戳作为分数肯定排在昨天的帖子前面
// 要想昨天的帖子排在今天的帖子前面就要分数比今天的高，该帖子就需要获得一天时间对应的票数
// 一天86400秒，就需要加86400分，如果我们设定100票就让昨天的帖子出现在前面，那一票就是864
// 我们这里设值200篇会让昨天的帖子在今天的没有获得赞成票的帖子前面

/* 投票的几种情况
direction=1时，有两种情况
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在投赞成票
direction=0时，有两种情况
	1. 之前投赞成票，现在要取消投票
	2. 之前投反对票，现在要取消投票
direction=-1时，有两种情况
	1. 之前没有投票，现在投反对票
	2. 之前投赞成票，现在投反对票

投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票
因为几年前的帖子几乎没人再看再投票了，这时候突然来一个请求个
一个星期之后就能够把数据持久化到mysql中了
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除对应的 KeyPostVotedZSetPrefix
*/

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
