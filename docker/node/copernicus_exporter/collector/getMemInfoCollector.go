package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/copernet/copernicus/rpc/btcjson"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
	"time"
)

type GetMemPoolInfoCollector struct {
	memSize       *prometheus.Desc
	memUsage      *prometheus.Desc
	maxMemPool    *prometheus.Desc
	memPoolMinFee *prometheus.Desc
}

func NewGetMemPoolInfoCollector() *GetMemPoolInfoCollector {
	return &GetMemPoolInfoCollector{
		memSize:       newDesc("get_mem_pool_info_mem_size", "Get the size of current memPool"),
		memUsage:      newDesc("get_mem_pool_info_mem_usage", "Get the usage of current memPool"),
		maxMemPool:    newDesc("get_mem_pool_info_max_mem_pool", "Get the max size of current memPool"),
		memPoolMinFee: newDesc("get_mem_pool_info_mem_pool_min_fee", "Get the min fee of current memPool"),
	}
}

func (collector *GetMemPoolInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.memSize
	ch <- collector.memUsage
	ch <- collector.maxMemPool
	ch <- collector.memPoolMinFee
}

func (collector *GetMemPoolInfoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "$GOPATH/bin/coperctl getmempoolinfo")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err := cmd.Start()
	checkErr(err)
	err = cmd.Wait()

	var ret btcjson.GetMempoolInfoResult
	if err := json.NewDecoder(stdout).Decode(&ret); err != nil {
		logs.Error(err)
		return
	}

	memSize := ret.Size
	memUsage := ret.Usage
	maxMemPool := ret.MaxMempool
	memPoolMinFee := ret.MempoolMinFee

	ch <- prometheus.MustNewConstMetric(collector.memSize, prometheus.CounterValue, float64(memSize))
	ch <- prometheus.MustNewConstMetric(collector.memUsage, prometheus.CounterValue, float64(memUsage))
	ch <- prometheus.MustNewConstMetric(collector.maxMemPool, prometheus.CounterValue, float64(maxMemPool))
	ch <- prometheus.MustNewConstMetric(collector.memPoolMinFee, prometheus.CounterValue, float64(memPoolMinFee))
}
