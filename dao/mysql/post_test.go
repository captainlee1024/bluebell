package mysql

import (
	"testing"

	"github.com/captainlee1024/bluebell/models"
	"github.com/captainlee1024/bluebell/settings"
)

func init() {
	dbcfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "644315",
		DbName:       "bluebell",
		Port:         3306,
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}
	err := Init(&dbcfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	// 注意！运行的时候，只会运行单元测试里的代码，没有对db进行初始化，会报空指针错误
	// 我们要写个初始化函数进行初始化
	post := models.Post{
		ID:          12,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Conent:      "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert into mysql success")
}
