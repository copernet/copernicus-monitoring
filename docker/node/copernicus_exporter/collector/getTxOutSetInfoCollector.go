package collector

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/copernet/copernicus/rpc/btcjson"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
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
		height:       newDesc("get_txout_set_info_height", "Get the height of current txout set"),
		transactions: newDesc("get_txout_set_info_transactions", "Get the transactions of current txout set"),
		txouts:       newDesc("get_txout_set_info_txouts", "Get the txouts of current txout set"),
		bogosize:     newDesc("get_txout_set_info_bogosize", "Get the bogosize of current txout set"),
		diskSize:     newDesc("get_txout_set_info_diskSize", "Get the disk size of current txout set"),
		totalAmount:  newDesc("get_txout_set_info_totalAmount", "Get the total amount of current txout set"),
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
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()

	cmd := exec.Command("/bin/bash", "-c", "$GOPATH/bin/coperctl gettxoutsetinfo")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err := cmd.Start()
	checkErr(err)
	err = cmd.Wait()

	var ret btcjson.GetTxOutSetInfoResult

	decoder := json.NewDecoder(stdout)
	decoder.UseNumber()
	err = decoder.Decode(&ret)
	if err != nil {
		logs.Error(err)
		return
	}

	height := ret.Height
	transactions := ret.Transactions
	txouts := ret.TxOuts
	bogosize := ret.BogoSize
	disksize := ret.DiskSize
	totalamount := ret.TotalAmount

	ch <- prometheus.MustNewConstMetric(collector.height, prometheus.CounterValue, float64(height))
	ch <- prometheus.MustNewConstMetric(collector.transactions, prometheus.CounterValue, float64(transactions))
	ch <- prometheus.MustNewConstMetric(collector.txouts, prometheus.CounterValue, float64(txouts))
	ch <- prometheus.MustNewConstMetric(collector.bogosize, prometheus.CounterValue, float64(bogosize))
	ch <- prometheus.MustNewConstMetric(collector.diskSize, prometheus.CounterValue, float64(disksize))
	ch <- prometheus.MustNewConstMetric(collector.totalAmount, prometheus.CounterValue, totalamount)
}
