// Package controller provides ...
package controller

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/captainlee1024/bluebell/models"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器T
var trans ut.Translator

// InitTrans 初始化翻译器
// param local 指定翻译成哪种语言
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	// 拿到 gin 框架中的校验器引擎，用来注册我们指定的翻译器
	// binding.Validator.Engine()拿到引擎之后转换成需要的 Validator 类型
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取 json tag 的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "_" {
				return ""
			}
			return name
		})

		// 为 ParamSignUp 注册校验方法
		v.RegisterStructValidation(ParamSignUpStructLevelValidation, models.ParamSigUp{})

		// 为校验器注册自定义的字段级别的校验方法
		if err := v.RegisterValidation("checkPassword", customFunc); err != nil {
			return err
		}

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale) // 使用指定的或者从 'Accept-Language' 获得的信息去初始化一个翻译器
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		// 校验器是真正做校验的，这里注册到在 gin 中拿到的校验器中
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		if err != nil {
			return err
		}
		// 在 InitTrans 中通过调用 RegisterTranslation() 方法来注册我们为自定义字段级别校验方法自定义的翻译方法
		// 注意！因为这里会使用到 trans 实例
		// 所以这一步注册要方放到 trans 初始化的后面
		if err := v.RegisterTranslation(
			"checkPassword",
			trans,
			registerTranslator("checkPassword", "密码必须不小于6位"),
			translate,
		); err != nil {
			return err
		}
		return
	}
	return
}

// 去除错误信息里的结构体前缀
func remvoeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// 为 ParamSignUp 自定义一个校验方法，用于比较密码和确认密码，并返回错误信息
func ParamSignUpStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.ParamSigUp)

	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的 param
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
	}
}

// customFunc 自定义字段（Password）级别校验方法
// 检验密码长度是否小于6
func customFunc(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 {
		return false
	}
	return true
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		if err := ut.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(ut ut.Translator, fe validator.FieldError) string {
	msg, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

/*
type SignUpParam struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func main() {
	if err := InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		var u SignUpParam
		if err := c.ShouldBind(&u); err != nil {
			// 获取validator.ValidationErrors类型的errors
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				// 非validator.ValidationErrors类型错误直接返回
				c.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
				return
			}
			// validator.ValidationErrors类型错误则进行翻译
			c.JSON(http.StatusOK, gin.H{
				"msg":errs.Translate(trans),
			})
			return
		}
		// 保存入库等具体业务逻辑代码...

		c.JSON(http.StatusOK, "success")
	})

	_ = r.Run(":8999")
}
*/
