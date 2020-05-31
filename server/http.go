package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mistikel/go-kit/endpoints"
)

type HTTPServer struct {
	handler http.Handler
	logger  log.Logger
}

func NewHTTPServer(e endpoints.Endpoints, otTracer stdopentracing.Tracer, logger log.Logger) *HTTPServer {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r := mux.NewRouter()
	r.Methods("POST").Path("/search").Handler(httptransport.NewServer(
		e.SearchMovie,
		decodeHTTPSearchMovieRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Search", logger)))...,
	))

	return &HTTPServer{
		handler: r,
		logger:  logger,
	}
}

func (h *HTTPServer) Serve() {
	h.logger.Log("server", "HTTP", "addr", ":8080")
	http.ListenAndServe(":8080", h.handler)
}

func decodeHTTPSearchMovieRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.SearchMovieRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
