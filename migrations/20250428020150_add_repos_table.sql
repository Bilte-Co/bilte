-- +goose Up
-- +goose StatementBegin
create table repos (
    id int generated always as identity primary key,
    ulid text unique not null,
    repo_name text not null,
    repo_full_name text not null,
    repo_owner text not null,
    default_branch text not null,
    url text not null,
    synced_at timestamptz
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table repos;

-- +goose StatementEnd
