package scheduler

import (
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/yolossn/query2metric/pkg/config"
	"github.com/yolossn/query2metric/pkg/query"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

type Scheduler struct {
	conf config.Config
}

func (s Scheduler) Start() error {
	errorChan := make(chan bool, 1)
	successChan := make(chan bool, 1)
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
			err = prometheus.Register(gaugeMetric)
			if err != nil {
				return errors.Wrap(err, "Error registering metric")
			}
			ticker := time.NewTicker(time.Duration(metric.Time) * time.Second)
			run(ticker, gaugeMetric, dbConnection, metric, successChan, errorChan)
		}
	}

	go errorCounter(errorChan)
	go successCounter(successChan)

	return nil
}

func run(tick *time.Ticker, gauge prometheus.Gauge, quer query.CountQuery, metric config.Metric, successChan, errorChan chan bool) {

	go func() {
		for {
			select {
			case <-tick.C:
				out, err := quer.Count(metric)
				if err != nil {
					errorChan <- true
					log.WithFields(log.Fields{"db": metric.Database, "metric": metric.Name, "query": metric.Query}).Error(err)
				} else {
					gauge.Set(float64(out))
					successChan <- true
					log.WithFields(log.Fields{"db": metric.Database, "metric": metric.Name, "query": metric.Query}).Debug("success")
				}
			}
		}
	}()
}

func FromConfig(conf config.Config) Scheduler {
	return Scheduler{conf}
}

func errorCounter(errorChan chan bool) {

	errorCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "query2metric",
			Name:      "error_count",
			Help:      "No of errors when converting query to metric",
		},
	)

	prometheus.MustRegister(errorCounter)
	for {
		switch {
		case <-errorChan:
			errorCounter.Inc()
		}
	}

}

func successCounter(successChan chan bool) {

	successCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "query2metric",
			Name:      "success_count",
			Help:      "No of successful queries coverted to metrics",
		},
	)

	prometheus.MustRegister(successCounter)

	for {
		switch {
		case <-successChan:
			successCounter.Inc()
		}
	}

}
