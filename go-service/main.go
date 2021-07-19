package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)

func task(ctx context.Context, n int) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "task")
	defer span.Finish()
	duration := time.Duration(n+1/2)
	span.LogKV("event", fmt.Sprintf("sleep:%s", duration.String()))
	time.Sleep(time.Second * duration)
	return n * n * n
}

func getMatrix(ctx context.Context, filler, rows, columns int) matrix {
	matrix := make([]int, rows)
	for i := 0; i < len(r); i++ {
		c := make([]int, columns)
	}
	return
}

func tracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// derive a new context from the request
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "middleware")
		//defer span.Finish()
		r = r.WithContext(ctx)

		//ctx, _ := tracer.Extract(
		//	opentracing.HTTPHeaders,
		//	opentracing.HTTPHeadersCarrier(r.Header))
		//childSpan := tracer.StartSpan("middleware", ext.RPCServerOption(ctx))
		//defer childSpan.Finish()

		span.Finish() // end now?
		next.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "index")
	//opentracing.ContextWithSpan()
	//defer span.Finish()
	q := r.URL.Query()
	num := q.Get("num")
	if num != "" {
		span.LogKV("event", fmt.Sprintf("num:%s", num))
		numAsInt, err := strconv.Atoi(num)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request. Query param must be numeric: ?num=<number>"))
		}
		n := task(opentracing.ContextWithSpan(ctx, span), numAsInt)
		w.Write([]byte(fmt.Sprintf("Answer: %d", n)))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request. Query param is required: ?num=<number>"))
	}
}

func init() {
	name, ok := os.LookupEnv("JAEGER_SERVICE_NAME")
	if !ok {
		log.Fatal("env variable JAEGER_SERVICE_NAME is required")
	}
	cfg := config.Configuration{
		ServiceName: name,
		Sampler: &config.SamplerConfig{
			Type: "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}
	//cfg, err := config.FromEnv()
	//if err != nil {
	//	log.Fatalf("could not read jaeger configs from env: %v", err)
	//}
	//cfg.Sampler.Param = 1
	var err error
	tracer, closer, err = cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Metrics(metrics.NullFactory))
	if err != nil {
		log.Fatalf("could not init tracer: %v", err)
	}
	opentracing.InitGlobalTracer(tracer)
}

func main()  {
	defer closer.Close()
	http.Handle("/", tracingMiddleware(http.HandlerFunc(index)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}