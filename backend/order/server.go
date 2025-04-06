// backend/order/server.go
package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "backend/pb/orderpb"
    "google.golang.org/grpc"
)

type orderServer struct {
    orderpb.UnimplementedOrderServiceServer
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
    log.Println("Creating order for:", req.Item)
    return &orderpb.OrderResponse{
        OrderId: "ORD123",
        Status:  "Order Created",
    }, nil
}

func (s *orderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
    log.Println("Canceling order:", req.OrderId)
    return &orderpb.CancelOrderResponse{
        Status: "Order Cancelled",
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    orderpb.RegisterOrderServiceServer(grpcServer, &orderServer{})
    fmt.Println("Order service running on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
