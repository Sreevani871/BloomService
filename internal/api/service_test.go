package bloomservice

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

var (
	mockTraceId     = TraceID("E6S3zHkB1gqXhRfi8oGI")
	mockIndex       = Index("rum-freshdesk-2021.06.01-000059")
	mockDeleteIndex = Index("rum-freshdesk-2021.08.02-000077")
	mockGetResponse = []Index{Index("rum-freshdesk-2021.06.01-000059")}
)

func initiaizeTestService() Service {
	var logger log.Logger
	qs := NewInmemService(logger)
	return qs
}

var logger log.Logger
var bloomsvc = NewInmemService(logger)

func TestPostTraceIdToIndex(t *testing.T) {
	// svc := initiaizeTestService()

	type contextKey string
	ctx := context.Background()
	res, err := bloomsvc.PostTraceIdToIndex(context.WithValue(ctx, contextKey("foo"), "bar"), mockTraceId, mockIndex)
	assert.NoError(t, err)
	assert.Equal(t, "Success", res)
}

func TestPostTraceIdToIndexInconsistentTraceId(t *testing.T) {
	// svc := initiaizeTestService()

	type contextKey string
	ctx := context.Background()
	_, err := bloomsvc.PostTraceIdToIndex(context.WithValue(ctx, contextKey("foo"), "bar"), "", mockIndex)
	assert.Equal(t, err, ErrInconsistentTraceId)
}

func TestPostTraceIdToIndexInconsistentIndex(t *testing.T) {
	// svc := initiaizeTestService()

	type contextKey string
	ctx := context.Background()
	_, err := bloomsvc.PostTraceIdToIndex(context.WithValue(ctx, contextKey("foo"), "bar"), mockTraceId, "")
	assert.Equal(t, err, ErrInconsistentIndex)
}

func TestGetIndex(t *testing.T) {
	// svc := initiaizeTestService()

	type contextKey string
	ctx := context.Background()
	res, err := bloomsvc.GetIndex(context.WithValue(ctx, contextKey("foo"), "bar"), mockTraceId)
	assert.NoError(t, err)
	assert.Equal(t, res, mockGetResponse)
}

func TestDeleteIndex(t *testing.T) {
	// svc := initiaizeTestService()

	type contextKey string
	ctx := context.Background()
	res, err := bloomsvc.DeleteIndex(context.WithValue(ctx, contextKey("foo"), "bar"), mockIndex)
	assert.NoError(t, err)
	assert.Equal(t, "Success", res)
}

func TestDeleteIndexNotFound(t *testing.T) {
	// svc := initiaizeTestService()

	type contextKey string
	ctx := context.Background()
	_, err := bloomsvc.DeleteIndex(context.WithValue(ctx, contextKey("foo"), "bar"), mockDeleteIndex)
	assert.Equal(t, err, ErrNotFound)
}
