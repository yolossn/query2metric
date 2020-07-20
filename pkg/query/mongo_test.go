package query_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yolossn/query2metric/pkg/config"
	"github.com/yolossn/query2metric/pkg/query"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewMongoConn(t *testing.T) {
	t.Parallel()

	// Env check
	connectionEnv := "TEST_INVALID_MONGO_CONN"

	conn, err := query.NewMongoConn(connectionEnv)
	require.Error(t, err)
	require.Equal(t, err, query.ENV_NOT_SET)
	require.Nil(t, conn)

	// invalid URI check
	connectionStr := "mongo://test@123:27017/not_exists"

	conn, err = newMongoConnectionHelper(t, connectionEnv, connectionStr)
	require.Error(t, err)
	require.Nil(t, conn)

	// valid connect check
	connectionStr = "mongodb://test@123:27017/not_exists?directConnection=blah"

	conn, err = newMongoConnectionHelper(t, connectionEnv, connectionStr)
	require.NoError(t, err)
	require.NotNil(t, conn)

}

func TestMongoCount(t *testing.T) {

	ctx := context.Background()
	connectionEnv := "TEST_MONGO_CONN"
	connectionString := os.Getenv(connectionEnv)

	// Setup

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	require.NotNil(t, mongoClient)
	require.NoError(t, err)

	err = mongoClient.Connect(ctx)
	require.NoError(t, err)

	coll := mongoClient.Database("yolo_db").Collection("test")

	// Clean
	err = coll.Drop(ctx)
	require.NoError(t, err)

	// Test Data
	_, err = coll.InsertOne(ctx, map[string]interface{}{"name": "test", "is_active": true})
	require.NoError(t, err)
	_, err = coll.InsertOne(ctx, map[string]interface{}{"name": "test1", "is_active": false})
	require.NoError(t, err)

	// New Conn
	conn, err := newMongoConnectionHelper(t, connectionEnv, connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn)

	// Test
	// Positive test cases
	testQuery := config.Metric{
		Database:   "yolo_db",
		Collection: "test",
		Query:      `{"is_active":true}`,
	}
	out, err := conn.Count(testQuery)
	require.Equal(t, int64(1), out)
	require.NoError(t, err)

	testQuery.Query = ""
	out, err = conn.Count(testQuery)
	require.Equal(t, int64(2), out)
	require.NoError(t, err)

	// negative test cases
	testQuery.Query = "{test:'invalidjson'"
	out, err = conn.Count(testQuery)
	require.Equal(t, int64(0), out)
	require.Error(t, err)

	testQuery.Query = `{"test":{"$count":"{"}}`
	out, err = conn.Count(testQuery)
	require.Equal(t, int64(0), out)
	require.Error(t, err)

	testQuery.Query = `{"test":{"$count":"{"}}`
	out, err = conn.Count(testQuery)
	require.Equal(t, int64(0), out)
	require.Error(t, err)

}

func newMongoConnectionHelper(t *testing.T, connectionEnv, connectionString string) (query.CountQuery, error) {
	t.Helper()

	err := os.Setenv(connectionEnv, connectionString)
	require.NoError(t, err)

	return query.NewMongoConn(connectionEnv)
}
