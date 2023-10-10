package types

import (
	"github.com/K3das/bigcord/scraping/jobs"
	"github.com/K3das/bigcord/scraping/store/schema"
	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	ID      string                `json:"id"`
	GuildID string                `json:"guild_id"`
	Type    discordgo.ChannelType `json:"type"`

	Name string `json:"name,omitempty"`

	ParentID  string `json:"parent_id,omitempty"`
	HasAccess bool   `json:"has_access,omitempty"`

	State  schema.State `json:"state,omitempty"`
	Offset string       `json:"offset,omitempty"`

	JobState jobs.JobState `json:"job_state,omitempty"`
	JobError string        `json:"job_error,omitempty"`
}

type Guild struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Channels []Channel `json:"channels"`
}

type GenericResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Job struct {
	ID string `json:"id"`

	State    jobs.JobState `json:"state,omitempty"`
	JobError string        `json:"job_error,omitempty"`
}
