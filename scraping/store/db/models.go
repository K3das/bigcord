// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"github.com/K3das/bigcord/scraping/store/schema"
	"github.com/bwmarrin/discordgo"
)

type State struct {
	ChannelID     string
	GuildID       string
	Type          discordgo.ChannelType
	Name          string
	State         schema.State
	MessageOffset string
}
