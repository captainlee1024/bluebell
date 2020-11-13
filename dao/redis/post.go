// Package redis provides ...
package redis

import (
	"strconv"
	"time"

	"github.com/captainlee1024/bluebell/models"
	"github.com/go-redis/redis"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 2. 确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// .3 ZREVRANGE 查询
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder 获取指定区间的帖子id，并按分数（根据参数决定　分数/时间）从高到低排序
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1. 根据用户传递过来的order字段从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	/* 抽取出来，封装成getIDsFormKey
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// .3 ZREVRANGE 查询
	return rdb.ZRevRange(key, start, end).Result()
	*/
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
// 注意当使用for循环发送查询的时候注意使用pipeline减少发送请求的RTT
func GetPostVoteData(ids []string) (data []int64, err error) {
	/*
		data = make([]int64, 0, len(ids))
		for _, id := range ids {
			key := getRedisKey(KeyPostVotedZSetPrefix + id)
			// 查找key中分数是1的元素的数量->统计每篇帖子投赞成票的数量
			v := rdb.ZCount(key, "1", "1").Val()
			data = append(data, v)
		}
	*/

	// 使用pipeline一次发送多条命，令减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用zinterstore 吧分区的帖子set与帖子分数的zset生成一个zset
	// 针对新的zset按之前的逻辑取数据
	// 但是这个zinterstore这个函数比较重，所以我们要减少运行次数

	// 社区的key
	ckey := getRedisKey(KeyCommunitySetPrefis + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderKey)
		pipeline.Expire(key, 60*time.Second) // 查一次缓存60S，60秒内在需要数据就从缓存里拿数据，不运行函数
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	// 存在的话直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
