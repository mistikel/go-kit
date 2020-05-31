package imdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	// todo > move to env
	IMDB_API_URL = "http://www.omdbapi.com/"
	IMDB_API_KEY = "faf7e5bb&s"
)

var (
	encode = func(context.Context, *http.Request, interface{}) error { return nil }
	decode = func(_ context.Context, r *http.Response) (interface{}, error) {
		var sr SearchResult
		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			return SearchResult{}, err
		}
		return sr, nil
	}
)

type Service interface {
	SearchMovie(ctx context.Context, keyword string, page int64) ([]Movie, error)
}

func New(logger log.Logger) Service {
	var svc Service
	{
		svc = &imService{}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

type imService struct{}

func (s *imService) SearchMovie(ctx context.Context, keyword string, page int64) ([]Movie, error) {
	rawurl := fmt.Sprintf("%s?apikey=%s?&s=%s&page=%d", IMDB_API_URL, IMDB_API_KEY, keyword, page)
	client := httptransport.NewClient("GET", mustParse(rawurl), encode, decode)
	resp, err := client.Endpoint()(ctx, struct{}{})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(SearchResult)
	if !ok {
		return nil, errors.New("Can not parse result")
	}
	return response.Result, nil
}

type SearchResult struct {
	Result []Movie `json:"Search"`
}

type Movie struct {
	ImdbID string `json:"imdbID"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Year   string `json:"year"`
	Poster string `json:"poster"`
}

func mustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
