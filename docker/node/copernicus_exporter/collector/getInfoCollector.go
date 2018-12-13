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

type GetInfoCollector struct {
	blocks          *prometheus.Desc
	protocolVersion *prometheus.Desc
	connections     *prometheus.Desc
}

func newDesc(fqName, help string) *prometheus.Desc {
	return prometheus.NewDesc(fqName, help, nil, nil)
}

func NewGetInfoCollector() *GetInfoCollector {
	return &GetInfoCollector{
		blocks:          newDesc("get_info_blocks", "Get the number of current blocks"),
		protocolVersion: newDesc("get_info_protocol_version", "Get current protocol version"),
		connections:     newDesc("get_info_connections", "Get the number of current connections"),
	}
}

func (collector *GetInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blocks
	ch <- collector.protocolVersion
	ch <- collector.connections
}

func (collector *GetInfoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "$GOPATH/bin/coperctl getinfo")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err := cmd.Start()
	checkErr(err)
	err = cmd.Wait()
	//logs.Info("cmd wait err:%s", err.Error())

	var ret btcjson.InfoChainResult
	if err := json.NewDecoder(stdout).Decode(&ret); err != nil {
		logs.Error(err)
		return
	}

	blocks := ret.Blocks
	protocolVersion := ret.ProtocolVersion
	connections:= ret.Connections

	ch <- prometheus.MustNewConstMetric(collector.blocks, prometheus.CounterValue, float64(blocks))
	ch <- prometheus.MustNewConstMetric(collector.protocolVersion, prometheus.CounterValue, float64(protocolVersion))
	ch <- prometheus.MustNewConstMetric(collector.connections, prometheus.CounterValue, float64(connections))
}

func checkErr(err error) {
	if err != nil {
		logs.Error("error:%s", err.Error())
	}
}
