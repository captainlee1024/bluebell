package controller

import (
	"errors"
	"fmt"

	"github.com/captainlee1024/bluebell/dao/mysql"
	"github.com/captainlee1024/bluebell/logic"
	"github.com/captainlee1024/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 注册
// 获取参数并进行校验，进行用户注册并返回响应
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSigUp)
	if err := c.ShouldBindJSON(p); err != nil { // 只能对参数的类型格式进行校验
		// 请求参数有误，直接返回响应
		zap.L().Error("signup with invalid param", zap.Error(err))
		// 使用翻译器
		// 判断 err 是不是 validatoe.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		fmt.Println(p)
		if !ok { // 如果不是该校验器类型的错误，我们就直接返回
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			// 使用封装的的响应方法
			// 返回参数有误
			ResponseError(c, CodeInvalidParam)
			return
		}

		// 是校验器类型的错误就翻译
		//c.JSON(http.StatusOK, gin.H{
		//	//"msg": "请求参数有误",
		//	//"msg": errs.Translate(trans), // 翻译错误
		//	"msg": remvoeTopStruct(errs.Translate(trans)), // 翻译错误信息，并使用 removeTopStruct 函数去除字段名中的结构体名称标识
		//})
		// 使用封装的自定义错误提示信息和状态码
		ResponseErrorWithMsg(c, CodeInvalidParam, remvoeTopStruct(errs.Translate(trans)))
		return
	}
	// 手动对请求参数进行详细的业务规则校验
	// 例如我们这里是注册，业务规则如下
	// 用户名和密码不能为空，密码和确认密码必须相同

	/*
		if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
			zap.L().Error("signup with invalid param")
			c.JSON(http.StatusOK, gin.H{
				"msg": "请求参数有误",
			})
			return
		}
	*/
	//fmt.Println(p)
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		if errors.Is(err, mysql.ErrorUserExist) { // 用户已存在
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录请求的函数
func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断 err 是不是 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam) // 参数错误
			return
		}
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": remvoeTopStruct(errs.Translate(trans)), // 翻译错误
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, remvoeTopStruct(errs.Translate(trans))) // 具体哪些字段错误
		return
	}
	// 2. 业务逻辑处理
	//token, err := logic.Login(p)
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或密码错误",
		//})
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3. 返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登录成功",
	//})
	//ResponseSuccess(c, token)
	//ResponseSuccess(c, gin.H{
	//	"token":    token,
	//	"username": p.Username,
	//})
	ResponseSuccess(c, gin.H{
		// 如果 id 值大于 1<<53-1 的话，传入前端就会失真，因为JS最大就1<<53-1
		// 我们会把userid序列化转换成字符串传到前端，接受的时候再反序列化转化成int64
		// 可以自己实现序列化和反序列化函数，也可以使用tag string
		"userid":   fmt.Sprintf("%d", user.UserID),
		"username": user.Username,
		"token":    user.Token,
	})
}
