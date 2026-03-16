package sqlite_assertions

import (
	"errors"

	"modernc.org/sqlite"
)

func validate(err error, code int) bool {
	if terr, ok := errors.AsType[*sqlite.Error](err); ok {
		return terr.Code() == code
	}

	return false
}
