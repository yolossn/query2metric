package query_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xo/dburl"
	"github.com/yolossn/query2metric/pkg/config"
	"github.com/yolossn/query2metric/pkg/query"
)

func TestNewSQLConn(t *testing.T) {
	t.Parallel()

	// Env check
	connectionEnv := "TEST_INVALID_SQL_CONN"

	conn, err := query.NewSQLQuery(connectionEnv)
	require.Error(t, err)
	require.Equal(t, err, query.ENV_NOT_SET)
	require.Nil(t, conn)

	// Invalid URI
	connectionStr := "postg://127.0.0.1/sdffsdf"
	conn, err = newSQLConnectionHelper(t, connectionEnv, connectionStr)
	require.Error(t, err)
	require.Nil(t, conn)

	// Valid URI
	connectionStr = os.Getenv("TEST_SQL_CONN")
	conn, err = newSQLConnectionHelper(t, connectionEnv, connectionStr)
	require.NoError(t, err)
	require.NotNil(t, conn)

}

func TestSQLCount(t *testing.T) {

	// Setup
	connectionEnv := "TEST_SQL_CONN"
	connectionStr := os.Getenv(connectionEnv)
	db, err := dburl.Open(connectionStr)
	require.NotNil(t, db)
	require.NoError(t, err)

	// Clean
	dropTableStmt := `DROP TABLE IF EXISTS test`
	res, err := db.Exec(dropTableStmt)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Test Data
	createTableStmt := `CREATE TABLE test(name TEXT, is_active BOOLEAN)`
	_, err = db.Exec(createTableStmt)
	require.NoError(t, err)

	insertStmt := `INSERT INTO test (name, is_active) VALUES ($1, $2)`
	_, err = db.Exec(insertStmt, "test1", true)
	require.NoError(t, err)
	_, err = db.Exec(insertStmt, "test2", false)
	require.NoError(t, err)

	// New Conn
	conn, err := newSQLConnectionHelper(t, connectionEnv, connectionStr)
	require.NoError(t, err)
	require.NotNil(t, conn)

	// Test
	// Positive test cases
	testQuery := config.Metric{
		Query: `select * from test where is_active = true`,
	}
	out, err := conn.Count(testQuery)
	require.Equal(t, int64(1), out)
	require.NoError(t, err)

	testQuery.Query = `select * from test`
	out, err = conn.Count(testQuery)
	require.Equal(t, int64(2), out)
	require.NoError(t, err)

	testQuery.Query = ""
	out, err = conn.Count(testQuery)
	require.Error(t, err)
	require.Equal(t, int64(0), out)

	// Negative test cases
	testQuery.Query = `test`
	out, err = conn.Count(testQuery)
	require.Error(t, err)
	require.Equal(t, int64(0), out)

}

func newSQLConnectionHelper(t *testing.T, connectionEnv, connectionString string) (query.CountQuery, error) {
	t.Helper()

	err := os.Setenv(connectionEnv, connectionString)
	require.NoError(t, err)

	return query.NewSQLQuery(connectionEnv)
}
