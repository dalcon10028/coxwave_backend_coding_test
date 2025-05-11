package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/dalcon10028/coxwave_backend_coding_test/internal/database"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := database.Connect(database.NewConfig(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestCampaignRepository_Create(t *testing.T) {
	// given
	db := setupDB(t)
	repo := NewCampaignRepository(db)
	ctx := context.Background()

	// when
	campaign, err := repo.Create(ctx, 100, time.Now().Unix())

	// then
	if err != nil {
		t.Fatal(err)
	}
	if campaign.TotalCoupons != 100 {
		t.Error("Expected total coupons to be 100, but got", campaign.TotalCoupons)
	}
	if campaign.RemainingCoupons != 100 {
		t.Error("Expected remaining coupons to be 100, but got", campaign.RemainingCoupons)
	}
}

func TestCampaignRepository_Get(t *testing.T) {
	// given
	db := setupDB(t)
	repo := NewCampaignRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, 100, time.Now().Unix())
	if err != nil {
		t.Fatal(err)
	}

	// when
	got, err := repo.Get(ctx, created.ID)

	// then
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Error("Expected campaign ID to be", created.ID, "but got", got.ID)
	}
	if got.TotalCoupons != 100 {
		t.Error("Expected total coupons to be 100, but got", got.TotalCoupons)
	}
}

func TestCampaignRepository_IssueCoupon(t *testing.T) {
	// given
	db := setupDB(t)
	repo := NewCampaignRepository(db)
	ctx := context.Background()

	campaign, err := repo.Create(ctx, 2, time.Now().Unix())
	if err != nil {
		t.Fatal(err)
	}

	// when
	err = repo.IssueCoupon(ctx, campaign.ID, "CODE1")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.IssueCoupon(ctx, campaign.ID, "CODE2")
	if err != nil {
		t.Fatal(err)
	}

	// then
	got, err := repo.Get(ctx, campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.RemainingCoupons != 0 {
		t.Error("Expected remaining coupons to be 0, but got", got.RemainingCoupons)
	}

	// 발급된 쿠폰 코드 확인
	codes, err := repo.GetIssuedCodes(ctx, campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 2 {
		fmt.Println(codes)
		t.Error("Expected 2 issued codes, but got", len(codes))
	}
}
