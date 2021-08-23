package integration_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/newqf"
	qsql "github.com/tobgu/qframe/config/sql"
)

const (
	pg_host        = "localhost"
	pg_port        = 5432
	pg_user        = "qframe"
	pg_password    = "qframe"
	pg_dbname      = "qframe"
	mysql_hostname = "localhost:3306"
	mysql_user     = "qframe"
	mysql_password = "qframe"
	mysql_dbname   = "qframe"
)

func GetPostgres(t *testing.T) *sql.DB {
	t.Helper()
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_password, pg_dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		t.Error(err)
	}
	err = db.Ping()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return db
}

func GetMySQL(t *testing.T) *sql.DB {
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysql_user, mysql_password, mysql_hostname, mysql_dbname)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		t.Error(err)
	}
	err = db.Ping()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return db
}

func GetSQLite(t *testing.T) *sql.DB {
	db, _ := sql.Open("sqlite3", ":memory:")
	db.Exec(`
		CREATE TABLE IF NOT EXISTS mockdata
		(
			int_col integer NOT NULL, 
			float_col double precision,
			string_col text NOT NULL,
			bool_col boolean,
			PRIMARY KEY (int_col)
		);
		CREATE INDEX IF NOT EXISTS mock_data_string_index ON mockdata (string_col);`)
	return db
}

type DBTestConfig struct {
	name             string
	getConn          func(t *testing.T) *sql.DB
	paramSymbol      string
	additionalConfig []qsql.ConfigFunc
}

// Test that we can select rows using supplied args
func TestReadSQLWithArgs(t *testing.T) {
	tcs := []DBTestConfig{
		{"postgres", GetPostgres, "$1", []qsql.ConfigFunc{}},
		{"mysql", GetMySQL, "?", []qsql.ConfigFunc{qsql.Coerce(qsql.CoercePair{Column: "bool_col", Type: qsql.Int64ToBool})}},
		{"sqlite", GetSQLite, "?", []qsql.ConfigFunc{}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			db := tc.getConn(t)
			paramSymbol := tc.paramSymbol

			tx, err := db.Begin()
			if err != nil {
				t.Error(err)
			}

			ins := `INSERT INTO mockdata (int_col, float_col, string_col, bool_col) VALUES
			(1, 1.1, 'one', true),
			(2, 2.2, 'two', true),
			(3, 3.3, 'three', false)
		`
			stmt, err := tx.Prepare(ins)
			if err != nil {
				t.Error(err)
			}
			stmt.Exec()

			// data inserted - now retrieve using qframe

			// for postgres this resolves to:
			// SELECT int_col, float_col, string_col, bool_col FROM mockdata WHERE string_col=$1
			// for mysql this will be:
			// SELECT int_col, float_col, string_col, bool_col FROM mockdata WHERE string_col=?
			sel := fmt.Sprintf("SELECT int_col, float_col, string_col, bool_col FROM mockdata WHERE string_col=%v", paramSymbol)
			configs := append(tc.additionalConfig, qsql.Query(sel))
			qf := qframe.ReadSQLWithArgs(tx, []interface{}{"two"}, configs...)

			data := map[string]interface{}{
				"int_col":    []int{2},
				"float_col":  []float64{2.2},
				"string_col": []string{"two"},
				"bool_col":   []bool{true}}
			inOrder := newqf.ColumnOrder("int_col", "float_col", "string_col", "bool_col")
			expected := qframe.New(data, inOrder)
			eq, reason := qf.Equals(expected)
			if !eq {
				t.Errorf("\nexpected:\n %s\n got: \n %s\nreason: %s", expected.String(), qf.String(), reason)
			}
			tx.Rollback()
			db.Close()
		})
	}
}
