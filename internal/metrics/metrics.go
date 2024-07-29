package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	OrdersIssued = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_issued_total",
			Help: "Total number of orders issued.",
		},
	)

	CachedOrdersIssued = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cached_orders_issued_total",
			Help: "Total number of cached orders issued.",
		})

	ChangesOrdersInCache = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "changes_orders_in_cache_total",
			Help: "Total number of changes in cache orders",
		})
)

func CountMetrics() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(OrdersIssued)
	reg.MustRegister(CachedOrdersIssued)
	reg.MustRegister(ChangesOrdersInCache)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8081", nil))

}
