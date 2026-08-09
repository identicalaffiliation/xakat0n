-- +goose Up
CREATE TYPE queue_status AS ENUM (
    'QUEUED',
    'OFFERED',
    'CHECKOUT',
    'EXPIRED',
    'PURCHASED',
    'SOLD_OUT',
    'CANCELLED'
    );

CREATE TABLE IF NOT EXISTS queues (
    id UUID PRIMARY KEY NOT NULL,
    item_id UUID NOT NULL REFERENCES items(id),
    user_id UUID NOT NULL,
    status queue_status NOT NULL DEFAULT 'QUEUED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_queue_unique_user_product
    ON queues (item_id, user_id)
    WHERE status IN ('QUEUED', 'OFFERED', 'CHECKOUT', 'SOLD_OUT');

CREATE INDEX idx_queue_product_status
    ON queues (item_id, status)
    WHERE status IN ('OFFERED', 'CHECKOUT');

CREATE INDEX idx_queue_product_status_created
    ON queues (item_id, status, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_queue_product_status_created;
DROP INDEX IF EXISTS idx_queue_product_status;
DROP INDEX IF EXISTS idx_queue_unique_user_product;
DROP TABLE IF EXISTS queues;
DROP TYPE IF EXISTS queue_status;

