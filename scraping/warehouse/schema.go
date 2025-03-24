package warehouse

import (
	"database/sql"
	"time"
)

type Message struct {
	MessageID            string       `ch:"message_id"`
	ChannelID            string       `ch:"channel_id"`
	AuthorID             string       `ch:"author_id"`
	Content              string       `ch:"content"`
	CreatedAt            time.Time    `ch:"created_at"`
	EditedAt             sql.NullTime `ch:"edited_at"`
	TTS                  bool         `ch:"tts"`
	MentionsEveryone     bool         `ch:"mentions_everyone"`
	UserMentions         []string     `ch:"user_mentions"`
	AttachmentsJSON      string       `ch:"attachments_json"`
	EmbedsJSON           string       `ch:"embeds_json"`
	StickerItemsJSON     string       `ch:"sticker_items_json"`
	Pinned               bool         `ch:"pinned"`
	WebhookID            string       `ch:"webhook_id"`
	Type                 int64        `ch:"type"`
	MessageReferenceJSON string       `ch:"message_reference_json"`
	MessageReferenceID   string       `ch:"message_reference_id"`
	ThreadID             string       `ch:"thread_id"`
	EmbedLinksUrl        []string     `ch:"embed_links.url"`
	EmbedLinksType       []string     `ch:"embed_links.type"`
	ReactionsID          []string     `ch:"reactions.id"`
	ReactionsName        []string     `ch:"reactions.name"`
	ReactionsCount       []int64      `ch:"reactions.count"`
	FullJson             string       `ch:"full_json"`
}

type Channel struct {
	ChannelID   string   `ch:"channel_id"`
	GuildID     string   `ch:"guild_id"`
	Type        int64    `ch:"type"`
	Name        string   `ch:"name"`
	Topic       string   `ch:"topic"`
	NSFW        bool     `ch:"nsfw"`
	OwnerID     string   `ch:"owner_id"`
	ParentID    string   `ch:"parent_id"`
	AppliedTags []string `ch:"applied_tags"`
	FullJson    string   `ch:"full_json"`
}
