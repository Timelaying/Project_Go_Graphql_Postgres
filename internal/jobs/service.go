package jobs

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

var (
	ErrNotFound = errors.New("NOT_FOUND")
	ErrBadInput = errors.New("BAD_INPUT")
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, company, role string, link *string) (Job, error) {
	company = strings.TrimSpace(company)
	role = strings.TrimSpace(role)
	if company == "" || role == "" {
		return Job{}, ErrBadInput
	}
	normalizedLink, err := normalizeLink(link)
	if err != nil {
		return Job{}, err
	}
	job, err := s.repo.Create(ctx, company, role, normalizedLink)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Job, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrBadInput
	}
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return job, nil
}

func (s *Service) List(ctx context.Context, status *Status, q *string) ([]Job, error) {
	var query *string
	if q != nil {
		trimmed := strings.TrimSpace(*q)
		if trimmed != "" {
			query = &trimmed
		}
	}
	return s.repo.List(ctx, status, query)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status Status) (Job, error) {
	id = strings.TrimSpace(id)
	if id == "" || !validStatus(status) {
		return Job{}, ErrBadInput
	}
	job, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Job{}, ErrNotFound
		}
		return Job{}, err
	}
	return job, nil
}

func (s *Service) UpdateLink(ctx context.Context, id string, link string) (Job, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Job{}, ErrBadInput
	}
	normalizedLink, err := normalizeLink(&link)
	if err != nil {
		return Job{}, err
	}
	job, err := s.repo.UpdateLink(ctx, id, normalizedLink)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Job{}, ErrNotFound
		}
		return Job{}, err
	}
	return job, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrBadInput
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *Service) StatsByStatus(ctx context.Context) ([]StatusCount, error) {
	return s.repo.StatsByStatus(ctx)
}

func (s *Service) SeedDemoJobs(ctx context.Context) (int, error) {
	return s.repo.SeedDemoJobs(ctx)
}

func normalizeLink(link *string) (*string, error) {
	if link == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*link)
	if trimmed == "" {
		return nil, ErrBadInput
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, ErrBadInput
	}
	return &trimmed, nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusApplied, StatusInterview, StatusOffer, StatusRejected:
		return true
	default:
		return false
	}
}
