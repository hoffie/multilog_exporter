package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func runServer(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Use /metrics\n"))
	})
	http.Handle("/metrics", promhttp.Handler())
	log.WithFields(log.Fields{"addr": *metricsListenAddr}).Info("Serving metrics")
	log.Fatal(http.ListenAndServe(*metricsListenAddr, nil))
}
