package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"Tugas2_PWL/backend/pb/orderpb"
	"Tugas2_PWL/backend/pb/paymentpb"
	"Tugas2_PWL/backend/pb/shippingpb"

	"google.golang.org/grpc"
)

func main() {
	// Ambil argumen dari CLI
	item := flag.String("item", "Default Item", "Nama item")
	amount := flag.Float64("amount", 0, "Jumlah pembayaran")
	address := flag.String("address", "Default Address", "Alamat pengiriman")
	simulateFailShipping := flag.Bool("fail-shipping", false, "Simulasikan kegagalan shipping")
	flag.Parse()

	fmt.Println("✅ Orchestrator Server is running...")

	// === gRPC Connect ===
	orderConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	paymentConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Payment Service: %v", err)
	}
	defer paymentConn.Close()
	paymentClient := paymentpb.NewPaymentServiceClient(paymentConn)

	shippingConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Shipping Service: %v", err)
	}
	defer shippingConn.Close()
	shippingClient := shippingpb.NewShippingServiceClient(shippingConn)

	ctx := context.Background()

	// === STEP 1: Buat order ===
	orderRes, err := orderClient.CreateOrder(ctx, &orderpb.OrderRequest{
		Item: *item,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create order: %v", err)
	}
	fmt.Printf("[✓] Order Created for: %s - Status: %s\n", *item, orderRes.Status)

	// === STEP 2: Proses pembayaran ===
	paymentRes, err := paymentClient.ProcessPayment(ctx, &paymentpb.PaymentRequest{
		OrderId: orderRes.OrderId,
		Amount:  *amount,
	})
	if err != nil {
		fmt.Printf("[x] Payment Failed: Rp%.0f - Reason: %v\n", *amount, err)

		// Kompensasi: Batalkan order
		_, _ = orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{
			OrderId: orderRes.OrderId,
		})
		fmt.Printf("[!] Order Cancelled: %s\n", orderRes.OrderId)
		return
	}
	fmt.Printf("[✓] Payment Success: Rp%.0f - Status: %s\n", *amount, paymentRes.Status)

	// === STEP 3: Kirim barang ===
	if *simulateFailShipping {
		fmt.Println("❌ Simulating shipping failure...")
		err = fmt.Errorf("simulated shipping failure")
	} else {
		var shippingRes *shippingpb.ShipResponse
		shippingRes, err = shippingClient.Ship(ctx, &shippingpb.ShipRequest{
			Address: *address,
		})
		if err == nil {
			fmt.Printf("[✓] Shipped to: %s - Status: %s\n", *address, shippingRes.Status)
		}
	}

	if err != nil {
		fmt.Printf("[x] Shipping Failed: %v\n", err)

		// Kompensasi: Refund pembayaran
		_, _ = paymentClient.RefundPayment(ctx, &paymentpb.RefundRequest{
			PaymentId: paymentRes.PaymentId,
		})
		fmt.Printf("[!] Refunded payment ID: %s - Status: Refund Success\n", paymentRes.PaymentId)

		// Kompensasi: Batalkan order
		_, _ = orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{
			OrderId: orderRes.OrderId,
		})
		fmt.Printf("[!] Order Cancelled: %s\n", orderRes.OrderId)
		return
	}
}
