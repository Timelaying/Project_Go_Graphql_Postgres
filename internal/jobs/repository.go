package jobs

import (
	"context"
	"strconv"
	"strings"

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

func (r *Repo) UpdateStatus(ctx context.Context, id string, status Status) (*Job, error) {
	var j Job
	err := r.pool.QueryRow(ctx,
		`UPDATE jobs SET status=$1 WHERE id=$2
		 RETURNING id, company, role, link, status, created_at`,
		string(status), id,
	).Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *Repo) UpdateLink(ctx context.Context, id, link string) (*Job, error) {
	var j Job
	err := r.pool.QueryRow(ctx,
		`UPDATE jobs SET link=$1 WHERE id=$2
		 RETURNING id, company, role, link, status, created_at`,
		link, id,
	).Scan(&j.ID, &j.Company, &j.Role, &j.Link, &j.Status, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM jobs WHERE id=$1`, id)
	return err
}

type StatusCount struct {
	Status Status
	Count  int
}

func (r *Repo) StatsByStatus(ctx context.Context) ([]StatusCount, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM jobs GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []StatusCount{}
	for rows.Next() {
		var sc StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

func (r *Repo) SeedDemo(ctx context.Context) (int, error) {
	demos := []struct{ company, role string }{
		{"Google", "Software Engineer"},
		{"Apple", "Backend Developer"},
		{"Meta", "Staff Engineer"},
		{"Amazon", "Site Reliability Engineer"},
		{"Netflix", "Platform Engineer"},
	}

	count := 0
	for _, d := range demos {
		_, err := r.Create(ctx, d.company, d.role, nil)
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
