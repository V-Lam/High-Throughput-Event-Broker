CREATE TYPE message_status AS ENUM ('pending', 'processing', 'completed', 'failed');

CREATE TABLE event_messages (
    id BIGSERIAL PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    payload TEXT NOT NULL,
    status message_status DEFAULT 'pending',
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_status_topic ON event_messages (status, topic);
CREATE INDEX idx_messages_created_at ON event_messages (created_at);
