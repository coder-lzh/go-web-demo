package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go-web-demo/internal/logger"
)

func Logging() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		startTime := time.Now()

		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		tracer := otel.Tracer("go-web-demo")
		spanName := fmt.Sprintf("%s %s", ctx.Method(), string(ctx.Path()))
		otelCtx, span := tracer.Start(c, spanName)
		defer span.End()

		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		var requestBody string
		bodyBytes := ctx.Request.Body()
		if len(bodyBytes) > 0 {
			requestBody = string(bodyBytes)
		}

		queryParams := string(ctx.Request.URI().QueryString())

		var errStack string
		defer func() {
			if r := recover(); r != nil {
				errStack = string(debug.Stack())
				ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    500,
					"message": fmt.Sprintf("panic: %v", r),
				})
			}
		}()

		ctx.Next(otelCtx)

		duration := time.Since(startTime)

		fields := []zap.Field{
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("method", string(ctx.Method())),
			zap.String("path", string(ctx.Path())),
			zap.Int("status_code", ctx.Response.StatusCode()),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("client_ip", getClientIP(ctx)),
			zap.String("query_params", queryParams),
			zap.String("request_body", requestBody),
			zap.String("response", string(ctx.Response.Body())),
		}

		message := fmt.Sprintf("%s %s", ctx.Method(), ctx.Path())
		if errStack != "" {
			fields = append(fields, zap.String("error_stack", errStack))
			logger.Error(message, convertFields(fields)...)
		} else if ctx.Response.StatusCode() >= 500 {
			logger.Error(message, convertFields(fields)...)
		} else if ctx.Response.StatusCode() >= 400 {
			logger.Warn(message, convertFields(fields)...)
		} else {
			logger.Info(message, convertFields(fields)...)
		}
	}
}

func getClientIP(ctx *app.RequestContext) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "Proxy-Client-IP", "WL-Proxy-Client-IP"} {
		if ip := string(ctx.GetHeader(header)); ip != "" {
			if idx := indexOf(ip, ','); idx != -1 {
				return ip[:idx]
			}
			return ip
		}
	}

	ip := ctx.ClientIP()
	if ip == "" {
		ip = "unknown"
	}

	return ip
}

func convertFields(fields []zap.Field) []interface{} {
	result := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		result = append(result, f)
	}
	return result
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}
