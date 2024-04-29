package sdkcodes

// Code 状态码
type Code interface {
	// Code 状态码
	Code() int64
	// Msg 消息
	Msg() string
	// I18n 国际化key
	I18n() string
	/// WithMsg 替换消息
	WithMsg(format string, args ...any) Code

	error
}
