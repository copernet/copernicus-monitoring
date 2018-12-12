package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
)

type ChainCollector struct {
	blockHeight  *prometheus.Desc
	headerHeight *prometheus.Desc
}

func newDesc(fqName, help string) *prometheus.Desc {
	return prometheus.NewDesc(fqName, help, nil, nil)
}

func newChainCollector() *ChainCollector {
	return &ChainCollector{
		blockHeight: newDesc("c_chain_active_block_height", "Active chain's block height"),
		headerHeight: newDesc("c_chain_active_header_height","Active chain's header height"),
	}
}

func (collector *ChainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blockHeight
	ch <- collector.headerHeight
}

func (collector *ChainCollector) Collect(ch chan<- prometheus.Metric) {
	blockHeight := rand.Float64()
	headerHeight := rand.Float64()

	ch <- prometheus.MustNewConstMetric(collector.blockHeight, prometheus.CounterValue, blockHeight)
	ch <- prometheus.MustNewConstMetric(collector.headerHeight, prometheus.CounterValue, headerHeight)
}
