package main

import (
	"context"
	"fmt"
	"log"

	"Tugas2_PWL/backend/pb/orderpb"
	"Tugas2_PWL/backend/pb/paymentpb"
	"Tugas2_PWL/backend/pb/shippingpb"

	"google.golang.org/grpc"
)

func main() {
	// Connect ke Order Service
	orderConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	// Connect ke Payment Service
	paymentConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Payment Service: %v", err)
	}
	defer paymentConn.Close()
	paymentClient := paymentpb.NewPaymentServiceClient(paymentConn)

	// Connect ke Shipping Service
	shippingConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Shipping Service: %v", err)
	}
	defer shippingConn.Close()
	shippingClient := shippingpb.NewShippingServiceClient(shippingConn)

	ctx := context.Background()

	// === STEP 1: Buat order ===
	orderRes, err := orderClient.CreateOrder(ctx, &orderpb.OrderRequest{
		Item: "Keyboard Gaming",
	})
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	fmt.Println("Order created:", orderRes.OrderId, "-", orderRes.Status)

	// === STEP 2: Proses pembayaran ===
	paymentRes, err := paymentClient.ProcessPayment(ctx, &paymentpb.PaymentRequest{
		OrderId: orderRes.OrderId,
		Amount:  300000,
	})
	if err != nil {
		log.Fatalf("Failed to process payment: %v", err)

		// Kompensasi: Batalkan order
		_, _ = orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{
			OrderId: orderRes.OrderId,
		})
		return
	}
	fmt.Println("Payment processed:", paymentRes.PaymentId, "-", paymentRes.Status)

	// === STEP 3: Kirim barang ===
    shippingRes, err := shippingClient.Ship(ctx, &shippingpb.ShipRequest{
        Address: "Jl. Mawar No. 123",
    })
    
	if err != nil {
		log.Fatalf("Failed to ship order: %v", err)

		// Kompensasi: Refund pembayaran
		_, _ = paymentClient.RefundPayment(ctx, &paymentpb.RefundRequest{
			PaymentId: paymentRes.PaymentId,
		})
		// Kompensasi: Batalkan order
		_, _ = orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{
			OrderId: orderRes.OrderId,
		})
		return
	}
	fmt.Println("Order shipped:", shippingRes.ShippingId, "-", shippingRes.Status)
}
