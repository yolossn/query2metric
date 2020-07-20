package scheduler_test

import (
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
	"github.com/xo/dburl"
	"github.com/yolossn/query2metric/pkg/config"
	"github.com/yolossn/query2metric/pkg/scheduler"
)

func TestFromConfig(t *testing.T) {
	t.Parallel()
	config := config.Config{
		Connections: []config.Connection{
			{
				Name: "test2",
			},
		},
	}
	sh := scheduler.FromConfig(config)
	require.NotNil(t, sh)
}

func TestSchedulerStart(t *testing.T) {
	t.Parallel()

	conf, err := config.FromFile("../../test/test_config.yaml")
	require.NoError(t, err)
	sch := scheduler.FromConfig(*conf)
	sch.Start()

	// Setup
	connectionEnv := "TEST_SQL_CONN"
	connectionStr := os.Getenv(connectionEnv)
	db, err := dburl.Open(connectionStr)
	require.NotNil(t, db)
	require.NoError(t, err)

	dropTableStmt := `DROP TABLE IF EXISTS queries`
	res, err := db.Exec(dropTableStmt)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Test Data
	createTableStmt := `CREATE TABLE queries(name TEXT, is_active BOOLEAN)`
	_, err = db.Exec(createTableStmt)
	require.NoError(t, err)

	insertStmt := `INSERT INTO queries (name, is_active) VALUES ($1, $2)`
	_, err = db.Exec(insertStmt, "test1", true)
	require.NoError(t, err)
	_, err = db.Exec(insertStmt, "test2", false)
	require.NoError(t, err)

	// Test
	time.Sleep(time.Duration(5) * time.Second)
	out, err := prometheus.DefaultGatherer.Gather()
	require.NotNil(t, out)
	require.NoError(t, err)

	metric := filterMetrics(t, out, []string{"postgres1_template_count"})
	require.NotNil(t, metric)
	require.Equal(t, 1, len(metric))
	require.NotZero(t, metric[0].Metric[0].Gauge.Value)

	currentValue := metric[0].Metric[0].Gauge.Value
	// add more data
	_, err = db.Exec(insertStmt, "test3", false)
	require.NoError(t, err)

	// wait for the metric to be recalculated from query
	time.Sleep(time.Duration(5) * time.Second)
	out, err = prometheus.DefaultGatherer.Gather()
	require.NotNil(t, out)
	require.NoError(t, err)

	metric = filterMetrics(t, out, []string{"postgres1_template_count"})
	require.NotNil(t, metric)
	require.Equal(t, 1, len(metric))
	require.NotZero(t, metric[0].Metric[0].Gauge.Value)
	// check the value is not equal to old value
	require.NotEqual(t, *currentValue, *metric[0].Metric[0].Gauge.Value)

	// TearDown
	res, err = db.Exec(dropTableStmt)
	require.NoError(t, err)
	require.NotNil(t, res)

}

func filterMetrics(t *testing.T, metrics []*dto.MetricFamily, names []string) []*dto.MetricFamily {

	t.Helper()

	var filtered []*dto.MetricFamily
	for _, m := range metrics {
		drop := true
		for _, name := range names {
			if m.GetName() == name {
				drop = false
				break
			}
		}
		if !drop {
			filtered = append(filtered, m)
		}
	}
	return filtered
}
