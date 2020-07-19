package scheduler

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yolossn/query2metric/pkg/config"
	"github.com/yolossn/query2metric/pkg/query"
)

type Scheduler struct {
	conf config.Config
}

func (s Scheduler) Start() error {
	for _, conn := range s.conf.Connections {
		var dbConnection query.CountQuery
		var err error
		switch conn.Type {
		case config.MONGO:
			dbConnection, err = query.NewMongoConn(conn.ConnectionString)
			if err != nil {
				return err
			}
		case config.SQL:
			dbConnection, err = query.NewSQLQuery(conn.ConnectionString)
			if err != nil {
				return err
			}
		default:
			continue
		}

		for _, metric := range conn.Metrics {

			gaugeMetric := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: conn.Name,
					Name:      metric.Name,
					Help:      metric.HelpString,
				},
			)
			prometheus.MustRegister(gaugeMetric)
			ticker := time.NewTicker(time.Duration(metric.Time) * time.Second)
			run(ticker, gaugeMetric, dbConnection, metric)
		}
	}
	return nil
}

func run(tick *time.Ticker, gauge prometheus.Gauge, quer query.CountQuery, metric config.Metric) {

	go func() {
		for {
			select {
			case <-tick.C:
				out, err := quer.Count(metric)
				if err != nil {
					fmt.Println(err)
				} else {
					gauge.Set(float64(out))
				}
			}
		}
	}()
}

func FromConfig(conf config.Config) Scheduler {
	return Scheduler{conf}
}
