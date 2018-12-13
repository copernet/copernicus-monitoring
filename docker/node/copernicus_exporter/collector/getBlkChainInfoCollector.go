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

type GetBlkChainInfoCollector struct {
	//bestHash   *prometheus.Desc
	//chain      *prometheus.Desc
	headers    *prometheus.Desc
	difficulty *prometheus.Desc
	//chainWork  *prometheus.Desc
}

func NewGetBlkChainInfoCollector() *GetBlkChainInfoCollector {
	return &GetBlkChainInfoCollector{
		//bestHash:   newDesc("get_block_chain_info_best_hash", "Get the best hash of current blocks"),
		//chain:      newDesc("get_block_chain_info_chain", "Whether the current network belongs to the main network or the test network"),
		headers:    newDesc("get_block_chain_info_headers", "Get the best hash of current headers"),
		difficulty: newDesc("get_block_chain_info_difficulty", "Get the current difficulty of the whole network"),
		//chainWork:  newDesc("get_block_chain_info_chain_work", "Get the current chain work of the whole network"),
	}
}

func (collector *GetBlkChainInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- collector.bestHash
	//ch <- collector.chain
	ch <- collector.headers
	ch <- collector.difficulty
	//ch <- collector.chainWork
}

func (collector *GetBlkChainInfoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "$GOPATH/bin/coperctl getblockchaininfo")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err := cmd.Start()
	checkErr(err)
	err = cmd.Wait()

	var ret btcjson.GetBlockChainInfoResult
	if err := json.NewDecoder(stdout).Decode(&ret); err != nil {
		logs.Error(err)
		return
	}

	headers := ret.Headers
	difficulty := ret.Difficulty

	ch <- prometheus.MustNewConstMetric(collector.headers, prometheus.CounterValue, float64(headers))
	ch <- prometheus.MustNewConstMetric(collector.difficulty, prometheus.CounterValue, float64(difficulty))
}
