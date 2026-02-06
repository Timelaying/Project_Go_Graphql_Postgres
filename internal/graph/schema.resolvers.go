package graph

import (
	"context"
	"errors"
	"time"

	"job-tracker/internal/graph/model"
	"job-tracker/internal/jobs"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateJob is the resolver for the createJob field.
func (r *mutationResolver) CreateJob(ctx context.Context, company string, role string, link *string) (*model.Job, error) {
	job, err := r.JobService.Create(ctx, company, role, link)
	if err != nil {
		return nil, mapServiceError(err)
	}
	return toModelJob(job), nil
}

// UpdateJobStatus is the resolver for the updateJobStatus field.
func (r *mutationResolver) UpdateJobStatus(ctx context.Context, id string, status model.Status) (*model.Job, error) {
	job, err := r.JobService.UpdateStatus(ctx, id, fromModelStatus(status))
	if err != nil {
		return nil, mapServiceError(err)
	}
	return toModelJob(job), nil
}

// UpdateJobLink is the resolver for the updateJobLink field.
func (r *mutationResolver) UpdateJobLink(ctx context.Context, id string, link string) (*model.Job, error) {
	job, err := r.JobService.UpdateLink(ctx, id, link)
	if err != nil {
		return nil, mapServiceError(err)
	}
	return toModelJob(job), nil
}

// DeleteJob is the resolver for the deleteJob field.
func (r *mutationResolver) DeleteJob(ctx context.Context, id string) (bool, error) {
	if err := r.JobService.Delete(ctx, id); err != nil {
		return false, mapServiceError(err)
	}
	return true, nil
}

// SeedDemoJobs is the resolver for the seedDemoJobs field.
func (r *mutationResolver) SeedDemoJobs(ctx context.Context) (int, error) {
	count, err := r.JobService.SeedDemoJobs(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Jobs is the resolver for the jobs field.
func (r *queryResolver) Jobs(ctx context.Context, status *model.Status, q *string) ([]*model.Job, error) {
	var statusFilter *jobs.Status
	if status != nil {
		converted := fromModelStatus(*status)
		statusFilter = &converted
	}
	jobsList, err := r.JobService.List(ctx, statusFilter, q)
	if err != nil {
		return nil, err
	}
	out := make([]*model.Job, 0, len(jobsList))
	for _, job := range jobsList {
		out = append(out, toModelJob(job))
	}
	return out, nil
}

// Job is the resolver for the job field.
func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	job, err := r.JobService.Get(ctx, id)
	if err != nil {
		return nil, mapServiceError(err)
	}
	if job == nil {
		return nil, nil
	}
	return toModelJob(*job), nil
}

// StatsByStatus is the resolver for the statsByStatus field.
func (r *queryResolver) StatsByStatus(ctx context.Context) ([]*model.StatusCount, error) {
	stats, err := r.JobService.StatsByStatus(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*model.StatusCount, 0, len(stats))
	for _, stat := range stats {
		out = append(out, &model.StatusCount{
			Status: toModelStatus(stat.Status),
			Count:  stat.Count,
		})
	}
	return out, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

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
