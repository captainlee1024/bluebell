// Package router provides ...
package router

import (
	"net/http"

	"github.com/captainlee1024/bluebell/controller"
	_ "github.com/captainlee1024/bluebell/docs"
	"github.com/captainlee1024/bluebell/logger"
	"github.com/captainlee1024/bluebell/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

func Setup(mod string) *gin.Engine {
	if mod == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 当系统配置为 release 的时候配置为发布模式，其他都设置成 debug 模式
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Error("init validatoe trans failed ", zap.Error(err))
		panic(err)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 全局注册限流中间件 并设置令牌桶总容量为1，两秒钟放置一个令牌
	//r.Use(middlewares.RateLimitMiddleware(2*time.Second, 1))

	// 加载静态资源
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	// 注册根路由，访问index页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 注册 swagger api 相关路由
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 注册业务路由
	//r.POST("/signup", controller.SignUpHandler)
	v1.POST("/signup", controller.SignUpHandler)
	// 注册登录路由
	//r.POST("/login", controller.LoginHandler)
	v1.POST("/login", controller.LoginHandler)

	// 这些模块应该是不需要登录的，所以我们把它拿出来
	// 社区分类列表
	v1.GET("/community", controller.CommunityHandler)
	// 社区类别详情
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	// 查看帖子详情
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	// 查看帖子列表
	v1.GET("/posts", controller.GetPostListHandler)

	// 根据参数决定按分数还是时间排序
	v1.GET("/posts2", controller.GetPostListHandler2)

	// 使用 JWT 认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		// 社区分类列表
		//v1.GET("/community", controller.CommunityHandler)
		// 社区类别详情
		//v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 创建帖子
		v1.POST("/post", controller.CreatePostHandler)
		// 查看帖子详情
		//v1.GET("/post/:id", controller.GetPostDetailHandler)
		// 查看帖子列表
		//v1.GET("/posts", controller.GetPostListHandler)

		// 贴子投票
		v1.POST("/vote", controller.PostVoteHandler)

		// 根据时间或者分数获取帖子列表
		//v1.GET("/posts2", controller.GetPostListHandler2)
		/*
			v1.GET("/posts2", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"msg": "ok",
				})
			})
		*/
	}

	/*
		r.GET("/", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
			//time.Sleep(5 * time.Second)
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		})
	*/

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	// 注册 pprof 相关路由
	pprof.Register(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
