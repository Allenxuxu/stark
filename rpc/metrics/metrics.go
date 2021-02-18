package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultMetricsPath = "/metrics"

var (
	rg = prometheus.NewRegistry()
)

func init() {
	rg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	rg.MustRegister(prometheus.NewGoCollector())
}

func PrometheusMustRegister(cs ...prometheus.Collector) {
	rg.MustRegister(cs...)
}

func Run(path, address string) error {
	if len(path) == 0 {
		path = defaultMetricsPath
	}

	http.Handle(path, promhttp.HandlerFor(rg, promhttp.HandlerOpts{}))
	return http.ListenAndServe(address, nil)
}
