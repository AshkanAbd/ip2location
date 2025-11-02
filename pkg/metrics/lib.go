package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var httpStatusMetric *prometheus.CounterVec
var countryMetric *prometheus.CounterVec
var ipSourceMetric *prometheus.CounterVec

var registry = prometheus.NewRegistry()

func RegisterMetrics() {
	httpStatusMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_status_count",
		Help: "The total number of handled http requests per status",
	}, []string{"status"})
	countryMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "country_count",
		Help: "The total number of ip check per country",
	}, []string{"country"})
	ipSourceMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "hit_rate",
		Help: "The total number of ip check from different sources",
	}, []string{"source"})

	registry.MustRegister(httpStatusMetric)
	registry.MustRegister(countryMetric)
	registry.MustRegister(ipSourceMetric)
}

func IncreaseHttpStatus(status string) {
	httpStatusMetric.With(prometheus.Labels{"status": status}).Inc()
}

func IncreaseCountryCount(country string) {
	countryMetric.With(prometheus.Labels{"country": country}).Inc()
}

func IncreaseIPSourceCount(source string) {
	ipSourceMetric.With(prometheus.Labels{"source": source}).Inc()
}

func GetRegistry() *prometheus.Registry {
	return registry
}
