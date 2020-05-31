package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/mistikel/go-kit/endpoints"
	"github.com/mistikel/go-kit/pkg/imdb"
	stdopentracing "github.com/opentracing/opentracing-go"
)

func TestHTTP(t *testing.T) {
	tracer := stdopentracing.GlobalTracer()
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	service := imdb.New(logger)
	endpoint := endpoints.NewEndpoints(service, tracer, logger)

	srv := NewHTTPServer(endpoint, tracer, logger)
	server := httptest.NewServer(srv.handler)
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL+"/search", strings.NewReader(`{"keyword": "Batman","page": 1}`))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("error request", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var movies endpoints.SearchMovieResponse
	err = json.Unmarshal(body, &movies)
	if err != nil {
		t.Error("error marshal response")
	}
}
