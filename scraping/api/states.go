package api

import (
	"github.com/K3das/bigcord/scraping/api/types"
	"github.com/gofiber/fiber/v2"
)

func (a *API) GetStatesState(c *fiber.Ctx) error {
	id := c.Params("state_id")

	jobState, err := a.jobs.GetJob(id)
	if err != nil {
		return &fiber.Error{
			Code:    404,
			Message: "Job not found",
		}
	}

	job := &types.Job{
		ID:    id,
		State: jobState.State,
	}

	if jobState.Error != nil {
		job.JobError = jobState.Error.Error()
	}

	return c.JSON(&types.GenericResponse{
		Success: true,
		Data:    job,
	})
}
