-- +goose Up
ALTER TABLE events DROP CONSTRAINT events_type_check;
ALTER TABLE events ADD CONSTRAINT events_type_check
    CHECK (type IN ('match_score', 'champion', 'runner_up'));

-- +goose Down
ALTER TABLE events DROP CONSTRAINT events_type_check;
ALTER TABLE events ADD CONSTRAINT events_type_check
    CHECK (type IN ('match_score', 'champion'));
