CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    subject_id UUID NOT NULL,
    subject_type VARCHAR(20) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_subject ON refresh_tokens(subject_id, subject_type);
CREATE UNIQUE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
