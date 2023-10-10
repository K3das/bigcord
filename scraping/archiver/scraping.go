package archiver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/K3das/bigcord/scraping/store/db"
	storeSchema "github.com/K3das/bigcord/scraping/store/schema"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func (a *Archiver) ScrapeChannel(ctx context.Context, c string) (err error) {
	channel, err := a.discord.Channel(c, discordgo.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error resolving channel: %w", err)
	}

	if channel.Type == discordgo.ChannelTypeGuildForum {
		return a.ScrapeForumChannel(ctx, channel)
	}

	log := a.log.With("channel_id", channel.ID).With("channel_type", ChannelTypeString(channel.Type))

	{
		// Apparently there is no better way to do this
		batch, err := a.warehouse.PrepareBatch(ctx, "INSERT INTO channels")
		if err != nil {
			return fmt.Errorf("error preparing batch: %w", err)
		}

		channelRow := TranslateDiscordChannel(channel)
		err = batch.AppendStruct(&channelRow)
		if err != nil {
			return fmt.Errorf("error appending channel to batch: %w", err)
		}

		err = batch.Send()
		if err != nil {
			return fmt.Errorf("error inserting channel: %w", err)
		}
	}

	maxID := "0"

	state, err := a.store.GetState(ctx, channel.ID)
	if err == nil {
		maxID = state.MessageOffset
		log.Infof("recovering state as %s", maxID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error getting existing state: %w", err)
	} else {
		log.Info("no in-progress state found, scraping from 0")
	}

	defer func() {
		err := a.store.SetStateIfNotCompleted(context.Background(), db.SetStateIfNotCompletedParams{
			ChannelID: channel.ID,
			State:     storeSchema.StateCrashed,
		})
		if err != nil {
			log.Error(fmt.Errorf("error setting crashed state: %w", err))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		log.Debugf("scraping after %s", maxID)
		maxID, err = a.ScrapePage(context.Background(), channel, maxID)
		if err != nil {
			return fmt.Errorf("error scraping page: %w", err)
		}

		if maxID == "0" {
			break
		}

		err = a.store.SetState(context.Background(), db.SetStateParams{
			GuildID:       channel.GuildID,
			ChannelID:     channel.ID,
			Type:          channel.Type,
			Name:          channel.Name,
			State:         storeSchema.StateInProgress,
			MessageOffset: maxID,
		})
		if err != nil {
			return fmt.Errorf("error updating state: %w", err)
		}
	}

	log.Debug("done scraping, marking as completed")
	if err := a.store.SetState(ctx, db.SetStateParams{
		GuildID:       channel.GuildID,
		ChannelID:     channel.ID,
		Type:          channel.Type,
		Name:          channel.Name,
		State:         storeSchema.StateCompleted,
		MessageOffset: "",
	}); err != nil {
		return fmt.Errorf("error marking state as completed: %w", err)
	}

	return nil
}

func (a *Archiver) ScrapePage(ctx context.Context, channel *discordgo.Channel, after string) (max string, err error) {
	channelMessages, err := a.discord.ChannelMessages(channel.ID, 100, "", after, "", discordgo.WithContext(ctx))
	if err != nil {
		return "0", fmt.Errorf("error fetching messages: %w", err)
	}

	batch, err := a.warehouse.PrepareBatch(ctx, "INSERT INTO messages")
	if err != nil {
		return "0", fmt.Errorf("error preparing batch: %w", err)
	}

	maxID := "0"
	for _, message := range channelMessages {
		{
			id := channel.ID
			if channel.Type == discordgo.ChannelTypeGuildPublicThread ||
				channel.Type == discordgo.ChannelTypeGuildPrivateThread {
				id = "P:" + channel.ParentID
			}
			messagesProcessedCounter.WithLabelValues(ChannelTypeString(channel.Type), id).Inc()
		}

		if message.Thread != nil {
			err := a.ScrapeChannel(ctx, message.Thread.ID)
			if err != nil {
				return "0", fmt.Errorf("error scraping thread: %w", err)
			}
		}

		messageRow, err := TranslateDiscordMessage(message)
		if err != nil {
			return "0", fmt.Errorf("error translating message: %w", err)
		}

		err = a.DownloadAttachments(ctx, message)
		if err != nil {
			return "", err
		}

		err = batch.AppendStruct(messageRow)
		if err != nil {
			return "0", fmt.Errorf("error appending message to batch: %w", err)
		}

		msgIDInt, err1 := strconv.ParseInt(message.ID, 10, 64)
		maxIDInt, err2 := strconv.ParseInt(maxID, 10, 64)

		if err1 != nil || err2 != nil {
			return "0", fmt.Errorf("error parsing snowflakes as integers: %w, %w", err1, err2)
		}

		if msgIDInt > maxIDInt {
			maxID = message.ID
		}
	}

	if maxID != "0" {
		err = batch.Send()
		if err != nil {
			return "0", fmt.Errorf("error inserting messages: %w", err)
		}
	}

	return maxID, nil
}

func (a *Archiver) ScrapeForumChannel(ctx context.Context, channel *discordgo.Channel) (err error) {
	log := a.log.With("guild_id", channel.GuildID).With("channel_id", channel.ID)

	log.Debug("scraping as forum channel")

	guildThreads, err := a.discord.GuildThreadsActive(channel.GuildID, discordgo.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error getting active threads: %w", err)
	}

	forumThreads := make(map[string]struct{})
	for _, thread := range guildThreads.Threads {
		if thread.ParentID != channel.ID {
			continue
		}

		if thread.ThreadMetadata == nil || thread.ThreadMetadata.Archived {
			continue
		}

		if _, ok := forumThreads[thread.ID]; ok {
			continue
		}

		forumThreads[thread.ID] = struct{}{}
	}

	log.Debugf("got %d active threads, getting archived threads", len(forumThreads))

	before := time.Now()
	for {
		log.Debugf("getting threads before %s", before.String())

		archived, err := a.discord.ThreadsArchived(channel.ID, &before, 100)
		if err != nil {
			return fmt.Errorf("error fetching archived threads: %w", err)
		}
		for _, thread := range archived.Threads {
			if _, ok := forumThreads[thread.ID]; ok {
				continue
			}

			if thread.ThreadMetadata != nil && thread.ThreadMetadata.ArchiveTimestamp.Before(before) {
				before = thread.ThreadMetadata.ArchiveTimestamp
			}

			forumThreads[thread.ID] = struct{}{}
		}

		if !archived.HasMore {
			break
		}
	}

	log.Debugf("finally have %d thread channels", len(forumThreads))

	for threadID := range forumThreads {
		err := a.ScrapeChannel(ctx, threadID)
		if err != nil {
			return fmt.Errorf("error scraping thread: %w", err)
		}
	}

	log.Debugf("done scraping forum threads")

	return nil
}
