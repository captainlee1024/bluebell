package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var (
	ErrorKeyNotExist = errors.New("token不存在，请登录")
)

func SetAToken(userID string, aToken string) error {
	//err := rdb.Set(key, aToken, time.Hour*24*7).Err()
	err := rdb.Set(getRedisKey(KeyUserAccessTokenSetPrefix+userID), aToken, time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetAToken(userID string) (redisToken string, err error) {
	var val string
	//val, err = rdb.Get(key).Result()
	val, err = rdb.Get(getRedisKey(KeyUserAccessTokenSetPrefix + userID)).Result()
	if err == redis.Nil {
		return "", ErrorKeyNotExist
	} else if err != nil {
		return
	} else {
		redisToken = val
		return
	}
}
