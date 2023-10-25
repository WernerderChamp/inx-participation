package prometheus

import (
	"github.com/iotaledger/hive.go/app"
)

// ParametersPrometheus contains the definition of the parameters used by Prometheus.
type ParametersPrometheus struct {
	// BindAddress defines the bind address on which the Prometheus exporter listens on.
	BindAddress string `default:"inx-participation:9312" usage:"the bind address on which the Prometheus HTTP server listens on"`
	// Enabled defines whether the prometheus plugin is enabled.
	Enabled bool `default:"true" usage:"whether the prometheus plugin is enabled"`
	// ParticipationMetrics defines whether to include the participation metrics.
	ParticipationMetrics bool `default:"true" usage:"whether to include the participation metrics"`
	// GoMetrics defines whether to include go metrics.
	GoMetrics bool `default:"false" usage:"whether to include go metrics"`
	// ProcessMetrics defines whether to include process metrics.
	ProcessMetrics bool `default:"false" usage:"whether to include process metrics"`
	// PromhttpMetrics defines whether to include promhttp metrics.
	PromhttpMetrics bool `default:"false" usage:"whether to include promhttp metrics"`
}

var ParamsPrometheus = &ParametersPrometheus{}

var params = &app.ComponentParams{
	Params: map[string]any{
		"prometheus": ParamsPrometheus,
	},
	Masked: nil,
}
