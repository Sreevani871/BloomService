package bloomservice

import (
	"context"
	"log"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostEndpoint   endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostEndpoint:   MakePostEndpoint(s),
		GetEndpoint:    MakeGetEndpoint(s),
		DeleteEndpoint: MakeDeleteEndpoint(s),
	}
}

// MakePostProfileEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		log.Print(request)
		msg, e := s.PostTraceIdToIndex(ctx, TraceID(req.TraceId), Index(req.Index))
		return postResponse{Msg: msg}, e
	}
}

// MakeGetProfileEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		p, e := s.GetIndex(ctx, TraceID(req.Id))
		return getResponse{Indices: &p, Err: e}, nil
	}
}

// MakeDeleteProfileEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteRequest)
		msg, e := s.DeleteIndex(ctx, Index(req.IDX))
		return deleteResponse{Msg: msg}, e
	}
}

type postRequest struct {
	TraceId string `json:"traceId"`
	Index   string `json:"index"`
}

type postResponse struct {
	Err error  `json:"err,omitempty"`
	Msg string `json:"msg,omitempty"`
}

func (r postResponse) error() error { return r.Err }

type getRequest struct {
	Id string `json:"traceId"`
}

type getResponse struct {
	Indices *[]Index `json:"indices,omitempty"`
	Err     error    `json:"err,omitempty"`
}

type deleteRequest struct {
	IDX string `json:"index"`
}

type deleteResponse struct {
	Err error  `json:"err,omitempty"`
	Msg string `json:"msg,omitempty"`
}
