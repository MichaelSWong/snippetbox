-- +goose UP
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
  token CHAR(43) PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
