package bloomservice

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"

	"github.com/bits-and-blooms/bloom/v3"
)

type Service interface {
	PostTraceIdToIndex(ctx context.Context, t TraceID, i Index) (string, error)
	GetIndex(ctx context.Context, t TraceID) (*[]Index, error)
	DeleteIndex(ctx context.Context, i Index) (string, error)
}

// TraceID is the shared trace ID of all spans in the trace.
type TraceID string

// Index is the ES index.
type Index string

var (
	ErrInconsistentTraceId = errors.New("Inconsistent traceId")
	ErrInconsistentIndex   = errors.New("Inconsistent index")
	ErrAlreadyExists       = errors.New("Already exists")
	ErrNotFound            = errors.New("Not found")
)

type inmemService struct {
	logger log.Logger
	m      map[Index]bloom.BloomFilter
}

func NewInmemService(logger log.Logger) Service {
	return &inmemService{
		logger: logger,
		m:      map[Index]bloom.BloomFilter{},
	}
}

func (s *inmemService) PostTraceIdToIndex(ctx context.Context, t TraceID, i Index) (string, error) {
	logger := log.With(s.logger, "method", "PostTraceIdToIndex")

	logger.Log("traceId", t)
	logger.Log("index", i)

	if len(t) == 0 {
		return "", ErrInconsistentTraceId
	}

	if len(t) == 0 {
		return "", ErrInconsistentIndex
	}

	_, ok := s.m[i]
	if !ok {
		s.m[i] = *bloom.NewWithEstimates(1_000_000, 0.0001)
	}
	filter := s.m[i]

	if filter.Test([]byte(t)) {
		return "Success", nil
	}
	filter.Add([]byte(t))

	return "Success", nil
}

func (s *inmemService) GetIndex(ctx context.Context, t TraceID) (*[]Index, error) {
	var indices []Index
	for idx, filter := range s.m {
		if filter.Test([]byte(t)) {
			indices = append(indices, idx)
		}
	}
	return &indices, nil
}

func (s *inmemService) DeleteIndex(ctx context.Context, i Index) (string, error) {
	_, ok := s.m[i]
	if !ok {
		return "", ErrNotFound
	}
	delete(s.m, i)
	return "Success", nil
}
