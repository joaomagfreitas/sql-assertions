package pq_assertions

import (
	"errors"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

func validate(err error, code pqerror.Code) bool {
	if terr, ok := errors.AsType[*pq.Error](err); ok {
		return terr.Code == code
	}

	return false
}
