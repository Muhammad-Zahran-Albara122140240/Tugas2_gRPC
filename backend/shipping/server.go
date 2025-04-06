package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "Tugas2_PWL/backend/pb/shippingpb"
    "google.golang.org/grpc"
)

type shippingServer struct {
    shippingpb.UnimplementedShippingServiceServer
}

func (s *shippingServer) Ship(ctx context.Context, req *shippingpb.ShipRequest) (*shippingpb.ShipResponse, error) {
    log.Println("Shipping order to:", req.Address)
    return &shippingpb.ShipResponse{
        ShippingId: "SHIP123",
        Status:     "Order Shipped",
    }, nil
}

func (s *shippingServer) CancelShipping(ctx context.Context, req *shippingpb.CancelShipRequest) (*shippingpb.CancelShipResponse, error) {
    log.Println("Canceling shipping for Shipping ID:", req.ShippingId)
    return &shippingpb.CancelShipResponse{
        Status: "Shipping Cancelled",
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    shippingpb.RegisterShippingServiceServer(grpcServer, &shippingServer{})
    fmt.Println("Shipping service running on port 50053...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
