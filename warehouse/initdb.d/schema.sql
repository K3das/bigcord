CREATE TABLE IF NOT EXISTS messages
(
    `message_id` String,
    `channel_id` String,
    `author_id` String,

    `content` String,

    `created_at` DateTime('UTC'),
    `edited_at` DateTime('UTC'),

    `tts` Boolean,
    `mentions_everyone` Boolean,
    `user_mentions` Array(String),

    `attachments_json` String,
    `embeds_json` String,
    `sticker_items_json` String,

    `pinned` Boolean,

    `webhook_id` String,

    `type` Int64,

    `message_reference_json` String,
    `message_reference_id` String,
    `thread_id` String,

    `embed_links` Nested(url String, type String),
    `reactions` Nested(id String, name String, count Int64)
) ENGINE = ReplacingMergeTree ORDER BY (message_id, channel_id);

CREATE TABLE IF NOT EXISTS channels
(
    `channel_id` String,
    `guild_id` String,

    `type` Int64,

    `name` String,
    `topic` String,
    `nsfw` Boolean,

    `owner_id` String,
    `parent_id` String,

    `applied_tags` Array(String)
) ENGINE = ReplacingMergeTree ORDER BY channel_id;
