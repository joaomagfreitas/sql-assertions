package pq_assertions

import "github.com/lib/pq/pqerror"

// IsConstraintCheck reports whether err is a CHECK violation (SQLSTATE 23514).
// This occurs when a CHECK constraint expression evaluates to false
// during INSERT or UPDATE.
func IsConstraintCheck(err error) bool {
	return validate(err, pqerror.CheckViolation)
}

// IsConstraintForeignKey reports whether err is a FOREIGN KEY violation (SQLSTATE 23503).
// This occurs when a FOREIGN KEY constraint fails, such as inserting or
// updating a row that references a non-existent parent row.
func IsConstraintForeignKey(err error) bool {
	return validate(err, pqerror.ForeignKeyViolation)
}

// IsConstraintUnique reports whether err is a UNIQUE KEY violation (SQLSTATE 23505).
// This occurs when a UNIQUE (primary or unique index) constraint is violated by inserting or updating
// a row whose value duplicates another row.
func IsConstraintUnique(err error) bool {
	return validate(err, pqerror.UniqueViolation)
}

// IsConstraintNotNull reports whether err is a NOT NULL violation (SQLSTATE 23502).
// This occurs when a NULL value is written to a column declared NOT NULL.
func IsConstraintNotNull(err error) bool {
	return validate(err, pqerror.NotNullViolation)
}

// IsConstraintDataType reports whether err is a datatype mismatch error (SQLSTATE 42804).
// This occurs when a value is incompatible with the column's declared type.
func IsConstraintDataType(err error) bool {
	return validate(err, pqerror.DatatypeMismatch)
}
