-- +goose Up
-- +goose StatementBegin
SELECT 'stage2 placeholder';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'stage2 placeholder rollback';
-- +goose StatementEnd