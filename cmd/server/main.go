package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	couponv1 "github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1"
	"github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1/couponv1connect"
)

type CouponServer struct{}

func (s *CouponServer) Hello(
	ctx context.Context,
	req *connect.Request[couponv1.HelloRequest],
) (*connect.Response[couponv1.HelloResponse], error) {
	res := connect.NewResponse(&couponv1.HelloResponse{
		Message: fmt.Sprintf("안녕하세요, %s님!", req.Msg.Name),
	})
	return res, nil
}

func main() {
	couponServer := &CouponServer{}
	mux := http.NewServeMux()
	path, handler := couponv1connect.NewCouponServiceHandler(couponServer)
	mux.Handle(path, handler)

	var url string = "localhost:8080"

	log.Println("Server is running on ", url)
	http.ListenAndServe(
		url,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
