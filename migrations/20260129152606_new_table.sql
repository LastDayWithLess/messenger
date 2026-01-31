-- +goose Up
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP(0)
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    text VARCHAR(5000) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP(0)
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_chats_created_at ON chats(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_chats_created_at;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
