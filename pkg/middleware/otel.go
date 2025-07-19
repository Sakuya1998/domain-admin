package middleware

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"domain-admin/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)



// InitTracer 初始化OpenTelemetry追踪器
// serviceName: 服务名称
// otlpEndpoint: OTLP端点地址
// 返回清理函数
func InitTracer(serviceName, otlpEndpoint string) func() {
	ctx := context.Background()

	// 配置选项
	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(otlpEndpoint),
	}

	// 根据环境变量决定是否使用安全连接
	if os.Getenv("OTEL_INSECURE") == "true" {
		options = append(options, otlptracehttp.WithInsecure())
	}

	exp, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		logger.Errorf("Failed to create OTLP exporter", "error", err)
		// 返回空操作函数，避免程序崩溃
		return func() {}
	}

	// 获取采样率配置
	sampleRate := getSampleRate()

	// 获取环境信息
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// 获取服务版本
	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "0.1.0"
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(env),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceInstanceIDKey.String(getHostname()),
		)),
	)

	otel.SetTracerProvider(tp)

	logger.Infof("OpenTelemetry tracer initialized",
		"service", serviceName,
		"endpoint", otlpEndpoint,
		"environment", env,
		"sampleRate", sampleRate,
	)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			logger.Errorf("Error shutting down tracer provider", "error", err)
		} else {
			logger.Info("OpenTelemetry tracer shutdown successfully")
		}
	}
}

// getSampleRate 从环境变量获取采样率
func getSampleRate() float64 {
	rateStr := os.Getenv("OTEL_SAMPLE_RATE")
	if rateStr == "" {
		return 1.0 // 默认100%采样
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil || rate < 0 || rate > 1 {
		logger.Warnf("Invalid sample rate, using default", "rate", rateStr)
		return 1.0
	}

	return rate
}

// getHostname 获取主机名用于实例ID
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Warnf("Failed to get hostname", "error", err)
		return "unknown"
	}
	return hostname
}

func OTLPMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 过滤掉健康检查等不需要追踪的请求
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/health") ||
			strings.HasPrefix(path, "/metrics") ||
			strings.HasPrefix(path, "/swagger") {
			c.Next()
			return
		}
		
		// 对其他请求应用OpenTelemetry中间件
		otelgin.Middleware("domain-admin-gin")(c)
	})
}
