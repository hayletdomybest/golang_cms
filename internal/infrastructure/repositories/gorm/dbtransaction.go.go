package gorm

import (
	"weex_admin/internal/domain/repositories"
	"weex_admin/internal/shared/errors"

	"gorm.io/gorm"
)

type DBTransaction struct {
	db *gorm.DB
}

func (t *DBTransaction) Begin() (interface{}, error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return nil, errors.ErrInternal.Wrap(tx.Error).WithMessage("failed to start transaction")
	}
	return tx, nil
}

func (t *DBTransaction) Commit(tx interface{}) error {
	txDb, err := t.getTx(tx)
	if err != nil {
		return err
	}
	if err := txDb.Commit().Error; err != nil {
		return errors.ErrInternal.
			WithMessage("TransactionManager fail to Commit").
			Wrap(err)
	}

	return nil
}

func (t *DBTransaction) Rollback(tx interface{}) error {
	txDb, err := t.getTx(tx)
	if err != nil {
		return err
	}
	return errors.ErrInternal.
		WithMessage("TransactionManager fail to Rollback").
		Wrap(txDb.Rollback().Error)
}

var _ repositories.DBTransactioner = (*DBTransaction)(nil)

func NewDBTransaction(db *gorm.DB) repositories.DBTransactioner {
	return &DBTransaction{db: db}
}

func (t *DBTransaction) getTx(tx interface{}) (*gorm.DB, error) {
	if tx == nil {
		return nil, errors.ErrPanic.WithMessage("transaction is nil")
	}

	if txDb, ok := tx.(*gorm.DB); ok {
		return txDb, nil
	} else {
		return nil, errors.ErrPanic.WithMessage("invalid transaction type")
	}
}
