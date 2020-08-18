package context

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
)

type contextKey struct{}

type contextMap map[string]interface{}

const (
	Trace       = "trace"
	Transaction = "transaction"
	User        = "user"
	Service     = "service"
	Hostname    = "hostname"
	IP          = "ip"

	ginCtxKey = "ctx"

	emptyStr = ""
)

var (
	tspContextKey = contextKey{}
)

func findOrCreateContextMap(ctx context.Context) (contextMap, context.Context) {
	ctxMap, ok := ctx.Value(tspContextKey).(contextMap)
	if !ok {
		ctxMap = make(contextMap)
		ctx = context.WithValue(ctx, tspContextKey, ctxMap)
	}
	return ctxMap, ctx
}

func findContextMap(ctx context.Context) (contextMap, bool) {
	ctxMap, ok := ctx.Value(tspContextKey).(contextMap)
	return ctxMap, ok
}

// WithTrace adds trace into context
func WithTrace(ctx context.Context, trace opentracing.Span) context.Context {
	ctxMap, ctx := findOrCreateContextMap(ctx)
	ctxMap[Trace] = trace
	return ctx
}

// WithoutTrace removes trace from context
func WithoutTrace(ctx context.Context) context.Context {
	return WithoutValue(ctx, Trace)
}

// GetTrace returns trace in context
func GetTrace(ctx context.Context) (opentracing.Span, bool) {
	if ctxMap, ok := findContextMap(ctx); ok {
		if trace, ok := ctxMap[Trace]; ok {
			ret, ok := trace.(opentracing.Span)
			return ret, ok
		}
	}
	return nil, false
}

// WithTransaction adds transaction into context
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	ctxMap, ctx := findOrCreateContextMap(ctx)
	ctxMap[Transaction] = tx
	return ctx
}

// WithoutTransaction removes transaction from context
func WithoutTransaction(ctx context.Context) context.Context {
	return WithoutValue(ctx, Transaction)
}

// GetTransaction returns transaction in context
func GetTransaction(ctx context.Context) (*gorm.DB, bool) {
	if ctxMap, ok := findContextMap(ctx); ok {
		if tx, ok := ctxMap[Transaction]; ok {
			ret, ok := tx.(*gorm.DB)
			return ret, ok
		}
	}
	return nil, false
}

// WithService adds service into context
func WithService(ctx context.Context, service string) context.Context {
	ctxMap, ctx := findOrCreateContextMap(ctx)
	ctxMap[Service] = service
	return ctx
}

// WithoutService removes service from context
func WithoutService(ctx context.Context) context.Context {
	return WithoutValue(ctx, Service)
}

// GetService returns service in context
func GetService(ctx context.Context) (string, bool) {
	if ctxMap, ok := findContextMap(ctx); ok {
		if service, ok := ctxMap[Service]; ok {
			ret, ok := service.(string)
			return ret, ok
		}
	}
	return emptyStr, false
}

// WithHostname adds hostname into context
func WithHostname(ctx context.Context, hostname string) context.Context {
	ctxMap, ctx := findOrCreateContextMap(ctx)
	ctxMap[Hostname] = hostname
	return ctx
}

// WithoutHostName removes hostname from context
func WithoutHostName(ctx context.Context) context.Context {
	return WithoutValue(ctx, Hostname)
}

// GetHostName returns host name in context
func GetHostName(ctx context.Context) (string, bool) {
	if ctxMap, ok := findContextMap(ctx); ok {
		if hostname, ok := ctxMap[Hostname]; ok {
			ret, ok := hostname.(string)
			return ret, ok
		}
	}
	return emptyStr, false
}

// WithIP adds ip into context
func WithIP(ctx context.Context, ip string) context.Context {
	ctxMap, ctx := findOrCreateContextMap(ctx)
	ctxMap[IP] = ip
	return ctx
}

// WithoutIP removes ip from context
func WithoutIP(ctx context.Context) context.Context {
	return WithoutValue(ctx, IP)
}

// GetIP returns ip in context
func GetIP(ctx context.Context) (string, bool) {
	if ctxMap, ok := findContextMap(ctx); ok {
		if ip, ok := ctxMap[IP]; ok {
			ret, ok := ip.(string)
			return ret, ok
		}
	}
	return emptyStr, false
}

// WithValue adds any key/value into context
func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	ctxMap, ctx := findOrCreateContextMap(ctx)
	ctxMap[key] = value
	return ctx
}

// WithoutValue removes key from context
func WithoutValue(ctx context.Context, key string) context.Context {
	if ctxMap, ok := findContextMap(ctx); ok {
		delete(ctxMap, key)
	}
	return ctx
}

// GetValue returns key/value in context
func GetValue(ctx context.Context, key string) (interface{}, bool) {
	if ctxMap, ok := findContextMap(ctx); ok {
		val, ok := ctxMap[key]
		return val, ok
	}
	return nil, false
}

// ShallowCopyCtx returns a copied context from an exist context, without transaction and trace
func ShallowCopyCtx(ctx context.Context) context.Context {
	newCtx := context.Background()
	if ctxMap, ok := findContextMap(ctx); ok {
		newCtxMap := make(contextMap)
		for key, val := range ctxMap {
			// ignore transaction
			if key != Transaction {
				newCtxMap[key] = val
			}
		}
		newCtx = context.WithValue(newCtx, tspContextKey, newCtxMap)
	}
	return newCtx
}

// InjectGinContext adds context into gin.Context
func InjectGinContext(gc *gin.Context, ctx context.Context) {
	gc.Set(ginCtxKey, ctx)
}

// ExtractFromGinContext gets context from gin.Context
func ExtractFromGinContext(gc *gin.Context) (context.Context, bool) {
	val := gc.Value(ginCtxKey)
	if val == nil {
		return nil, false
	}
	ctx, ok := val.(context.Context)
	return ctx, ok
}
