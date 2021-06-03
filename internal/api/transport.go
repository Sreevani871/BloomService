package bloomservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	apiRouter.Methods(http.MethodPost).PathPrefix("/sherlock/routes").Handler(httptransport.NewServer(
		e.PostEndpoint,
		decodePostRequest,
		encodeResponse,
		options...,
	))
	apiRouter.Methods(http.MethodGet).PathPrefix("/sherlock/routes").Handler(httptransport.NewServer(
		e.GetEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))
	apiRouter.Methods(http.MethodDelete).PathPrefix("/sherlock/routes/{id}").Handler(httptransport.NewServer(
		e.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getRequest
	traceId := r.URL.Query().Get("traceId")
	req = getRequest{
		Id: traceId,
	}
	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteRequest{IDX: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentTraceId:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

type errorer interface {
	error() error
}
