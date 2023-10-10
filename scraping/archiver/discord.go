package archiver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	api "github.com/K3das/bigcord/scraping/api/types"
	"github.com/K3das/bigcord/scraping/warehouse"
	"github.com/bwmarrin/discordgo"
)

func ChannelTypeString(channelType discordgo.ChannelType) string {
	switch channelType {
	case discordgo.ChannelTypeGuildText:
		return "GUILD_TEXT"
	case discordgo.ChannelTypeDM:
		return "DM"
	case discordgo.ChannelTypeGuildVoice:
		return "GUILD_VOICE"
	case discordgo.ChannelTypeGroupDM:
		return "GROUP_DM"
	case discordgo.ChannelTypeGuildCategory:
		return "GUILD_CATEGORY"
	case discordgo.ChannelTypeGuildNews:
		return "GUILD_NEWS"
	case discordgo.ChannelTypeGuildStore:
		return "GUILD_STORE"
	case discordgo.ChannelTypeGuildNewsThread:
		return "GUILD_NEWS_THREAD"
	case discordgo.ChannelTypeGuildPublicThread:
		return "GUILD_PUBLIC_THREAD"
	case discordgo.ChannelTypeGuildPrivateThread:
		return "GUILD_PRIVATE_THREAD"
	case discordgo.ChannelTypeGuildStageVoice:
		return "GUILD_STAGE_VOICE"
	case discordgo.ChannelTypeGuildForum:
		return "GUILD_FORUM"
	default:
		return "UNKNOWN"
	}
}

type SearchResponse struct {
	TotalResults int `json:"total_results"`
}

func (a *Archiver) CountMessages(ctx context.Context, guild, channel string) (totalMessages int, err error) {
	log := a.log.With("guild_id", guild).With("channel_id", channel)
	endpoint := discordgo.EndpointGuilds + fmt.Sprintf("%s/messages/search", guild)
	res, err := a.discord.RequestWithBucketID(
		"GET",
		endpoint+fmt.Sprintf("?min_id=0&include_nsfw=true&limit=1&channel_id=%s", channel),
		nil,
		endpoint,
		discordgo.WithContext(ctx),
	)
	if err != nil {
		return 0, fmt.Errorf("error fetching message count: %w", err)
	}

	data := &SearchResponse{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return 0, fmt.Errorf("error parsing search response: %w", err)
	}

	log.Debugf("got message counts %d", data.TotalResults)

	return data.TotalResults, nil
}

