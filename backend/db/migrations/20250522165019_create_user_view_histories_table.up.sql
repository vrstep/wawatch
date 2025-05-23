CREATE TABLE IF NOT EXISTS user_view_histories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INT NOT NULL,
    anime_external_id INT NOT NULL,
    last_viewed_at TIMESTAMPTZ,
    view_count INT DEFAULT 1,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT uq_user_anime_view UNIQUE (user_id, anime_external_id)
);

CREATE INDEX IF NOT EXISTS idx_user_view_histories_deleted_at ON user_view_histories (deleted_at);

CREATE INDEX IF NOT EXISTS idx_user_view_histories_last_viewed_at ON user_view_histories (last_viewed_at);