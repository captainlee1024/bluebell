package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHander(t *testing.T) {
	// 注意！在编写单元测试的时候，我们项目划分的细，容易导致循环引用的问题
	// 这里的router.Setup()就会导致循环引用问题
	// 所以这里我们自己起一个router
	//r := router.Setup(settings.Conf.Name)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
"community_id": 1,
"title": "test",
"content": "just a test"
	}`

	// 把body转换成io.Reader类型，接收的参数是字节类型的，先转换成byte的一个数组
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// 我们没设置登录的用户ID，所以在解析的时候一定会报错，提示需要登录
	// 判断响应的内容是不是按预期返回了需要登录的信息
	// 两种方法

	// 1. 直接使用Contanins判断响应内容中包含不包含指定的字符串
	//assert.Contains(t, w.Body.String(), "请登录")

	// 2. 将响应内容反序列化到ResponseData，然后判断字段与语气是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)

	// 我们还可以进行其他测试，例如我们可以把body里的字段删除一个，那么程序就走不到NeedLogin那个错误
	// 会提前返回参数有误的错误，可以进行测试
	// ...
}
