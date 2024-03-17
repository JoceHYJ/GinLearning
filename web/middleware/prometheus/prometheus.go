package prometheus

import (
	"GinLearning/web"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func (m MiddlewareBuilder) Build() web.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Namespace: m.Namespace,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(vector)

	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()
				status := ctx.RespStatusCode
				// 响应时间
				pattern := ctx.MatchedRoute
				// 命中路由 404
				if pattern == "" {
					pattern = "unknown"
				}
				vector.WithLabelValues(pattern, ctx.Req.Method,
					strconv.Itoa(status)).Observe(float64(duration))
			}()
			next(ctx)
		}
	}
}
