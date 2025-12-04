package utils

import "context"

type contextKey string

const (
	CtxCauserID  contextKey = "causer_id"
	CtxIP        contextKey = "ip"
	CtxUserAgent contextKey = "user_agent"
	CtxTraceID   contextKey = "trace_id"
)

func GetCauserID(ctx context.Context) string {
	if v := ctx.Value(CtxCauserID); v != nil {
		return v.(string)
	}
	return ""
}

func GetIP(ctx context.Context) string {
	if v := ctx.Value(CtxIP); v != nil {
		return v.(string)
	}
	return ""
}

func GetUserAgent(ctx context.Context) string {
	if v := ctx.Value(CtxUserAgent); v != nil {
		return v.(string)
	}
	return ""
}

func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(CtxTraceID); v != nil {
		return v.(string)
	}
	return ""
}
