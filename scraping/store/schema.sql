CREATE TABLE IF NOT EXISTS states
(
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    type INT NOT NULL,
    name TEXT NOT NULL,
    state INT DEFAULT 1 NOT NULL,
    message_offset TEXT DEFAULT '0' NOT NULL,
    PRIMARY KEY (channel_id)
)

