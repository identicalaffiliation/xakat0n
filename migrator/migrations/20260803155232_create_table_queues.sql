-- +goose Up
CREATE TYPE queue_status AS ENUM (
    'QUEUED',
    'OFFERED',
    'EXPIRED',
    'COMPLETED',
    'SOLD_OUT',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS queues (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status queue_status NOT NULL DEFAULT 'QUEUED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_queue_product_status_created
    ON queues (product_id, status, created_at);

CREATE UNIQUE INDEX idx_queue_unique_active
    ON queues (product_id, user_id)
    WHERE status IN ('QUEUED', 'OFFERED');

CREATE UNIQUE INDEX idx_queue_one_offered
    ON queues (product_id)
    WHERE status = 'OFFERED';

-- +goose Down
DROP INDEX IF EXISTS idx_queue_product_status_position;
DROP INDEX IF EXISTS idx_queue_unique_active;
DROP TABLE IF EXISTS queues;
DROP TYPE IF EXISTS queue_status;
