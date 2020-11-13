// @Title controller
// @Description 帖子相关
// @Author CaptainLee1024 2020-10-07
// @Update CaptainLee1024 2020-10-07
package controller

import (
	"strconv"

	"github.com/captainlee1024/bluebell/logic"
	"github.com/captainlee1024/bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子接口
// @Description 根据用户输入的数据创建一个帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object body models.Post false "帖子参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验
	// 定义一个模型
	// ==> shouldbindjson 需要绑定模型，这里需要定义模型，放在 models 里
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil { // validator --> binding tag
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create pot with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c 中取到当前发请求的用户的 ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回相应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的接口
// @Summary 获取帖子详情的接口
// @Description 根据传入的postid查询帖子的详细信息
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path string true "帖子ID"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子的id）并进行校验
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		return
	}

	// 2. 根据 id 取出帖子的数据（数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的接口
// @Summary 获取帖子列表的接口
// @Description 获取所有帖子列表，根据传递的参数进行分页，按照发部顺序进行排序
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	// 抽取出来封装成一个函数
	page, size := getPageInfo(c)
	/*
		pageStr := c.Query("page")
		sizeStr := c.Query("size")

		var (
			size int64
			page int64
			err  error
		)

		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			page = 1
		}
		size, err = strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			size = 10
		}
	*/

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
	// 返回响应
}

// GetPostListHandler2 升级版帖子列表接口 根据前端传来的参数动态的获取帖子列表 按创建时间或者按分数进行排序
// @Summary 升级版帖子列表接口
// @Description 可按社区时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET 请求参数(query string)：/api/v1/posts2?page=1&size=10&order=score
	// 获取参数
	// 这里可以利用初始化结构体时指定默认值
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, // 不要出现 magic string
	}
	// 使用shouldBindQuery 需要通过tag`form`指定获取的名字
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// c.ShouldBind() 根据请求的数据选择相应的方法获取数据
	// c.ShouldBinJSON() 如果请求中携带的是json格式的数据，才能用这个方法取到数据，使用`json`tag

	/*
		// 修改，两个接口合二为一：两个都是查询帖子列表的接口，参数的不同查询逻辑不同所以可以使用一个controller
		// 传入参数，在logic层写个函数进行判断调用哪个逻辑处理的代码
		// 但是这些代码也属于业务逻辑相关代码，最好放在logic层，我们在logic层封装成一个函数
			if p.CommunityID == 0 {
				data, err := logic.GetPostList2(p)
			} else {
				data, err := logic.GetCommunityPostList(p)
			}
	*/

	// 获取数据
	//data, err := logic.GetPostList2(p)
	data, err := logic.GetPostListNew(p) // 更新：两个帖子列表查询接口合二为一
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
	// 返回响应
}

/*
	两个查询帖子列表的逻辑层合并之前的社区查询的Handler
// GetCommunityPostListHandler 根据社区去查询帖子列表
func GetCommunityPostListHandler(c *gin.Context) {
	// GET 请求参数(query string)：/api/v1/posts2?page=1&size=10&order=score
	// 获取参数
	// 这里可以利用初始化结构体时指定默认值
	p := &models.ParamCommunityPostList{
		ParamPostList: models.ParamPostList{
			Page:  1,
			Size:  10,
			Order: models.OrderTime,
		},
	}
	// 使用shouldBindQuery 需要通过tag`form`指定获取的名字
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// c.ShouldBind() 根据请求的数据选择相应的方法获取数据
	// c.ShouldBinJSON() 如果请求中携带的是json格式的数据，才能用这个方法取到数据，使用`json`tag

	// 获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
	// 返回响应
}
*/
