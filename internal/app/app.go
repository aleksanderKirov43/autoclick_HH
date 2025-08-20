package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"autoclick_HH/internal/config"
	"autoclick_HH/internal/hhclient"
	repo "autoclick_HH/internal/repository"
	"autoclick_HH/internal/service"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// --- подключение к БД ---
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("postgres недоступен: %w", err)
	}

	pgRepo := repo.NewPostgresRepo(db)
	client := hhclient.New(cfg.ClientID, cfg.ClientSecret, cfg.Username, cfg.Password)

	// Лимит на запуск — передаём явно (100)
	const maxDaily = 100
	autoService := service.NewAutoService(client, pgRepo, maxDaily)

	token, err := client.GetToken()
	if err != nil {
		return fmt.Errorf("ошибка получения токена: %w", err)
	}
	log.Println("Успешная авторизация")

	// основной фильтр
	keywords := []string{"go", "golang"}
	vacancies, err := client.SearchVacancies(token, keywords)
	if err != nil {
		return fmt.Errorf("ошибка поиска вакансий: %w", err)
	}

	if err := autoService.AutoRespond(token, cfg.ResumeID, vacancies); err != nil {
		return fmt.Errorf("ошибка автооткликов: %w", err)
	}

	log.Println("✅ Работа завершена. Перезапусти приложение через 24 часа.")
	return nil
}
