package archiver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var messagesProcessedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "scraping_messages_total",
	Help: "Number of messages processed",
}, []string{"channel_type", "channel_id"})

var discordRequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "scraping_discord_requests_total",
	Help: "Number of requests sent",
}, []string{"code", "method"})
var discordRequestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "scraping_discord_request_duration_seconds",
	Help:    "Discord request duration histogram",
	Buckets: prometheus.DefBuckets,
}, []string{"method"})
