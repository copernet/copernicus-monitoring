package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/copernet/copernicus-monitoring/docker/node/copernicus_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	prometheus.MustRegister(collector.NewGetInfoCollector())
	prometheus.MustRegister(collector.NewGetBlkChainInfoCollector())
	prometheus.MustRegister(collector.NewGetMemPoolInfoCollector())
	prometheus.MustRegister(collector.NewGetMiningInfoCollector())
	prometheus.MustRegister(collector.NewGetTxOutSetInfoCollector())

	http.Handle("/metrics", promhttp.Handler())
	logs.Info("Beginning to serve on port :8081")
	logs.Error(http.ListenAndServe(":8081", nil))
}
