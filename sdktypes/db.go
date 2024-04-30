package sdktypes

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type BaseModel struct {
	Id         ID        `gorm:"primarykey;type:bigint;not null;comment:ID;" json:"id"`
	CreateTime time.Time `gorm:"->;type:timestamp;default:current_timestamp;not null;comment:创建时间;" json:"createTime"`
	UpdateTime time.Time `gorm:"->;type:timestamp;default:current_timestamp on update current_timestamp;not null;comment:修改时间;" json:"updateTime"`
}

func NewBaseModel(id ID) BaseModel {
	return BaseModel{
		Id:         id,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
}

type PagingSortKey struct {
	Key string `form:"key" json:"key" binding:"required"`
	Asc bool   `form:"asc" json:"asc"`
}

type PagingParam struct {
	Page  int64           `form:"page" json:"page" binding:"gte=0"`
	Rows  int64           `form:"rows" json:"rows" binding:"required,gte=1,lte=1000"`
	Sorts []PagingSortKey `form:"sorts" json:"sorts"`
}

// Offset 偏移
func (p *PagingParam) Offset() int64 {
	return p.Page * p.Rows
}

func (p *PagingParam) Limit() int64 {
	return p.Rows
}

// CheckSortKey 检查排序的key是否合法
func (p *PagingParam) CheckSortKey(allowKeys ...string) error {
	for _, item := range p.Sorts {
		if !slices.Contains(allowKeys, item.Key) {
			return fmt.Errorf("invalid sort key : `%s`", item.Key)
		}
	}
	return nil
}

// SortString 获取排序SQL
func (p *PagingParam) SortSQLString() *string {
	if len(p.Sorts) == 0 {
		return nil
	}
	var sb strings.Builder
	for i := 0; i < len(p.Sorts); i++ {
		item := p.Sorts[i]
		sb.WriteString(item.Key)
		if item.Asc {
			sb.WriteString(" asc")
		} else {
			sb.WriteString(" desc")
		}
		if len(p.Sorts)-1 != i {
			sb.WriteString(",")
		}
	}
	s := sb.String()
	return &s
}

// PagingResult 分页结果
type PagingResult[T any] struct {
	Total int64 `form:"total" json:"total"`
	Rows  []T   `form:"rows" json:"rows"`
}
