package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/dalcon10028/coxwave_backend_coding_test/internal/model"
)

type CampaignRepository struct {
	db *sql.DB
}

func NewCampaignRepository(db *sql.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

func (r *CampaignRepository) Create(ctx context.Context, totalCoupons int64, startAt int64) (*model.Campaign, error) {
	query := `
		INSERT INTO campaign (total_coupons, remaining_coupons, start_at, created_at)
		VALUES (?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, totalCoupons, totalCoupons, time.Unix(startAt, 0), now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Campaign{
		ID:               id,
		TotalCoupons:     totalCoupons,
		RemainingCoupons: totalCoupons,
		StartAt:          time.Unix(startAt, 0),
		CreatedAt:        now,
	}, nil
}

func (r *CampaignRepository) Get(ctx context.Context, id int64) (*model.Campaign, error) {
	query := `
		SELECT id, total_coupons, remaining_coupons, start_at, created_at
		FROM campaign
		WHERE id = ?
	`
	var campaign model.Campaign
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&campaign.ID,
		&campaign.TotalCoupons,
		&campaign.RemainingCoupons,
		&campaign.StartAt,
		&campaign.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

const (
	// 한글 유니코드 범위: 가(0xAC00) ~ 힣(0xD7A3)
	hangulStart = 0xAC00
	hangulEnd   = 0xD7A3
	numbers     = "0123456789"
)

func generateCouponCode(campaignID int64) string {
	rand.Seed(time.Now().UnixNano())

	hangulPart := make([]rune, 2)
	for i := range hangulPart {
		hangulPart[i] = rune(hangulStart + rand.Intn(hangulEnd-hangulStart+1))
	}

	numPart := make([]byte, 4)
	for i := range numPart {
		numPart[i] = numbers[rand.Intn(len(numbers))]
	}

	return fmt.Sprintf("%d%s%s", campaignID, string(hangulPart), string(numPart))
}

func (r *CampaignRepository) IssueCoupon(ctx context.Context, campaignID int64) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 남은 쿠폰 수 감소
	result, err := tx.ExecContext(ctx, `
		UPDATE campaign 
		SET remaining_coupons = remaining_coupons - 1
		WHERE id = ? AND remaining_coupons > 0 AND start_at <= ?
	`, campaignID, time.Now())
	if err != nil {
		return "", err
	}

	// 업데이트된 행이 있는지 확인
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", fmt.Errorf("failed to issue coupon, campaign id: %d", campaignID)
	}

	// 쿠폰 코드 생성
	code := generateCouponCode(campaignID)

	// 쿠폰 발급
	_, err = tx.ExecContext(ctx, `
		INSERT INTO coupon (campaign_id, code, created_at)
		VALUES (?, ?, ?)
	`, campaignID, code, time.Now())
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return code, nil
}

func (r *CampaignRepository) GetIssuedCodes(ctx context.Context, campaignID int64) ([]string, error) {
	query := `
		SELECT code
		FROM coupon
		WHERE campaign_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}
