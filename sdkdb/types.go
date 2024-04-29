package sdkdb

import (
	"gorm.io/gorm"
)

// TxCalls 事务调用链
type TxCalls interface {
	// Add 添加事务调用
	Add(call func(tx *gorm.DB) error)
	// Run 执行调用链
	Run(db *gorm.DB) error
}

type _TxCalls []func(tx *gorm.DB) error

func (calls *_TxCalls) Add(call func(tx *gorm.DB) error) {
	*calls = append(*calls, call)
}

func (calls *_TxCalls) Run(db *gorm.DB) error {
	if len(*calls) == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, call := range *calls {
			if err := call(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

// NewTxCalls 创建数据库调用链
func NewTxCalls() TxCalls {
	return new(_TxCalls)
}
