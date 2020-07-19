package query

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
	"github.com/yolossn/query2metric/pkg/config"
)

type sqlQuery struct {
	connection string
	db         *sql.DB
}

func NewSQLQuery(connnectionURL string) (CountQuery, error) {
	connnectionString := os.Getenv(connnectionURL)
	if connnectionString == "" {
		return nil, errors.New("connnectionString is empty")
	}
	db, err := dburl.Open(connnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "Error in establishing connection to db")
	}

	return &sqlQuery{connnectionURL, db}, nil
}

func (s *sqlQuery) Count(metric config.Metric) (int64, error) {
	if metric.Query == "" {
		return 0, errors.New("Query is empty")
	}

	countQuery := fmt.Sprintf("select count(*) from (%s) as count_query", metric.Query)
	row, err := s.db.Query(countQuery)
	if err != nil {
		return 0, errors.Wrap(err, "Error running query")
	}
	defer row.Close()
	var out int64
	if row.Next() {
		err := row.Scan(&out)
		if err != nil {
			return 0, errors.Wrap(err, "Error running query")
		}
	}
	return out, nil
}
