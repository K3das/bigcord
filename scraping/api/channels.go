package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/K3das/bigcord/scraping/api/types"
	"github.com/K3das/bigcord/scraping/jobs"
	"github.com/gofiber/fiber/v2"
	"time"
)

type GetChannelsOutput struct {
	Channels []types.Channel     `json:"channels"`
	LostJobs map[string]jobs.Job `json:"lost_jobs"`
}

func (a *API) GetChannels(c *fiber.Ctx) error {
	var channels []types.Channel
	storeChannels, err := a.store.ListStates(c.Context())
	if err != nil {
		return fmt.Errorf("error listing states: %w", err)
	}

	jobStates := a.jobs.GetJobs()

	for _, channel := range storeChannels {
		channel := types.Channel{
			ID:      channel.ChannelID,
			GuildID: channel.GuildID,
			Type:    channel.Type,
			Name:    channel.Name,
			State:   channel.State,
			Offset:  channel.MessageOffset,
		}

		if job, ok := jobStates[channel.ID]; ok {
			channel.JobState = job.State
			if job.Error != nil {
				channel.JobError = job.Error.Error()
			}
			delete(jobStates, channel.ID)
		}

		channels = append(channels, channel)
	}

	return c.JSON(&types.GenericResponse{
		Success: true,
		Data: GetChannelsOutput{
			Channels: channels,
			LostJobs: jobStates,
		},
	})
}

func (a *API) GetChannelsChannel(c *fiber.Ctx) error {
	id := c.Params("channel_id")

	storeState, err := a.store.GetState(c.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		return &fiber.Error{
			Code:    404,
			Message: "Channel state not found",
		}
	} else if err != nil {
		return fmt.Errorf("error getting existing state: %w", err)
	}
	channel := &types.Channel{
		ID:      storeState.ChannelID,
		GuildID: storeState.GuildID,
		Type:    storeState.Type,
		Name:    storeState.Name,
		State:   storeState.State,
		Offset:  storeState.MessageOffset,
	}

	jobState, err := a.jobs.GetJob(id)
	if err == nil {
		channel.JobState = jobState.State
		if jobState.Error != nil {
			channel.JobError = jobState.Error.Error()
		}
	}

	return c.JSON(&types.GenericResponse{
		Success: true,
		Data:    channel,
	})
}

func (a *API) PostChannels(c *fiber.Ctx) error {
	var channelIDs []string
	if err := c.BodyParser(&channelIDs); err != nil {
		return fmt.Errorf("error parsing body: %w", err)
	}

	for _, c := range channelIDs {
		err := a.jobs.Add(c, func(ctx context.Context) error {
			c := c
			return a.archiver.ScrapeChannel(ctx, c)
		})
		if err != nil {
			return fmt.Errorf("failed to add job %s: %w", c, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return c.JSON(&types.GenericResponse{
		Success: true,
		Data:    nil,
	})
}
