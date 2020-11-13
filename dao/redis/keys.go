// Package redis provides ...
package redis

// redis key

// redis key 注意使用注意使用命名空间的方式，方便查询和拆分
// 如果名字比较短，可以把类型加在和面例如：KeyPostTimeZSet
// 如果感觉比较长的话可以不写类型，在后面加上注释
// 例如 KeyPostTime // zset: 贴子及发帖的时间

const (
	Prefix                      = "bluebell:"
	KeyPostTimeZSet             = "post:time"         // zset: 贴子及发帖的时间
	KeyPostScoreZSet            = "post:score"        // zset: 贴子及投票的分数
	KeyPostVotedZSetPrefix      = "post:voted:"       // zset: 记录用户以投票类型，参数是 psot id
	KeyUserAccessTokenSetPrefix = "user:accessToken:" // set 记录用户登录时的accessToken
	KeyCommunitySetPrefis       = "community:"        // set:保存每个分区下帖子的id
)

// 给reids key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
