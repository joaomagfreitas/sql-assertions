package sqlite3_assertions

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

func validate(err error, primary sqlite3.ErrNo, extended sqlite3.ErrNoExtended) bool {
	var terr sqlite3.Error
	if !errors.As(err, &terr) {
		return false
	}

	return terr.Code == primary && terr.ExtendedCode == extended
}
