package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	prometheus.MustRegister(newChainCollector())

	http.Handle("/metrics", promhttp.Handler())
	logs.Info("Beginning to serve on port :8081")
	logs.Error(http.ListenAndServe(":8081", nil))
}
