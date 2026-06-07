-- +goose Up
ALTER TABLE world_cup_users ADD COLUMN nickname VARCHAR(64);

-- +goose Down
ALTER TABLE world_cup_users DROP COLUMN IF EXISTS nickname;
