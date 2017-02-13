package util

import (
	"google.golang.org/grpc"
	"net"
	"fmt"
	"golang.org/x/net/context"
	"github.com/vendasta/gosdks/logging"
)

//Creates a basic GRPC Server with the specified interceptors
func CreateGrpcServer(interceptors ...grpc.UnaryServerInterceptor) (*grpc.Server) {
    s := grpc.NewServer(
        grpc.UnaryInterceptor(chainUnaryServer(interceptors...)),
    )
    return s
} 

func StartGrpcServer(server *grpc.Server, port int) (error) {
    var lis net.Listener
    var err error
    
    if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
        logging.Errorf(context.Background(), "Error creating GRPC listening socket: %s", err.Error())
        return err
    }

    //The following call blocks until an error occurs
    return server.Serve(lis)
}

// AuthInterceptor handles auth through GRPC, does not work on local environments
// func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//     var err error
//     ctx, err = accessmanagement.AttachUserInfoToContext(ctx); if err != nil {
//         return nil, grpc.Errorf(codes.PermissionDenied, "Call isn't authenticated.")
//     }
// 	return handler(ctx, req)
// }

func NoAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    return handler(ctx, req)
}

// ChainUnaryServer combines multiple grpc.UnaryServerInterceptor into a single grpc.UnaryServerInterceptor (required by GRPC)
func chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		buildChain := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return current(currentCtx, currentReq, info, next)
			}
		}
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = buildChain(interceptors[i], chain)
		}
		return chain(ctx, req)
	}
}
