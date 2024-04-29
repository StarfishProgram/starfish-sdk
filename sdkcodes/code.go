package sdkcodes

import (
	"fmt"
)

type _Code struct {
	code int64
	msg  string
	i18n string
}

func (c *_Code) Code() int64 {
	return c.code
}
func (c *_Code) Msg() string {
	return c.msg
}
func (c *_Code) I18n() string {
	return c.i18n
}
func (c *_Code) WithMsg(format string, args ...any) Code {
	return New(c.code, fmt.Sprintf(format, args...), c.i18n)
}
func (c *_Code) Error() string {
	return fmt.Sprintf("状态码 = %d, 消息 = %s", c.code, c.msg)
}

// New 创建状态码
func New(code int64, msg string, i18n string) Code {
	return &_Code{code, msg, i18n}
}

var (
	OK              = New(0, "OK", "")
	Internal        = New(1, "服务异常", "")
	AccessLimited   = New(2, "访问受限", "AccessLimited")
	ParamInvalid    = New(3, "参数错误", "")
	TooManyRequests = New(4, "请求过于频繁", "TooManyRequests")
)
