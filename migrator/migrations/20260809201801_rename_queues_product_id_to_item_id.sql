-- +goose Up
ALTER TABLE queues RENAME COLUMN product_id TO item_id;

ALTER INDEX idx_queue_unique_user_product RENAME TO idx_queue_unique_user_item;
ALTER INDEX idx_queue_product_status RENAME TO idx_queue_item_status;
ALTER INDEX idx_queue_product_status_created RENAME TO idx_queue_item_status_created;

-- +goose Down
ALTER INDEX idx_queue_item_status_created RENAME TO idx_queue_product_status_created;
ALTER INDEX idx_queue_item_status RENAME TO idx_queue_product_status;
ALTER INDEX idx_queue_unique_user_item RENAME TO idx_queue_unique_user_product;

ALTER TABLE queues RENAME COLUMN item_id TO product_id;
