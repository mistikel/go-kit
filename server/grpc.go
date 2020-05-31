package server

import (
	"context"
	"net"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/mistikel/go-kit/endpoints"
	"github.com/mistikel/go-kit/pb"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	searchMovie grpctransport.Handler
	logger      log.Logger
}

func NewGRPCServer(e endpoints.Endpoints, otTracer stdopentracing.Tracer, logger log.Logger) *GRPCServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &GRPCServer{
		searchMovie: grpctransport.NewServer(
			e.SearchMovie,
			decodeGRPCSearchMovieRequest,
			encodeGRPCSearchMovieResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "SearchMovie", logger)))...,
		),
		logger: logger,
	}
}

func (g *GRPCServer) Serve() {
	grpcListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		g.logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	g.logger.Log("transport", "gRPC", "addr", ":8081")
	g.logger.Log("transport", "gRPC", "addr", ":8081")
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterMovieServer(baseServer, g)
	baseServer.Serve(grpcListener)
}

func (g *GRPCServer) SearchMovie(ctx context.Context, req *pb.SearchRequest) (*pb.SearchReplys, error) {
	_, resp, err := g.searchMovie.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SearchReplys), nil
}

func decodeGRPCSearchMovieRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SearchRequest)
	return endpoints.SearchMovieRequest{Keyword: req.Keyword, Page: req.Page}, nil
}

func encodeGRPCSearchMovieResponse(_ context.Context, response interface{}) (interface{}, error) {
	resps := response.(endpoints.SearchMovieResponse)
	searchReplys := []*pb.SearchReply{}
	for _, resp := range resps.Movies {
		s := &pb.SearchReply{
			Id:     resp.ImdbID,
			Title:  resp.Title,
			Type:   resp.Type,
			Year:   resp.Year,
			Poster: resp.Poster,
		}
		searchReplys = append(searchReplys, s)
	}
	return &pb.SearchReplys{Search: searchReplys}, nil
}
