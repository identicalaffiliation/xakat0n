-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    price BIGINT NOT NULL, -- в минимальных денежных единицах (копейки), без плавающей точки
    category TEXT,
    is_limited BOOLEAN NOT NULL DEFAULT false,
    stock INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS items;
