-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    price BIGINT NOT NULL CHECK (price >= 0), -- в минимальных денежных единицах (копейки), без плавающей точки
    category TEXT,
    is_limited BOOLEAN NOT NULL DEFAULT false,
    stock INT NOT NULL DEFAULT 1 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_items_on_created_at ON items(created_at);
CREATE INDEX idx_items_on_id ON items(id);

-- +goose Down
DROP TABLE IF EXISTS items;