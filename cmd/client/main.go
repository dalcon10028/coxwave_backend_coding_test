package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	couponv1 "github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1"
	"github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1/couponv1connect"
)

func main() {
	client := couponv1connect.NewCouponServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	ctx := context.Background()

	// 1. 캠페인 생성
	createRes, err := client.CreateCampaign(
		ctx,
		connect.NewRequest(&couponv1.CreateCampaignRequest{
			TotalCoupons: 100,
			StartAt:      time.Now().Unix(),
		}),
	)
	if err != nil {
		log.Fatalf("캠페인 생성 실패: %v", err)
	}
	campaign := createRes.Msg.Campaign
	log.Printf("캠페인 생성됨: ID=%d, 총 쿠폰=%d, 남은 쿠폰=%d",
		campaign.Id, campaign.TotalCoupons, campaign.RemainingCoupons)

	// 2. 캠페인 조회
	getRes, err := client.GetCampaign(
		ctx,
		connect.NewRequest(&couponv1.GetCampaignRequest{
			Id: campaign.Id,
		}),
	)
	if err != nil {
		log.Fatalf("캠페인 조회 실패: %v", err)
	}
	log.Printf("캠페인 조회됨: ID=%d, 총 쿠폰=%d, 남은 쿠폰=%d",
		getRes.Msg.Campaign.Id, getRes.Msg.Campaign.TotalCoupons, getRes.Msg.Campaign.RemainingCoupons)

	// 3. 쿠폰 발급 - 101개
	for i := 0; i < 101; i++ {
		issueRes, err := client.IssueCoupon(
			ctx,
			connect.NewRequest(&couponv1.IssueCouponRequest{
				CampaignId: campaign.Id,
			}),
		)
		if err != nil {
			if connect.CodeOf(err) == connect.CodeResourceExhausted {
				log.Printf("쿠폰 발급 중단: %v", err)
				break
			}
			log.Fatalf("쿠폰 발급 실패: %v", err)
		}
		log.Printf("쿠폰 발급됨: ID=%d, 쿠폰=%s", issueRes.Msg.Coupon.Id, issueRes.Msg.Coupon.Code)
	}

	// 3-1. 쿠폰 발급 후 캠페인 상태 확인
	getRes, err = client.GetCampaign(
		ctx,
		connect.NewRequest(&couponv1.GetCampaignRequest{
			Id: campaign.Id,
		}),
	)
	if err != nil {
		log.Fatalf("캠페인 조회 실패: %v", err)
	}
	log.Printf("쿠폰 발급 후 캠페인 상태: ID=%d, 총 쿠폰=%d, 남은 쿠폰=%d",
		getRes.Msg.Campaign.Id, getRes.Msg.Campaign.TotalCoupons, getRes.Msg.Campaign.RemainingCoupons)

	// 4. 발급된 쿠폰 목록 조회
	listRes, err := client.GetIssuedCoupons(
		ctx,
		connect.NewRequest(&couponv1.GetIssuedCouponsRequest{
			CampaignId: campaign.Id,
		}),
	)
	if err != nil {
		log.Fatalf("쿠폰 목록 조회 실패: %v", err)
	}
	log.Printf("발급된 쿠폰 목록: %v", listRes.Msg.Codes)
}
