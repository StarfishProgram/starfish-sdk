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

func (d *Dao[T]) GetExistsById(
	tx *gorm.DB,
	id sdktypes.ID,
	locking bool,
) bool {
	return d.GetExistsByIds(tx, []sdktypes.ID{id}, locking)
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

func (d *Dao[T]) GetById(
	tx *gorm.DB,
	id sdktypes.ID,
	locking bool,
) *T {
	rows := d.GetByIds(tx, []sdktypes.ID{id}, locking)
	if len(rows) > 0 {
		return rows[0]
	}
	return nil
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

func (d *Dao[T]) DeleteById(
	tx *gorm.DB,
	id sdktypes.ID,
	code sdkcodes.Code,
) int64 {
	return d.DeleteByIds(tx, []sdktypes.ID{id}, code)
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
	ids []sdktypes.ID,
	code sdkcodes.Code,
) int64 {
	var t T
	query := tx.Model(&t)
	query.Where("id in ?", ids)
	result := query.UpdateColumns(updates)
	sdk.AssertError(result.Error, code)
	return result.RowsAffected
}

func (d *Dao[T]) ChangeById(
	tx *gorm.DB,
	updates map[string]any,
	id sdktypes.ID,
	code sdkcodes.Code,
) int64 {
	return d.ChangeByIds(tx, updates, []sdktypes.ID{id}, code)
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

func (*Dao[T]) SaveInBatch(
	tx *gorm.DB,
	ts []*T,
	code sdkcodes.Code,
) {
	result := tx.CreateInBatches(&ts, 1000)
	sdk.AssertError(result.Error, code)
}
