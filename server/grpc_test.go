package server

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/mistikel/go-kit/endpoints"
	"github.com/mistikel/go-kit/pb"
	"github.com/mistikel/go-kit/pkg/imdb"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func setupGRPC(t *testing.T) (*grpc.Server, *bufconn.Listener) {
	tracer := stdopentracing.GlobalTracer()
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	service := imdb.New(logger)
	endpoint := endpoints.NewEndpoints(service, tracer, logger)
	server := NewGRPCServer(endpoint, tracer, logger)

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)

	s := grpc.NewServer()
	pb.RegisterMovieServer(s, server)
	go func() {
		if err := s.Serve(listener); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()
	return s, listener
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}
func TestGRPC(t *testing.T) {
	ctx := context.Background()
	server, listener := setupGRPC(t)
	defer server.Stop()

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewMovieClient(conn)

	params := &pb.SearchRequest{
		Keyword: "Batman",
		Page:    1,
	}
	resp, err := client.SearchMovie(ctx, params)
	if err != nil {
		t.Logf("Error: %v", err)
	}

	if len(resp.Search) == 0 {
		t.Logf("Error data not found")
	}
}
