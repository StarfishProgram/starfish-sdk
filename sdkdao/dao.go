package sdkdao

import (
	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdkdb"
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
	"gorm.io/gorm"
)

type Dao[T any] struct{}

func (*Dao[T]) GetExistsByIds(
	tx *gorm.DB,
	ids []sdktypes.ID,
	locking bool,
) bool {
	var t T
	var count int64
	query := tx.Model(&t)
	query.Where("id in ?", ids)
	if locking {
		query.Clauses(sdkdb.LockingForUpdate())
	}
	err := query.Count(&count).Error
	sdk.AssertError(err)
	return int64(len(ids)) == count
}

func (*Dao[T]) GetByIds(
	tx *gorm.DB,
	ids []sdktypes.ID,
	locking bool,
) []*T {
	var t T
	var rows []*T
	query := tx.Model(&t)
	query.Where("id in ?", ids)
	if locking {
		query.Clauses(sdkdb.LockingForUpdate())
	}
	err := query.Find(&rows).Error
	sdk.AssertError(err)
	if err != nil {
		return nil
	}
	return rows
}

func (*Dao[T]) DeleteByIds(
	tx *gorm.DB,
	ids []sdktypes.ID,
	code sdkcodes.Code,
) int64 {
	var t T
	result := tx.Delete(&t, ids)
	sdk.AssertError(result.Error, code)
	return result.RowsAffected
}

func (*Dao[T]) Delete(
	tx *gorm.DB,
	condition func(tx *gorm.DB),
	code sdkcodes.Code,
) int64 {
	var t T
	query := tx.Model(&t)
	condition(query)
	result := tx.Delete(&t)
	sdk.AssertError(result.Error, code)
	return result.RowsAffected
}

func (*Dao[T]) ChangeByIds(
	tx *gorm.DB,
	updates map[string]any,
	ids sdktypes.ID,
	code sdkcodes.Code,
) int64 {
	var t T
	query := tx.Model(&t)
	query.Where("id in ?", ids)
	result := query.UpdateColumns(updates)
	sdk.AssertError(result.Error, code)
	return result.RowsAffected
}

func (*Dao[T]) Change(
	tx *gorm.DB,
	updates map[string]any,
	condition func(tx *gorm.DB),
	code sdkcodes.Code,
) int64 {
	var t T
	query := tx.Model(&t)
	condition(query)
	result := query.UpdateColumns(updates)
	sdk.AssertError(result.Error, code)
	return result.RowsAffected
}

func (*Dao[T]) Save(
	tx *gorm.DB,
	t *T,
	code sdkcodes.Code,
) {
	result := tx.Create(t)
	sdk.AssertError(result.Error, code)
	sdk.Assert(result.RowsAffected == 1, code)
}
