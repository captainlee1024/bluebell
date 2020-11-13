// @Title controller
// @Description 社区相关
// @Author CaptainLee1024 2020-10-06
// @Update CaptainLee1024 2020-10-06
package controller

import (
	"strconv"

	"github.com/captainlee1024/bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 社区分类

func CommunityHandler(c *gin.Context) { // 参数的校验，适当的做一些路由的跳转，具体业务逻辑在 logic 中实现
	// 查询到所有的社区（community_id, communiry_name）以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		// 后端出错一般是不对外暴露详细信息，记录在日志文件中
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给前端
		return
	}

	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区ID
	// 获取参数
	idStr := c.Param("id")
	// 获取到的是字符串，要转换成 int64 类型的
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil { // 如果出错，是参数有误
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 查询社区的详细信息
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, data)

}
