package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "backend/pb/paymentpb"
    "google.golang.org/grpc"
)

type paymentServer struct {
    paymentpb.UnimplementedPaymentServiceServer
}

func (s *paymentServer) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.PaymentResponse, error) {
    log.Println("Processing payment for Order:", req.OrderId)
    return &paymentpb.PaymentResponse{
        PaymentId: "PAY123",
        Status:    "Payment Success",
    }, nil
}

func (s *paymentServer) RefundPayment(ctx context.Context, req *paymentpb.RefundRequest) (*paymentpb.RefundResponse, error) {
    log.Println("Refunding payment for Order:", req.OrderId)
    return &paymentpb.RefundResponse{
        Status: "Refund Success",
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    paymentpb.RegisterPaymentServiceServer(grpcServer, &paymentServer{})
    fmt.Println("Payment service running on port 50052...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
