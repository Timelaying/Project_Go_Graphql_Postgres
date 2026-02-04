package jobs

import (
	"context"
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
		where = append(where, `status = $`+itoa(i))
		args = append(args, string(*status))
		i++
	}
	if q != nil && strings.TrimSpace(*q) != "" {
		where = append(where, `(company ILIKE $`+itoa(i)+` OR role ILIKE $`+itoa(i)+`)`)
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

func itoa(i int) string {
	// tiny helper (avoid strconv import if you want; but strconv.Itoa is fine too)
	const digits = "0123456789"
	if i < 10 {
		return string(digits[i])
	}
	// for weekend scope, you can just use strconv.Itoa instead
	return "0"
}
