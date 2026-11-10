package graph

import (
	"errors"
	"time"

	"job-tracker/internal/graph/model"
	"job-tracker/internal/jobs"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func toModelJob(job jobs.Job) *model.Job {
	return &model.Job{
		ID:        job.ID,
		Company:   job.Company,
		Role:      job.Role,
		Link:      job.Link,
		Status:    toModelStatus(job.Status),
		CreatedAt: job.CreatedAt.Format(time.RFC3339),
	}
}

func toModelStatus(status jobs.Status) model.Status {
	switch status {
	case jobs.StatusApplied:
		return model.StatusApplied
	case jobs.StatusInterview:
		return model.StatusInterview
	case jobs.StatusOffer:
		return model.StatusOffer
	case jobs.StatusRejected:
		return model.StatusRejected
	default:
		return model.StatusApplied
	}
}

func fromModelStatus(status model.Status) jobs.Status {
	switch status {
	case model.StatusApplied:
		return jobs.StatusApplied
	case model.StatusInterview:
		return jobs.StatusInterview
	case model.StatusOffer:
		return jobs.StatusOffer
	case model.StatusRejected:
		return jobs.StatusRejected
	default:
		return jobs.StatusApplied
	}
}

func mapServiceError(err error) error {
	if errors.Is(err, jobs.ErrBadInput) {
		return &gqlerror.Error{
			Message:    "bad input",
			Extensions: map[string]any{"code": "BAD_INPUT"},
		}
	}
	if errors.Is(err, jobs.ErrNotFound) {
		return &gqlerror.Error{
			Message:    "not found",
			Extensions: map[string]any{"code": "NOT_FOUND"},
		}
	}
	return err
}
