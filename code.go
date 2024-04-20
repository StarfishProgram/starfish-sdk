package starfish_sdk

import "fmt"

type ICode interface {
	// Code 状态码
	Code() int
	// Msg 消息
	Msg() string
	// I18n 国际化
	I18n() string
	// WithMsg 替换消息
	WithMsg(msg string) ICode
	// WithMsgf 替换消息
	WithMsgf(format string, args ...interface{}) ICode
	// WithError 替换消息
	WithError(err error) ICode
	// Error 错误消息
	Error() string
}

type _code struct {
	code int
	msg  string
	i18n string
}

func (c *_code) Code() int {
	return c.code
}

func (c *_code) Msg() string {
	return c.msg
}

func (c *_code) I18n() string {
	return c.i18n
}

func (c *_code) WithMsg(msg string) ICode {
	return &_code{
		code: c.code,
		msg:  msg,
		i18n: c.i18n,
	}
}

func (c *_code) WithMsgf(format string, args ...interface{}) ICode {
	return c.WithMsg(fmt.Sprintf(format, args...))
}

func (c *_code) WithError(err error) ICode {
	if err != nil {
		return c.WithMsg(err.Error())
	}
	return c
}

func (c *_code) Error() string {
	return fmt.Sprintf("状态码 = %v, 消息 = %s", c.code, c.msg)
}

// 创建Code
func NewCode(code int, msg string, i18n string) ICode {
	return &_code{code, msg, i18n}
}

var (
	CodeOK           = NewCode(0, "OK", "")
	CodeServerError  = NewCode(1, "服务异常", "")
	CodeAccessInvlid = NewCode(2, "访问受限", "access-limit")
	CodeParamInvalid = NewCode(3, "参数错误", "")
)
