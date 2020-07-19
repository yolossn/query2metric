package main

import "github.com/yolossn/query2metric/pkg/config"

type Connection interface {
	Count(metric config.Metric) (int64, error)
}
