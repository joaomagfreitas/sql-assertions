package sqlite_assertions_test

import (
	"database/sql"
	"testing"

	sqlite_assertions "github.com/joaomagfreitas/sql-assertions/sqlite"
	_ "modernc.org/sqlite"
)

var constraintsSchema = `
PRAGMA foreign_keys = ON;

-- CHECK constraint
CREATE TABLE test_check (
    v INTEGER CHECK (v > 0)
);

-- PRIMARY KEY constraint
CREATE TABLE test_primary_key (
    id INTEGER PRIMARY KEY
);

-- UNIQUE constraint
CREATE TABLE test_unique (
    v TEXT UNIQUE
);

-- NOT NULL constraint
CREATE TABLE test_not_null (
    v TEXT NOT NULL
);

-- ROWID constraint (explicit rowid usage)
CREATE TABLE test_rowid (
    v TEXT
);

-- FOREIGN KEY constraint
CREATE TABLE test_parent (
    id INTEGER PRIMARY KEY
);

CREATE TABLE test_foreign_key (
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES test_parent(id)
);

-- DATATYPE constraint (requires STRICT tables)
CREATE TABLE test_datatype (
    v INTEGER
) STRICT;
`

func TestConstraints(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(constraintsSchema)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		dml      func() error
		asserter func(err error) bool
		desc     string
	}{
		{
			desc: "check",
			dml: func() error {
				_, err = db.Exec(`INSERT INTO test_check VALUES (-1)`)
				return err
			},
			asserter: sqlite_assertions.IsConstraintCheck,
		},
		{
			desc: "primary key",
			dml: func() error {
				_, err = db.Exec(`
				INSERT INTO test_primary_key VALUES (1);
				INSERT INTO test_primary_key VALUES (1);
				`)

				return err
			},
			asserter: sqlite_assertions.IsConstraintPrimaryKey,
		},
		{
			desc: "foreign key",
			dml: func() error {
				_, err = db.Exec(`
			INSERT INTO test_foreign_key VALUES (999);
		`)
				return err
			},
			asserter: sqlite_assertions.IsConstraintForeignKey,
		},
		{
			desc: "unique",
			dml: func() error {
				_, err = db.Exec(`
			INSERT INTO test_unique VALUES ('a');
			INSERT INTO test_unique VALUES ('a');
		`)
				return err
			},
			asserter: sqlite_assertions.IsConstraintUnique,
		},
		{
			desc: "rowid",
			dml: func() error {
				_, err = db.Exec(`
			INSERT INTO test_rowid(rowid, v) VALUES (1, 'a');
			INSERT INTO test_rowid(rowid, v) VALUES (1, 'b');
		`)
				return err
			},
			asserter: sqlite_assertions.IsConstraintRowId,
		},
		{
			desc: "not null",
			dml: func() error {
				_, err = db.Exec(`
			INSERT INTO test_not_null VALUES (NULL);
		`)
				return err
			},
			asserter: sqlite_assertions.IsConstraintNotNull,
		},
		{
			desc: "datatype",
			dml: func() error {
				_, err = db.Exec(`
			INSERT INTO test_datatype VALUES ('text');
		`)
				return err
			},
			asserter: sqlite_assertions.IsConstraintDataType,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			aerr := tC.dml()
			if aerr == nil {
				t.FailNow()
			}

			if !tC.asserter(aerr) {
				t.Fatal(aerr)
			}
		})
	}
}
