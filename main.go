package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/yolossn/query2metric/pkg/config"
	"github.com/yolossn/query2metric/pkg/scheduler"
)

func main() {
	conf, err := config.FromFile("config.yaml")
	if err != nil {
		log.Fatal("Error reading config file,err:", err)
	}

	configRunner := scheduler.FromConfig(*conf)
	err = configRunner.Start()
	if err != nil {
		log.Fatal("Error in runner,err:", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8090", nil))

}
