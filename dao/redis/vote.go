// Package redis provides ...
package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天

/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录  差值的绝对值：2  +432*2
direction=0时，有两种情况：
	1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  -432
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录  差值的绝对值：1  -432
	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录  差值的绝对值：2  -432*2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600 // 投票过期时间
	scorePerVote     = 432           // 每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

// CreatePost 初始化帖子分数的redis数据
func CreatePost(postID, communityID int64) error {
	// 下面的操作要同时成功，这里要用到事务
	pipeline := rdb.TxPipeline()
	// 帖子发布时间存入redis
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子的默认分数，就是帖子的发布时间
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPrefis + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)

	_, err := pipeline.Exec()
	return err
}

// VoteForPost
func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票的限制
	// 去 redis 取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 下面2里的操作部分和3需要放在同一个事务中，同时成功，同时失败
	// 2. 更新帖子的分数
	// 先查询当前用户给帖子的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostScoreZSet+postID), userID).Val()
	//zap.L().Debug("ov", zap.)
	// 如果这次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	// 如果当前值大于查询的值，则投票是正的
	var op float64
	if value > ov {
		op = 1
	} else { // 如果小于查询的值，说明该操作是要减分数
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值（绝对值）
	pipeline := rdb.TxPipeline()
	/*
		_, err := rdb.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID).Result()
		if err != nil {
			return err
		}
	*/
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3. 记录用户为该帖子投过票的数据
	if value == 0 {
		//_, err = rdb.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Result()
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		/*
			_, err = rdb.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
				Score:  value, // 赞成还是反对票
				Member: userID,
			}).Result()
		*/
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value, // 赞成还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
