-- +goose Up
-- Establish Goose version tracking before domain tables are introduced.
-- No application tables are needed for the runnable foundation.
SELECT 1;

-- +goose Down
-- Version zero preserves existing application data and Goose's history table.
SELECT 1;
