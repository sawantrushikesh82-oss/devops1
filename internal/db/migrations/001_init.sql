-- +goose Up
CREATE TYPE user_type AS ENUM ('external','internal');
CREATE TYPE user_role AS ENUM ('admin','viewer');
CREATE TYPE flow_type AS ENUM ('send_receive','send_only','receive_only');
CREATE TYPE log_status AS ENUM ('success','failed','access_denied');

CREATE TABLE users (
  id UUID PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  public_key TEXT NOT NULL DEFAULT '',
  user_type user_type NOT NULL,
  role user_role NOT NULL DEFAULT 'viewer',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mappings (
  id UUID PRIMARY KEY,
  external_user TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
  internal_user TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
  flow_type flow_type NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (external_user, internal_user)
);

CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  refresh_token TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE logs (
  id UUID PRIMARY KEY,
  username TEXT NOT NULL,
  filename TEXT NOT NULL,
  direction TEXT NOT NULL,
  bytes BIGINT NOT NULL DEFAULT 0,
  status log_status NOT NULL,
  error TEXT NOT NULL DEFAULT '',
  mapping_id UUID,
  session_id UUID,
  component TEXT NOT NULL,
  duration_ms BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_type ON users(user_type);
CREATE INDEX idx_users_active ON users(is_active);
CREATE INDEX idx_mappings_active ON mappings(is_active);
CREATE INDEX idx_mappings_flow ON mappings(flow_type);
CREATE INDEX idx_logs_username ON logs(username);
CREATE INDEX idx_logs_status ON logs(status);
CREATE INDEX idx_logs_created_at ON logs(created_at DESC);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- +goose Down
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS mappings;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS log_status;
DROP TYPE IF EXISTS flow_type;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_type;
