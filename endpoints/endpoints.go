package endpoints

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/mistikel/go-kit/pkg/imdb"
)

type Endpoints struct {
	SearchMovie endpoint.Endpoint
}

func NewEndpoints(s imdb.Service, otTracer stdopentracing.Tracer, logger log.Logger) Endpoints {
	var searchMovie endpoint.Endpoint
	{
		searchMovie = NewSearchMovieEndpoint(s)
		searchMovie = opentracing.TraceServer(otTracer, "SearchMovie")(searchMovie)
	}
	return Endpoints{
		SearchMovie: NewSearchMovieEndpoint(s),
	}
}

func NewSearchMovieEndpoint(s imdb.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SearchMovieRequest)
		m, e := s.SearchMovie(ctx, req.Keyword, req.Page)
		return SearchMovieResponse{Movies: m, Err: e}, nil
	}
}

type SearchMovieRequest struct {
	Keyword string `json:"keyword"`
	Page    int64  `json:"page"`
}

type SearchMovieResponse struct {
	Movies []imdb.Movie `json:"movies"`
	Err    error        `json:"error"`
}
