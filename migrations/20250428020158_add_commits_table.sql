-- +goose Up
-- +goose StatementBegin
create table commits (
    id int generated always as identity primary key,
    ulid text unique not null,
    commit_hash text not null,
    repository_id int not null references repos (id),
    author_name text not null,
    author_email text not null,
    message text not null,
    commit_at timestamptz not null,
    synced_at timestamptz
);

create index idx_commits_repository_id on commits (repository_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table commits;

-- +goose StatementEnd
