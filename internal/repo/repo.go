package repo

import (
	"database/sql"
	"log"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) AlreadyResponded(vacancyID, resumeID string) bool {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM responses WHERE vacancy_id = $1 AND resume_id = $2
		)`,
		vacancyID, resumeID,
	).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Ошибка проверки дубликата: %v", err)
		return true
	}
	return exists
}

func (r *PostgresRepo) SaveResponse(vacancyID, resumeID string) error {
	_, err := r.db.Exec(
		`INSERT INTO responses (vacancy_id, resume_id) VALUES ($1, $2)`,
		vacancyID, resumeID,
	)
	if err != nil {
		log.Printf("Ошибка сохранения отклика: %v", err)
	}
	return err
}
