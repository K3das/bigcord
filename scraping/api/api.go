package api

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/K3das/bigcord/scraping/api/types"
	"github.com/K3das/bigcord/scraping/archiver"
	"github.com/K3das/bigcord/scraping/jobs"
	"github.com/K3das/bigcord/scraping/store"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

type API struct {
	log *zap.SugaredLogger
	app *fiber.App

	archiver *archiver.Archiver
	store    *store.Store
	jobs     *jobs.JobStore
}

//go:embed ui/*
var embedUI embed.FS

func NewAPI(rawLog *zap.Logger, archiver *archiver.Archiver, store *store.Store, jobs *jobs.JobStore) *API {
	a := &API{
		archiver: archiver,
		store:    store,
		jobs:     jobs,
	}
	a.log = rawLog.Sugar().With("source", "api")

	a.app = fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := err.Error()

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				message = e.Message
			}

			err = ctx.Status(code).JSON(&types.GenericResponse{
				Success: false,
				Error:   message,
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})

	a.app.Use(fiberzap.New(fiberzap.Config{
		Logger: a.log.Desugar(),
		Levels: []zapcore.Level{zap.ErrorLevel, zap.WarnLevel, zap.DebugLevel},
	}))

	a.app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return os.Getenv("ENVIRONMENT") == "development"
		},
	}))

	api := a.app.Group("/api")
	api.Get("/discord/guilds/:guild_id<regex(\\d+)>", a.GetDiscordGuild)

	api.Get("/channels", a.GetChannels)
	api.Post("/channels", a.PostChannels)
	api.Get("/channels/:channel_id<regex(\\d+)>", a.GetChannelsChannel)
	
	api.Get("/jobs/:job_id<regex(\\d+)>", a.GetJobsJob)

	a.app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embedUI),
		PathPrefix: "ui",
	}))

	return a
}

func (a *API) Listen(ctx context.Context, addr string) error {
	err := a.app.ShutdownWithContext(ctx)
	if err != nil {
		return fmt.Errorf("error registering context: %w", err)
	}
	return a.app.Listen(addr)
}
