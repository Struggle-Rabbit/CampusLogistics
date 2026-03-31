package utils

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator 获取全局校验器
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		// 注册自定义校验器
		_ = validate.RegisterValidation("mobile", ValidateMobile)
	})
	return validate
}

// BindAndValidate 绑定并校验参数
func ShouldBind(c *gin.Context, obj any) bool {
	// 绑定 JSON/Form
	if err := c.ShouldBind(obj); err != nil {
		msg := Translate(err)
		Fail(c, msg)
		return false
	}

	// 手动校验
	v := GetValidator()
	if err := v.Struct(obj); err != nil {
		msg := Translate(err)
		Fail(c, msg)
		return false
	}
	return true
}

// Translate 翻译校验错误
func Translate(err error) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return constant.MsgParamFormatError
	}

	for _, fe := range validationErrs {
		field := fe.Field()
		tag := fe.Tag()
		param := fe.Param()

		switch tag {
		case "required":
			return fmt.Sprintf("%s 不能为空", field)
		case "email":
			return fmt.Sprintf("%s 格式不正确", field)
		case "len":
			return fmt.Sprintf("%s 长度必须为 %s 位", field, param)
		case "min":
			return fmt.Sprintf("%s 最小长度为 %s", field, param)
		case "max":
			return fmt.Sprintf("%s 最大长度为 %s", field, param)
		case "gt", "gte", "lt", "lte":
			return fmt.Sprintf("%s 数值范围不合法", field)
		case "mobile":
			return fmt.Sprintf("%s 不是合法手机号", field)
		default:
			return fmt.Sprintf("%s 格式错误", field)
		}
	}
	return constant.MsgParamValidationFailed
}

// ValidateMobile 手机号校验
var mobileReg = regexp.MustCompile(`^1[3-9]\d{9}$`)

func ValidateMobile(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return mobileReg.MatchString(val)
}
