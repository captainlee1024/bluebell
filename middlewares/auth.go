package middlewares

import (
	"fmt"
	"strings"

	"github.com/captainlee1024/bluebell/controller"
	"github.com/captainlee1024/bluebell/dao/redis"
	"github.com/captainlee1024/bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带 Token 有三种方式 1. 放在请求头中 2. 放在请求体中 3. 放在 URI
		// 这里假设 Token 放在 Header 的 Autherization 中，并使用 Bearer 开头
		// Authorization: Bearer xxx.xxx.xxx 或者是前端团队自己定义的其他格式，例如：X-TOKEN xx.xx.xx
		// 具体根据自己的团队和业务逻辑决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割，获取 JWT
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1] 是获取到的 tokenString, 我们使用之前定义好的解析 JWT 的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		// 解析有效之后在判断是否与 Redis 中存的 token 相等
		redisToken, err := redis.GetAToken(fmt.Sprint(mc.UserID))
		if err != nil {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		if parts[1] != redisToken {
			controller.ResponseError(c, controller.CodeLoginElsewhere)
			c.Abort()
			return
		}
		//fmt.Println(authHeader, parts)
		// 将当前请求的 userID 信息保存到请求的上下文 c 中
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以通过 c.Get(controller.CtxUserIDKey) 来获取当前请求用户的 userID
	}
}
