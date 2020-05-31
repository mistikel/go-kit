package imdb

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) SearchMovie(ctx context.Context, keyword string, page int64) (m []Movie, e error) {
	defer func() {
		mw.logger.Log("method", "SearchMovie", "keyword", keyword, "page", page, "err", e)
	}()
	return mw.next.SearchMovie(ctx, keyword, page)
}
