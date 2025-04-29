-- +goose Up
-- +goose StatementBegin
create table file_changes (
    id int generated always as identity primary key,
    ulid text unique not null,
    repository_id int not null references repos (id),
    commit_id int not null references commits (id),
    file_path text not null,
    language text not null,
    lines_added int not null,
    lines_removed int not null,
    lines_changed int not null,
    commit_hash text not null,
    vendor_files boolean not null,
    generated_files boolean not null,
    commit_at timestamptz not null,
    synced_at timestamptz
);

create index idx_file_changes_repository_id on file_changes (repository_id);

create index idx_file_changes_commit_id on file_changes (commit_id);

create index idx_file_changes_language on file_changes (language);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table file_changes;

-- +goose StatementEnd
