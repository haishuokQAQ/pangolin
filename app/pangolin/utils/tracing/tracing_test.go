package tracing

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestTracer(t *testing.T) {
	InitTracer("test")
	ctx := context.Background()

	sp, ctx := StartSpanFromContext(ctx, "op1")
	traceID, spanID, ok := GetTraceIDAndSpanID(sp)
	if !ok || traceID != spanID {
		t.Errorf("create root span error, ok: %v, trace id: %s, span id %s", ok, traceID, spanID)
	}

	sp1, _ := StartSpanFromContext(ctx, "op2")
	traceID1, spanID1, ok := GetTraceIDAndSpanID(sp1)
	if !ok || traceID1 != traceID {
		t.Errorf("create second span error, ok: %v, trace id: %s, span id %s", ok, traceID1, spanID1)
	}

	sp2, _ := StartSpanFromContext(ctx, "op2")
	traceID2, spanID2, ok := GetTraceIDAndSpanID(sp2)
	if !ok || traceID2 != traceID {
		t.Errorf("create third span error, ok: %v, trace id: %s, span id %s", ok, traceID2, spanID2)
	}

	sp.Finish()
	sp1.Finish()
	sp2.Finish()
}

func TestInjectAndExtractHeader(t *testing.T) {
	InitTracer("test")
	ctx := context.Background()
	sp, ctx := StartSpanFromContext(ctx, "op1")
	traceID, spanID, _ := GetTraceIDAndSpanID(sp)

	header := http.Header{}
	InjectHeader(ctx, header)
	traceStr := header.Get("Uber-Trace-Id")
	if traceStr == "" {
		t.Errorf("trace is not injected into header")
	}
	splits := strings.Split(traceStr, ":")
	if traceID != splits[0] || spanID != splits[1] {
		t.Errorf("trace id and span id are incorrect")
	}
	ctx = context.Background()
	ctx, err := ExtractHeader(ctx, header, "op2")
	if err != nil {
		t.Errorf("error to extract header")
	}
	newSp, ok := SpanFromContext(ctx)
	if !ok {
		t.Errorf("no span in context after extracting header")
	}
	newTraceId, newSpanId, _ := GetTraceIDAndSpanID(newSp)
	if newTraceId != traceID || newSpanId == spanID {
		t.Errorf("trace id and span id are incorrect")
	}
}

func TestExtractHeaderWithError(t *testing.T) {
	header := http.Header{}
	ctx := context.Background()
	ctx, err := ExtractHeader(ctx, header, "op2")
	if err == nil {
		t.Errorf("expect error when extracting header")
	}
	sp, ok := SpanFromContext(ctx)
	if ok || sp != nil {
		t.Errorf("expect span not found when extracting header")
	}
}

func TestInjectAndExtractMessage(t *testing.T) {
	InitTracer("test")
	ctx := context.Background()
	sp, ctx := StartSpanFromContext(ctx, "op1")
	traceID, spanID, _ := GetTraceIDAndSpanID(sp)
	msgStr := "msg"
	msg := []byte(msgStr)

	InjectMessage(ctx, &msg)
	if string(msg) == msgStr || !strings.HasSuffix(string(msg), msgStr) {
		t.Errorf("error to inject message")
	}

	ctx = context.Background()
	ctx, err := ExtractMessage(ctx, &msg, "subject")
	if err != nil {
		t.Errorf("error to extract message")
	}
	if string(msg) != msgStr {
		t.Errorf("error to extract message")
	}
	newSp, ok := SpanFromContext(ctx)
	if !ok {
		t.Errorf("no span in context after extracting header")
	}
	newTraceId, newSpanId, _ := GetTraceIDAndSpanID(newSp)
	if newTraceId != traceID || newSpanId == spanID {
		t.Errorf("trace id and span id are incorrect")
	}
}

func TestExtractMessageWithError(t *testing.T) {
	msgStr := "msg"
	msg := []byte(msgStr)
	ctx := context.Background()
	ctx, err := ExtractMessage(ctx, &msg, "subject")
	if err == nil {
		t.Errorf("expect error when extracting message")
	}
	sp, ok := SpanFromContext(ctx)
	if ok || sp != nil {
		t.Errorf("expect span not found when extracting message")
	}
}