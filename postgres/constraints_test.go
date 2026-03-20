package pq_assertions_test

import (
	"database/sql"
	"testing"

	pq_assertions "github.com/joaomagfreitas/sql-assertions/postgres"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var constraintsSchema = `
-- CHECK constraint
DROP TABLE IF EXISTS test_check;
CREATE TABLE test_check (
    v INTEGER CHECK (v > 0)
);

-- PRIMARY KEY constraint
DROP TABLE IF EXISTS test_primary_key;
CREATE TABLE test_primary_key (
    id INTEGER PRIMARY KEY
);

-- UNIQUE constraint
DROP TABLE IF EXISTS test_unique;
CREATE TABLE test_unique (
    v TEXT UNIQUE
);

-- NOT NULL constraint
DROP TABLE IF EXISTS test_not_null;
CREATE TABLE test_not_null (
    v TEXT NOT NULL
);

-- FOREIGN KEY constraint
DROP TABLE IF EXISTS test_foreign_key;
DROP TABLE IF EXISTS test_parent;

CREATE TABLE test_parent (
    id INTEGER PRIMARY KEY
);

CREATE TABLE test_foreign_key (
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES test_parent(id)
);

-- DATATYPE constraint
DROP TABLE IF EXISTS test_datatype;
CREATE TABLE test_datatype (
    v INTEGER
);
`

func TestConstraints(t *testing.T) {
	ct, err := postgres.Run(t.Context(),
		"postgres:18-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer testcontainers.TerminateContainer(ct)

	conn, err := ct.ConnectionString(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", conn)
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
			asserter: pq_assertions.IsConstraintCheck,
		},
		{
			desc: "foreign key",
			dml: func() error {
				_, err = db.Exec(`
					INSERT INTO test_foreign_key VALUES (999);
				`)
				return err
			},
			asserter: pq_assertions.IsConstraintForeignKey,
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
			asserter: pq_assertions.IsConstraintUnique,
		},
		{
			desc: "not null",
			dml: func() error {
				_, err = db.Exec(`
					INSERT INTO test_not_null VALUES (NULL);
				`)
				return err
			},
			asserter: pq_assertions.IsConstraintNotNull,
		},
		{
			desc: "datatype",
			dml: func() error {
				_, err = db.Exec(`
					INSERT INTO test_datatype VALUES (FALSE);
				`)
				return err
			},
			asserter: pq_assertions.IsConstraintDataType,
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
