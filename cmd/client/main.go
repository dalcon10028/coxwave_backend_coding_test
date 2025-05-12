package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"

	couponv1 "github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1"
	"github.com/dalcon10028/coxwave_backend_coding_test/gen/coupon/v1/couponv1connect"
)

func main() {
	client := couponv1connect.NewCouponServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)

	ctx := context.Background()

	// 캠페인 생성 (쿠폰 1000개)
	campaign, err := client.CreateCampaign(ctx, connect.NewRequest(&couponv1.CreateCampaignRequest{
		TotalCoupons: 1000,
		StartAt:      time.Now().Unix(),
	}))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("캠페인 생성 완료: ID=%d, 총 쿠폰=%d\n",
		campaign.Msg.Campaign.Id,
		campaign.Msg.Campaign.TotalCoupons,
	)

	// 동시성 테스트
	log.Println("동시성 테스트 시작...")
	start := time.Now()

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex
	errors := make(map[string]int)

	// 2000개의 고루틴으로 동시 요청 (캠페인은 1000개만 발급 가능)
	const totalRequests = 2000

	// 모든 요청을 동시에 시작
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			code, err := client.IssueCoupon(ctx, connect.NewRequest(&couponv1.IssueCouponRequest{
				CampaignId: campaign.Msg.Campaign.Id,
			}))
			if err != nil {
				mu.Lock()
				errors[err.Error()]++
				mu.Unlock()
				return
			}

			mu.Lock()
			successCount++
			mu.Unlock()
			log.Printf("쿠폰 발급 성공: %s\n", code.Msg.Coupon.Code)
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	// 결과 출력
	log.Printf("\n=== 동시성 테스트 결과 ===\n")
	log.Printf("총 실행 시간: %v\n", totalDuration)
	log.Printf("초당 처리량: %.2f req/s\n", float64(totalRequests)/totalDuration.Seconds())
	log.Printf("성공: %d건\n", successCount)
	log.Printf("실패: %d건\n", totalRequests-successCount)
	log.Printf("에러 종류:\n")
	for err, count := range errors {
		log.Printf("- %s: %d건\n", err, count)
	}

	// 발급된 쿠폰 목록 확인
	codes, err := client.GetIssuedCoupons(ctx, connect.NewRequest(&couponv1.GetIssuedCouponsRequest{
		CampaignId: campaign.Msg.Campaign.Id,
	}))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n발급된 쿠폰 수: %d\n", len(codes.Msg.Codes))
	if len(codes.Msg.Codes) > 0 {
		log.Println("발급된 쿠폰 코드 (처음 10개):")
		for i, code := range codes.Msg.Codes {
			if i >= 10 {
				break
			}
			log.Printf("- %s\n", code)
		}
		if len(codes.Msg.Codes) > 10 {
			log.Printf("... 외 %d개\n", len(codes.Msg.Codes)-10)
		}
	}
}
