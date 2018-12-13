package collector

import (
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/client_golang/prometheus"
)

func newDesc(fqName, help string) *prometheus.Desc {
	return prometheus.NewDesc(fqName, help, nil, nil)
}

func checkErr(err error) {
	if err != nil {
		logs.Error("error:%s", err.Error())
	}
}

