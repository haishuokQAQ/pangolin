package logger

import (
	"context"

	tspctx "gitlab.p1staff.com/tsp/common/context"
	"gitlab.p1staff.com/tsp/common/tracing"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithTraceInCtx(ctx context.Context) Logger

	// Flushing any buffered log entries. Applications should take care to call Flush before exiting.
	Flush()

	// Release resources
	Close()
}

const (
	traceIDKey  = "trace_id"
	spanIDKey   = "span_id"
	userIDKey   = "user_id"
	serviceKey  = "service"
	hostnameKey = "hostname"
	ipKey       = "ip"
)

func CtxToMap(ctx context.Context) map[string]interface{} {
	data := make(map[string]interface{}, 0)
	if ctx != nil {
		if span, ok := tspctx.GetTrace(ctx); ok {
			if traceID, spanID, ok := tracing.GetTraceIDAndSpanID(span); ok {
				data[traceIDKey] = traceID
				data[spanIDKey] = spanID
			}
		}
		if user, ok := tspctx.GetUser(ctx); ok {
			data[userIDKey] = user.GetID()
		}
		if service, ok := tspctx.GetService(ctx); ok {
			data[serviceKey] = service
		}
		if hostname, ok := tspctx.GetHostName(ctx); ok {
			data[hostnameKey] = hostname
		}
		if ip, ok := tspctx.GetIP(ctx); ok {
			data[ipKey] = ip
		}
	}
	return data
}
