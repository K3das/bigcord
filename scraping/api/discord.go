package api

import (
	"github.com/K3das/bigcord/scraping/api/types"
	"github.com/gofiber/fiber/v2"
)

func (a *API) GetDiscordGuild(c *fiber.Ctx) error {
	id := c.Params("guild_id")

	guild, err := a.archiver.GetGuild(c.Context(), id)
	if err != nil {
		return err
	}

	jobStates := a.jobs.GetJobs()
	for i, channel := range guild.Channels {
		j := jobStates[channel.ID]
		guild.Channels[i].JobState = j.State
	}

	return c.JSON(&types.GenericResponse{
		Success: true,
		Data:    guild,
	})
}
