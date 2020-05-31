package imdb

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestSearchMovie(t *testing.T) {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	s := New(logger)

	res, err := s.SearchMovie(context.Background(), "Batman", 1)
	if err != nil {
		t.Error("error: ", err)
	}

	if len(res) < 0 {
		t.Error("error: ", err)
	}
}
