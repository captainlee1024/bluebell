package controller

import (
	"github.com/captainlee1024/bluebell/logic"
	"github.com/captainlee1024/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

// PostVoteHandler 投票处理函数
func PostVoteHandler(c *gin.Context) {
	// 获取参数并进行参数校验
	p := new(models.ParamVoteData) // 使用new()创建结构体指针
	if err := c.ShouldBindJSON(p); err != nil {
		// 首先进行类断言，判断是不是validator的错误类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok { // 说明不是该错误类型，返回参数有误
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 如果是里面的错误类型，则返回具体的错误
		errData := remvoeTopStruct(errs.Translate(trans)) // 并去除错误提示中的结构体标示
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 获取当前请求用户的 ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost(userID, p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
