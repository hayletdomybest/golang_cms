package gorm

import (
	"strings"
	"weex_admin/internal/shared/errors"

	"github.com/go-sql-driver/mysql"
)

func IsNotFoundErr(err error) bool {
	return err != nil && err.Error() == "record not found"
}

type GormError struct {
	Number int
}

func TryWrapErr(err error, msgs ...string) *errors.Error {
	if err == nil {
		return nil
	}
	var msg string
	if len(msgs) > 0 {
		msg = strings.Join(msgs, " ")
	} else {
		msg = err.Error()
	}

	switch err.(*mysql.MySQLError).Number {
	case 1062:
		return errors.ErrAlreadyExist.Wrap(err)
	}

	if IsNotFoundErr(err) {
		return errors.ErrNotFound.
			Wrap(err).
			WithMessage(msg)
	}

	return errors.ErrInternal.
		Wrap(err)
}
