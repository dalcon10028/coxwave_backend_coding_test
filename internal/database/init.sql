CREATE TABLE IF NOT EXISTS campaign (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    total_coupons INTEGER NOT NULL,
    remaining_coupons INTEGER NOT NULL,
    start_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS coupon (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    campaign_id INTEGER NOT NULL,
    code TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_coupon_campaign_id ON coupon(campaign_id);
CREATE INDEX IF NOT EXISTS idx_coupon_code ON coupon(code); 