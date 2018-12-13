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

type GetTxOutSetInfoCollector struct {
	height       *prometheus.Desc
	transactions *prometheus.Desc
	txouts       *prometheus.Desc
	bogosize     *prometheus.Desc
	diskSize     *prometheus.Desc
	totalAmount  *prometheus.Desc
}

func NewGetTxOutSetInfoCollector() *GetTxOutSetInfoCollector {
	return &GetTxOutSetInfoCollector{
		height:       newDesc("get_txout_set_info_height", "Get the number of current blocks"),
		transactions: newDesc("get_txout_set_info_transactions", "Get current protocol version"),
		txouts:       newDesc("get_txout_set_info_txouts", "Get the number of current connections"),
		bogosize:     newDesc("get_txout_set_info_bogosize", "Get the number of current blocks"),
		diskSize:     newDesc("get_txout_set_info_diskSize", "Get current protocol version"),
		totalAmount:  newDesc("get_txout_set_info_totalAmount", "Get the number of current connections"),
	}
}

func (collector *GetTxOutSetInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.height
	ch <- collector.transactions
	ch <- collector.txouts
	ch <- collector.bogosize
	ch <- collector.diskSize
	ch <- collector.totalAmount
}

func (collector *GetTxOutSetInfoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "$GOPATH/bin/coperctl gettxoutsetinfo")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err := cmd.Start()
	checkErr(err)
	err = cmd.Wait()
	//logs.Info("cmd wait err:%s", err.Error())

	var ret btcjson.GetTxOutSetInfoResult
	if err := json.NewDecoder(stdout).Decode(&ret); err != nil {
		logs.Error(err)
		return
	}

	height := ret.Height
	transactions := ret.Transactions
	txouts := ret.TxOuts
	bogosize:=ret.BogoSize
	disksize:=ret.DiskSize
	totalamount:=ret.TotalAmount

	ch <- prometheus.MustNewConstMetric(collector.height, prometheus.CounterValue, float64(height))
	ch <- prometheus.MustNewConstMetric(collector.transactions, prometheus.CounterValue, float64(transactions))
	ch <- prometheus.MustNewConstMetric(collector.txouts, prometheus.CounterValue, float64(txouts))
	ch <- prometheus.MustNewConstMetric(collector.bogosize, prometheus.CounterValue, float64(bogosize))
	ch <- prometheus.MustNewConstMetric(collector.diskSize, prometheus.CounterValue, float64(disksize))
	ch <- prometheus.MustNewConstMetric(collector.totalAmount, prometheus.CounterValue, float64(totalamount))
}
