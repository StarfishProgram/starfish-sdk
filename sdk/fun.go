package sdk

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	decimal.DivisionPrecision = 16
}

// Waiting 信号阻塞等待
func Waiting() os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	sign := <-ch
	return sign
}

// Context 创建一个上下文
func Context() context.Context {
	return context.Background()
}

// Point 返回数据指针
func Point[P any](p P) *P {
	return &p
}

// If 如果条件成立，返回r1，不成立返回r2
func If[R any](logic bool, r1, r2 R) R {
	if logic {
		return r1
	} else {
		return r2
	}
}

// If 如果条件成立，返回r1，不成立返回r2
func IfCall[R any](logic bool, r1, r2 func() R) R {
	if logic {
		return r1()
	} else {
		return r2()
	}
}

func IfNil[R any](r *R, d R) R {
	if r != nil {
		return *r
	}
	return d
}

// NilDefault 如果data为nil, 则返回default
func NilDefault[P any](data *P, defaultVal P) P {
	if data == nil {
		return defaultVal
	} else {
		return *data
	}
}

// Filter 数据过滤
func Filter[S ~[]E, E any](datas S, f func(item E) bool) S {
	result := make(S, 0, len(datas))
	for i := range datas {
		data := datas[i]
		r := f(data)
		if r {
			result = append(result, data)
		}
	}
	return result
}

// Match 条件匹配
func Match[T comparable](v T, matchers ...T) bool {
	for i := range matchers {
		if v == matchers[i] {
			return true
		}
	}
	return false
}

// Map 数据转换
func Map[P any, R any](datas []P, f func(item P) R) []R {
	result := make([]R, 0, len(datas))
	for i := range datas {
		data := datas[i]
		r := f(data)
		result = append(result, r)
	}
	return result
}

// 安全调用Goroutine
func Go(call func()) {
	go func() {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if code, ok := err.(sdkcodes.Code); ok {
				sdklog.AddCallerSkip(3).Warn(code)
				return
			}
			sdklog.AddCallerSkip(2).Error(err)
		}()
		call()
	}()
}

// Assert
func Assert(expr bool, code ...sdkcodes.Code) {
	if !expr {
		if code != nil {
			panic(code[0])
		} else {
			panic(sdkcodes.Internal.WithMsg("expr is false"))
		}
	}
}

// AssertNil
func AssertNil[V any](v *V, code ...sdkcodes.Code) {
	if v == nil {
		if code != nil {
			panic(code[0])
		} else {
			panic(sdkcodes.Internal.WithMsg("value is nil"))
		}
	}
}

// AssertError
func AssertError(err error, code ...sdkcodes.Code) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if err != nil {
		if code != nil {
			panic(code[0])
		} else {
			panic(sdkcodes.Internal.WithMsg("%s", err.Error()))
		}
	}
}

// Reverse 翻转数组
func Reverse[S ~[]E, E any](s S) {
	slices.Reverse(s)
}
