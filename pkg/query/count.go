package query

import (
	"github.com/pkg/errors"
	"github.com/yolossn/query2metric/pkg/config"
)

type CountQuery interface {
	Count(metric config.Metric) (int64, error)
}

var ENV_NOT_SET = errors.New("connnectionString is empty")
