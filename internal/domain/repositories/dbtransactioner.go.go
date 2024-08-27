package repositories

import "context"

type ContextOption struct {
	Context context.Context
}

type TransactionOption struct {
	Tx interface{}
}

type DBTransactioner interface {
	Begin() (interface{}, error)
	Rollback(interface{}) error
	Commit(interface{}) error
}
