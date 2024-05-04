package sdktypes

import (
	"fmt"
	"strings"
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdk"
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
	Key  string `form:"key" json:"key"`
	Desc bool   `form:"desc" json:"desc"`
}

type PagingParam struct {
	Page  int64           `form:"current" json:"current"`
	Rows  int64           `form:"pageSize" json:"pageSize"`
	Sorts []PagingSortKey `form:"sorts" json:"sorts"`
}

// Offset 偏移
func (p *PagingParam) Offset() int {
	if p.Page < 1 || p.Page > 1000 {
		return 0
	}
	return int((p.Page - 1)) * p.Limit()
}

func (p *PagingParam) Limit() int {
	if p.Rows <= 0 || p.Rows > 10000 {
		return 30
	}
	return int(p.Rows)
}

// SortSQLString 获取排序SQL
func (p *PagingParam) SortSQLString(rules map[string]string) *string {
	if len(p.Sorts) == 0 {
		return nil
	}
	fieldSorts := make([]string, 0, len(p.Sorts))
	for i := 0; i < len(p.Sorts); i++ {
		item := p.Sorts[i]
		field, exists := rules[item.Key]
		if !exists {
			continue
		}
		fieldSort := fmt.Sprintf("%s %s", field, sdk.If(item.Desc, "desc", "asc"))
		fieldSorts = append(fieldSorts, fieldSort)
	}
	if len(fieldSorts) == 0 {
		return nil
	}
	sql := strings.Join(fieldSorts, ", ")
	return &sql
}

// PagingResult 分页结果
type PagingResult[T any] struct {
	Total int64 `form:"total" json:"total"`
	Rows  []T   `form:"rows" json:"rows"`
}
