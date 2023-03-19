package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var summaryObjectives = map[float64]float64{
	1:     0,
	0.999: 0.0001,
	0.995: 0.001,
	0.99:  0.001,
	0.95:  0.01,
	0.9:   0.01,
	0.75:  0.05,
	0.5:   0.005,
}

var TelegramHandlerDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "memes_bot",
	Subsystem:  "telegram_bot",
	Name:       "handler",
	Help:       "The latency of handler action",
	Objectives: summaryObjectives,
	MaxAge:     30 * time.Second,
	AgeBuckets: 3,
}, []string{"handler", "result"})
