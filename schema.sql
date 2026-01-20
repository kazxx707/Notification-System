-- Subscription table
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    channels TEXT NOT NULL, -- JSON array: ["email", "sms"]
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING or NOTIFIED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, item_id) -- One subscription per user per item
);

-- Notification table
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    channel VARCHAR(20) NOT NULL, -- email, sms, push
    status VARCHAR(20) NOT NULL, -- SUCCESS or FAILED
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_subscriptions_item_status ON subscriptions(item_id, status);
CREATE INDEX IF NOT EXISTS idx_notifications_user_item ON notifications(user_id, item_id);
