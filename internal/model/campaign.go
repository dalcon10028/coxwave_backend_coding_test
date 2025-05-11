package model

import "time"

type Campaign struct {
	ID               int64
	TotalCoupons     int64
	RemainingCoupons int64
	StartAt          time.Time
	CreatedAt        time.Time
}
