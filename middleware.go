package bloomservice

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostTraceIdToIndex(ctx context.Context, t TraceID, i Index) (msg string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostTraceIdToIndex", "traceId", t, "index", i, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostTraceIdToIndex(ctx, t, i)
}

func (mw loggingMiddleware) GetIndex(ctx context.Context, t TraceID) (l *[]Index, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetIndex", "tarceId", t, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetIndex(ctx, t)
}

func (mw loggingMiddleware) DeleteIndex(ctx context.Context, i Index) (msg string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteIndex", "index", i, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteIndex(ctx, i)
}
