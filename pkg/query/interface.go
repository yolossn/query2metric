package query

import "github.com/yolossn/query2metric/pkg/config"

type CountQuery interface {
	Count(metric config.Metric) (int64, error)
}
