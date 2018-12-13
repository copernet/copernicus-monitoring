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

type GetMiningInfoCollector struct {
	blocks                  *prometheus.Desc
	currentblocksize        *prometheus.Desc
	currentblocktx          *prometheus.Desc
	difficulty              *prometheus.Desc
	blockprioritypercentage *prometheus.Desc
	networkhashps           *prometheus.Desc
	pooledtx                *prometheus.Desc
}

func NewGetMiningInfoCollector() *GetMiningInfoCollector {
	return &GetMiningInfoCollector{
		blocks:                  newDesc("get_mining_info_blocks", "Get the number of current blocks"),
		currentblocksize:        newDesc("get_mining_info_currentblocksize", "Get current block size"),
		currentblocktx:          newDesc("get_mining_info_currentblocktx", "Get current block tx size"),
		difficulty:              newDesc("get_mining_info_difficulty", "Get current block difficulty"),
		blockprioritypercentage: newDesc("get_mining_info_blockprioritypercentage", "Get current block prioritypercentage"),
		networkhashps:           newDesc("get_mining_info_networkhashps", "Get the network hashps"),
		pooledtx:                newDesc("get_mining_info_pooledtx", "Get the pooledtx"),
	}
}

func (collector *GetMiningInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blocks
	ch <- collector.currentblocksize
	ch <- collector.currentblocktx
	ch <- collector.difficulty
	ch <- collector.blockprioritypercentage
	ch <- collector.networkhashps
	ch <- collector.pooledtx
}

func (collector *GetMiningInfoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "$GOPATH/bin/coperctl getmininginfo")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err := cmd.Start()
	checkErr(err)
	err = cmd.Wait()

	var ret btcjson.GetMiningInfoResult
	if err := json.NewDecoder(stdout).Decode(&ret); err != nil {
		logs.Error(err)
		return
	}

	blocks := ret.Blocks
	currentblocksize := ret.CurrentBlockSize
	currentblocktx := ret.CurrentBlockTx
	difficulty := ret.Difficulty
	blockprioritypercentage := ret.BlockPriorityPercentage
	networkhashps := ret.NetworkHashPS
	pooledtx := ret.PooledTx

	ch <- prometheus.MustNewConstMetric(collector.blocks, prometheus.CounterValue, float64(blocks))
	ch <- prometheus.MustNewConstMetric(collector.currentblocksize, prometheus.CounterValue, float64(currentblocksize))
	ch <- prometheus.MustNewConstMetric(collector.currentblocktx, prometheus.CounterValue, float64(currentblocktx))
	ch <- prometheus.MustNewConstMetric(collector.difficulty, prometheus.CounterValue, float64(difficulty))
	ch <- prometheus.MustNewConstMetric(collector.blockprioritypercentage, prometheus.CounterValue, float64(blockprioritypercentage))
	ch <- prometheus.MustNewConstMetric(collector.networkhashps, prometheus.CounterValue, float64(networkhashps))
	ch <- prometheus.MustNewConstMetric(collector.pooledtx, prometheus.CounterValue, float64(pooledtx))
}
