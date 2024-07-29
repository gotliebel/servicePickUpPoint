package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"io"
	"log"
)

func MustSetup(_ context.Context, serviceName string) io.Closer {
	cfg := traceconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &traceconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &traceconfig.ReporterConfig{
			// LogSpans: true,
		},
	}

	tracer, c, err := cfg.NewTracer(traceconfig.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger %s", err)
	}

	opentracing.SetGlobalTracer(tracer)
	return c
}