type Sticker struct {
	ID          string                  `json:"id"`
	PackID      string                  `json:"pack_id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Tags        string                  `json:"tags"`
	Type        discordgo.StickerType   `json:"type"`
	FormatType  discordgo.StickerFormat `json:"format_type"`
	Available   bool                    `json:"available"`
	GuildID     string                  `json:"guild_id"`
	SortValue   int                     `json:"sort_value"`
}

func TranslateDiscordStickers(input []*discordgo.Sticker) []*Sticker {
	output := make([]*Sticker, 0, len(input))

	for _, s := range input {
		sticker := &Sticker{
			ID:          s.ID,
			PackID:      s.PackID,
			Name:        s.Name,
			Description: s.Description,
			Tags:        s.Tags,
			Type:        s.Type,
			FormatType:  s.FormatType,
			Available:   s.Available,
			GuildID:     s.GuildID,
			SortValue:   s.SortValue,
		}
		output = append(output, sticker)
	}

	return output
}

func TranslateDiscordMessage(dgMessage *discordgo.Message) (*warehouse.Message, error) {
	createdAt, err := discordgo.SnowflakeTimestamp(dgMessage.ID)
	if err != nil {
		return nil, fmt.Errorf("error converting snowflake to timestamp: %w", err)
	}

	editedAt := sql.NullTime{}
	if dgMessage.EditedTimestamp != nil {
		editedAt = sql.NullTime{Time: *dgMessage.EditedTimestamp, Valid: true}
	}

	userMentions := make([]string, len(dgMessage.Mentions))
	for i, mention := range dgMessage.Mentions {
		userMentions[i] = mention.ID
	}

	embedLinksURL := make([]string, len(dgMessage.Embeds))
	embedLinksType := make([]string, len(dgMessage.Embeds))
	for i, embed := range dgMessage.Embeds {
		if embed.URL == "" {
			continue
		}
		embedLinksURL[i] = embed.URL
		embedLinksType[i] = string(embed.Type)
	}

	reactionsID := make([]string, len(dgMessage.Reactions))
	reactionsName := make([]string, len(dgMessage.Reactions))
	reactionsCount := make([]int64, len(dgMessage.Reactions))
	for i, reaction := range dgMessage.Reactions {
		reactionsID[i] = reaction.Emoji.ID
		reactionsName[i] = reaction.Emoji.Name
		reactionsCount[i] = int64(reaction.Count)
	}

	message := &warehouse.Message{
		MessageID:        dgMessage.ID,
		ChannelID:        dgMessage.ChannelID,
		Content:          dgMessage.Content,
		CreatedAt:        createdAt,
		EditedAt:         editedAt,
		TTS:              dgMessage.TTS,
		MentionsEveryone: dgMessage.MentionEveryone,
		UserMentions:     userMentions,
		Pinned:           dgMessage.Pinned,
		WebhookID:        dgMessage.WebhookID,
		Type:             int64(dgMessage.Type),
		EmbedLinksUrl:    embedLinksURL,
		EmbedLinksType:   embedLinksType,
		ReactionsID:      reactionsID,
		ReactionsName:    reactionsName,
		ReactionsCount:   reactionsCount,
	}

	if dgMessage.Attachments != nil {
		attachmentsJSON, err := json.Marshal(dgMessage.Attachments)
		if err != nil {
			return nil, fmt.Errorf("error marshaling attachments to json: %w", err)
		}
		message.AttachmentsJSON = string(attachmentsJSON)
	}

	if dgMessage.Embeds != nil {
		embedsJSON, err := json.Marshal(dgMessage.Embeds)
		if err != nil {
			return nil, fmt.Errorf("error marshaling embeds to json: %w", err)
		}
		message.StickerItemsJSON = string(embedsJSON)
	}

	if dgMessage.StickerItems != nil {
		stickerItemsJSON, err := json.Marshal(TranslateDiscordStickers(dgMessage.StickerItems))
		if err != nil {
			return nil, fmt.Errorf("error marshaling sticker items to json: %w", err)
		}
		message.StickerItemsJSON = string(stickerItemsJSON)
	}

	if dgMessage.Author != nil {
		message.AuthorID = dgMessage.Author.ID
	}

	if dgMessage.Thread != nil {
		message.ThreadID = dgMessage.Thread.ID
	}

	if dgMessage.MessageReference != nil {
		message.MessageReferenceID = dgMessage.MessageReference.MessageID

		messageReferenceJSON, err := json.Marshal(dgMessage.MessageReference)
		if err != nil {
			return nil, fmt.Errorf("error marshaling message reference to JSON: %w", err)
		}
		message.MessageReferenceJSON = string(messageReferenceJSON)
	}

	return message, nil
}

func TranslateDiscordChannel(discordChannel *discordgo.Channel) warehouse.Channel {
	channel := warehouse.Channel{
		ChannelID:   discordChannel.ID,
		GuildID:     discordChannel.GuildID,
		Type:        int64(discordChannel.Type),
		Name:        discordChannel.Name,
		Topic:       discordChannel.Topic,
		NSFW:        discordChannel.NSFW,
		OwnerID:     discordChannel.OwnerID,
		ParentID:    discordChannel.ParentID,
		AppliedTags: discordChannel.AppliedTags,
	}

	return channel
}

func (a *Archiver) GetGuild(ctx context.Context, guildID string) (guild *api.Guild, err error) {
	discordGuild, err := a.discord.Guild(guildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting guild: %w", err)
	}

	channels, err := a.GetGuildChannels(ctx, guildID)
	if err != nil {
		return nil, err
	}

	return &api.Guild{
		ID:       discordGuild.ID,
		Name:     discordGuild.Name,
		Channels: channels,
	}, nil
}

func (a *Archiver) GetGuildChannels(ctx context.Context, guildID string) (channels []api.Channel, err error) {
	guildChannels, err := a.discord.GuildChannels(guildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting channels: %w", err)
	}

	channels = make([]api.Channel, len(guildChannels))
	for i, channel := range guildChannels {
		// this is probably not a great way to do this, *but it works*
		_, err := a.discord.ChannelMessages(channel.ID, 1, "", "0", "")

		channels[i] = api.Channel{
			ID:        channel.ID,
			GuildID:   channel.GuildID,
			Type:      channel.Type,
			Name:      channel.Name,
			ParentID:  channel.ParentID,
			HasAccess: err == nil,
		}
	}

	return channels, nil
}
