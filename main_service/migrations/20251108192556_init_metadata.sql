-- +goose Up
CREATE TABLE IF NOT EXISTS pastes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  hash VARCHAR(64) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_pastes_hash ON pastes(hash);
CREATE INDEX idx_pastes_expires_at ON pastes(expires_at);
CREATE INDEX idx_pastes_expires_created ON pastes(expires_at, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_pastes_expires_created ON pastes;
DROP INDEX IF EXISTS idx_pastes_expires_at ON pastes;
DROP INDEX IF EXISTS idx_pastes_hash ON pastes;
DROP TABLE IF EXISTS pastes;
