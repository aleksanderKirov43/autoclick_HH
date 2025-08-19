package repo

import (
	"database/sql"
	"time"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CountResponsesToday() (int, error) {
	var count int
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24*time.Hour - time.Nanosecond)

	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM responses WHERE responded_at BETWEEN $1 AND $2",
		start, end,
	).Scan(&count)
	return count, err
}

func (r *PostgresRepo) AlreadyResponded(vacancyID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM responses WHERE vacancy_id=$1)",
		vacancyID,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresRepo) SaveResponse(vacancyID, resumeID string) error {
	_, err := r.db.Exec(
		"INSERT INTO responses(vacancy_id, resume_id) VALUES($1, $2)",
		vacancyID, resumeID,
	)
	return err
}
