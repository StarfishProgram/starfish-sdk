package sdktypes

import (
	"strings"
	"time"
)

type BaseModel struct {
	Id         ID        `gorm:"primarykey;type:bigint;not null;comment:ID;" json:"id" form:"id"`
	CreateTime time.Time `gorm:"->;type:timestamp;default:current_timestamp;not null;comment:创建时间;" json:"createTime" form:"createTime"`
	UpdateTime time.Time `gorm:"->;type:timestamp;default:current_timestamp on update current_timestamp;not null;comment:修改时间;" json:"updateTime" form:"updateTime"`
}

func NewBaseModel(id ID) BaseModel {
	return BaseModel{
		Id:         id,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
}

type PagingSortKey struct {
	Key string `form:"key" json:"key"`
	Asc bool   `form:"asc" json:"asc"`
}

type PagingParam struct {
	Page  int64           `form:"page" json:"page" binding:"gte=0,lte=1000"`
	Rows  int64           `form:"rows" json:"rows" binding:"required,gte=1,lte=10000"`
	Sorts []PagingSortKey `form:"sorts" json:"sorts"`
}

// Offset 偏移
func (p *PagingParam) Offset() int64 {
	return p.Page * p.Rows
}

func (p *PagingParam) Limit() int64 {
	return p.Rows
}

// SortSQLString 获取排序SQL
func (p *PagingParam) SortSQLString(rules map[string]string) *string {
	if len(p.Sorts) == 0 {
		return nil
	}
	var sb strings.Builder

	for i := 0; i < len(p.Sorts); i++ {
		item := p.Sorts[i]
		field, exists := rules[item.Key]
		if !exists {
			continue
		}
		sb.WriteString(field)
		if item.Asc {
			sb.WriteString(" asc")
		} else {
			sb.WriteString(" desc")
		}
		if len(p.Sorts)-1 != i {
			sb.WriteString(",")
		}
	}
	if sb.Len() == 0 {
		return nil
	}
	s := sb.String()
	return &s
}

// PagingResult 分页结果
type PagingResult[T any] struct {
	Total int64 `form:"total" json:"total"`
	Rows  []T   `form:"rows" json:"rows"`
}
