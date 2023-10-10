package archiver

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (a *Archiver) DownloadAttachments(ctx context.Context, message *discordgo.Message) error {
	for i, attachment := range message.Attachments {
		err := a.downloadAttachment(ctx, i, message, attachment)
		if err != nil {
			return fmt.Errorf("error downloading attachment: %w", err)
		}
	}
	return nil
}

// TODO: move this somewhere else, maybe discord.go

func isDiscordHostname(hostname string) bool {
	hostnames := []string{
		"cdn.discordapp.com",
		"media.discordapp.net",
	}
	for _, discordHostname := range hostnames {
		if strings.HasSuffix(hostname, discordHostname) {
			return true
		}
	}
	return false
}

func truncateFilename(filename string) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	if len(nameWithoutExt) > 20 {
		nameWithoutExt = nameWithoutExt[:20]
	}

	return nameWithoutExt + ext
}
func (a *Archiver) downloadAttachment(ctx context.Context, index int, message *discordgo.Message, attachment *discordgo.MessageAttachment) (err error) {
	downloadURL := attachment.URL

	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return fmt.Errorf("error parsing attachment URL: %w", err)
	}

	if !isDiscordHostname(parsedURL.Hostname()) {
		return nil
	}

	// TODO: don't hardcode path
	downloadPath := filepath.Join("/data/media", message.ChannelID, message.ID, parsedURL.Hostname(), fmt.Sprintf("%d-%s", index, truncateFilename(attachment.Filename)))

	a.log.Debugf("downloading %s to %s", downloadURL, downloadPath)

	err = os.MkdirAll(filepath.Dir(downloadPath), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			a.log.Error(fmt.Errorf("error closing file writer: %w", err))
		}
	}(file)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			a.log.Error(fmt.Errorf("error closing response reader: %w", err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		a.log.Errorf("failed to download attachment %s, status code: %d", downloadURL, resp.StatusCode)
		return nil
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
