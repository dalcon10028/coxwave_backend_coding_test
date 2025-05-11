package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	couponv1 "github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1"
	"github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1/couponv1connect"
)

func main() {
	client := couponv1connect.NewCouponServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	res, err := client.Hello(
		context.Background(),
		connect.NewRequest(&couponv1.HelloRequest{Name: "테스트"}),
	)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("response: %s", res.Msg.Message)
}
