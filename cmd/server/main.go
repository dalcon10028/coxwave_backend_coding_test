package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"

	couponv1 "github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1"
	"github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1/couponv1connect"
	"github.com/dalcon10028/coxwave_backend_coding_test/internal/database"
	"github.com/dalcon10028/coxwave_backend_coding_test/internal/repository/sqlite"
)

type CouponServer struct {
	campaignRepo *sqlite.CampaignRepository
}

func NewCouponServer(db *sql.DB) *CouponServer {
	return &CouponServer{
		campaignRepo: sqlite.NewCampaignRepository(db),
	}
}

func (s *CouponServer) CreateCampaign(
	ctx context.Context,
	req *connect.Request[couponv1.CreateCampaignRequest],
) (*connect.Response[couponv1.CreateCampaignResponse], error) {
	campaign, err := s.campaignRepo.Create(ctx, req.Msg.TotalCoupons, req.Msg.StartAt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&couponv1.CreateCampaignResponse{
		Campaign: &couponv1.Campaign{
			Id:               campaign.ID,
			TotalCoupons:     campaign.TotalCoupons,
			RemainingCoupons: campaign.RemainingCoupons,
			StartAt:          timestamppb.New(campaign.StartAt),
			CreatedAt:        timestamppb.New(campaign.CreatedAt),
		},
	})
	return res, nil
}

func (s *CouponServer) GetCampaign(
	ctx context.Context,
	req *connect.Request[couponv1.GetCampaignRequest],
) (*connect.Response[couponv1.GetCampaignResponse], error) {
	campaign, err := s.campaignRepo.Get(ctx, req.Msg.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&couponv1.GetCampaignResponse{
		Campaign: &couponv1.Campaign{
			Id:               campaign.ID,
			TotalCoupons:     campaign.TotalCoupons,
			RemainingCoupons: campaign.RemainingCoupons,
			StartAt:          timestamppb.New(campaign.StartAt),
			CreatedAt:        timestamppb.New(campaign.CreatedAt),
		},
	})
	return res, nil
}

func (s *CouponServer) IssueCoupon(
	ctx context.Context,
	req *connect.Request[couponv1.IssueCouponRequest],
) (*connect.Response[couponv1.IssueCouponResponse], error) {
	err := s.campaignRepo.IssueCoupon(ctx, req.Msg.CampaignId, req.Msg.Code)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&couponv1.IssueCouponResponse{
		Coupon: &couponv1.Coupon{
			CampaignId: req.Msg.CampaignId,
			Code:       req.Msg.Code,
			CreatedAt:  timestamppb.New(time.Now()),
		},
	})
	return res, nil
}

func (s *CouponServer) GetIssuedCoupons(
	ctx context.Context,
	req *connect.Request[couponv1.GetIssuedCouponsRequest],
) (*connect.Response[couponv1.GetIssuedCouponsResponse], error) {
	codes, err := s.campaignRepo.GetIssuedCodes(ctx, req.Msg.CampaignId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&couponv1.GetIssuedCouponsResponse{
		Codes: codes,
	})
	return res, nil
}

func main() {
	db, err := database.Connect(database.NewConfig("./data/coupon.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	couponServer := NewCouponServer(db)
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
