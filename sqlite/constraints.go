package sqlite_assertions

// IsConstraintCheck reports whether err is SQLITE_CONSTRAINT_CHECK.
// This occurs when a CHECK constraint expression evaluates to false
// during INSERT or UPDATE.
func IsConstraintCheck(err error) bool {
	return validate(err, 275)
}

// IsConstraintPrimaryKey reports whether err is SQLITE_CONSTRAINT_PRIMARYKEY.
// This occurs when an INSERT or UPDATE violates a PRIMARY KEY constraint.
func IsConstraintPrimaryKey(err error) bool {
	return validate(err, 1555)
}

// IsConstraintForeignKey reports whether err is SQLITE_CONSTRAINT_FOREIGNKEY.
// This occurs when a FOREIGN KEY constraint fails, such as inserting or
// updating a row that references a non-existent parent row.
func IsConstraintForeignKey(err error) bool {
	return validate(err, 787)
}

// IsConstraintUnique reports whether err is SQLITE_CONSTRAINT_UNIQUE.
// This occurs when a UNIQUE constraint is violated by inserting or updating
// a row whose value duplicates another row.
func IsConstraintUnique(err error) bool {
	return validate(err, 2067)
}

// IsConstraintRowId reports whether err is SQLITE_CONSTRAINT_ROWID.
// This occurs when an operation would create a duplicate rowid.
func IsConstraintRowId(err error) bool {
	return validate(err, 2579)
}

// IsConstraintNotNull reports whether err is SQLITE_CONSTRAINT_NOTNULL.
// This occurs when a NULL value is written to a column declared NOT NULL.
func IsConstraintNotNull(err error) bool {
	return validate(err, 1299)
}

// IsConstraintDataType reports whether err is SQLITE_CONSTRAINT_DATATYPE.
// This occurs in STRICT tables when a value is incompatible with the
// column's declared type.
func IsConstraintDataType(err error) bool {
	return validate(err, 3091)
}
