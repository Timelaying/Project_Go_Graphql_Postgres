package jobs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ pool *pgxpool.Pool }

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

func (r *Repo) Create(ctx context.Context, company, role string, link *string) (Job, error) {
	var j Job
	err := r.pool.QueryRow(ctx,
		`INSERT INTO jobs(company, role, link, status)
		 VALUES ($1,$2,$3,'APPLIED')
		 RETURNING id, company, role, link, status, created_at`,
		company, role, link,
	).Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt)
	return j, err
}

func (r *Repo) Get(ctx context.Context, id string) (*Job, error) {
	var j Job
	err := r.pool.QueryRow(ctx,
		`SELECT id, company, role, link, status, created_at
		 FROM jobs WHERE id=$1`, id,
	).Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &j, nil
}

func (r *Repo) List(ctx context.Context, status *Status, q *string) ([]Job, error) {
	base := `SELECT id, company, role, link, status, created_at FROM jobs`
	where := []string{}
	args := []any{}
	i := 1

	if status != nil {
		where = append(where, `status = $`+strconv.Itoa(i))
		args = append(args, string(*status))
		i++
	}
	if q != nil && strings.TrimSpace(*q) != "" {
		where = append(where, `(company ILIKE $`+strconv.Itoa(i)+` OR role ILIKE $`+strconv.Itoa(i)+`)`)
		args = append(args, "%"+strings.TrimSpace(*q)+"%")
		i++
	}

	if len(where) > 0 {
		base += " WHERE " + strings.Join(where, " AND ")
	}
	base += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Job{}
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

func (r *Repo) UpdateStatus(ctx context.Context, id string, status Status) (Job, error) {
	var j Job
	err := r.pool.QueryRow(ctx,
		`UPDATE jobs
		 SET status = $2
		 WHERE id = $1
		 RETURNING id, company, role, link, status, created_at`,
		id, status,
	).Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Job{}, ErrNotFound
		}
		return Job{}, err
	}
	return j, nil
}

func (r *Repo) UpdateLink(ctx context.Context, id string, link *string) (Job, error) {
	var j Job
	err := r.pool.QueryRow(ctx,
		`UPDATE jobs
		 SET link = $2
		 WHERE id = $1
		 RETURNING id, company, role, link, status, created_at`,
		id, link,
	).Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Job{}, ErrNotFound
		}
		return Job{}, err
	}
	return j, nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	commandTag, err := r.pool.Exec(ctx, `DELETE FROM jobs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) StatsByStatus(ctx context.Context) ([]StatusCount, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM jobs GROUP BY status ORDER BY status`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StatusCount
	for rows.Next() {
		var status Status
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		out = append(out, StatusCount{Status: status, Count: count})
	}
	return out, rows.Err()
}

func (r *Repo) SeedDemoJobs(ctx context.Context) (int, error) {
	type seed struct {
		company string
		role    string
		link    *string
		status  Status
	}

	seeds := []seed{
		{company: "Monzo", role: "Backend Engineer", link: strPtr("https://example.com/monzo"), status: StatusApplied},
		{company: "Stripe", role: "Platform Engineer", link: strPtr("https://example.com/stripe"), status: StatusInterview},
		{company: "Shopify", role: "Software Engineer", link: strPtr("https://example.com/shopify"), status: StatusOffer},
		{company: "Spotify", role: "Data Engineer", link: strPtr("https://example.com/spotify"), status: StatusRejected},
		{company: "GitHub", role: "Full Stack Engineer", link: strPtr("https://example.com/github"), status: StatusApplied},
		{company: "Airbnb", role: "Backend Engineer", link: strPtr("https://example.com/airbnb"), status: StatusInterview},
		{company: "Dropbox", role: "Infrastructure Engineer", link: strPtr("https://example.com/dropbox"), status: StatusApplied},
		{company: "Notion", role: "Product Engineer", link: strPtr("https://example.com/notion"), status: StatusOffer},
		{company: "Linear", role: "Frontend Engineer", link: strPtr("https://example.com/linear"), status: StatusRejected},
		{company: "Figma", role: "Growth Engineer", link: strPtr("https://example.com/figma"), status: StatusApplied},
	}

	inserted := 0
	for _, seed := range seeds {
		tag, err := r.pool.Exec(ctx,
			`INSERT INTO jobs (company, role, link, status)
			 VALUES ($1, $2, $3, $4)`,
			seed.company, seed.role, seed.link, seed.status,
		)
		if err != nil {
			return inserted, err
		}
		inserted += int(tag.RowsAffected())
	}
	return inserted, nil
}

func strPtr(value string) *string {
	return &value
}
