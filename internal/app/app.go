package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"autoclick_HH/internal/config"
	"autoclick_HH/internal/hhclient"
	"autoclick_HH/internal/repository"
	"autoclick_HH/internal/service"
)

func Run() error {
	// Загружаем конфиг
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Создаём клиент для hh.ru
	client := hhclient.New(cfg.ClientID, cfg.ClientSecret, cfg.Username, cfg.Password)

	// Получаем токен
	token, err := client.GetToken()
	if err != nil {
		return fmt.Errorf("ошибка получения токена: %w", err)
	}

	log.Println("Успешная авторизация. Access token:", token[:15], "...")

	// Ищем вакансии
	vacancies, err := client.SearchVacancies(token, "golang")
	if err != nil {
		return fmt.Errorf("ошибка поиска вакансий: %w", err)
	}

	fmt.Printf("Найдено вакансий: %d\n", len(vacancies))

	// --- подключение к PostgreSQL ---
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	postgresRepo := repo.NewPostgresRepo(db)
	autoService := service.NewAutoService(client, postgresRepo)

	// пример автоотклика
	resumeID := cfg.ResumeID
	if err := autoService.AutoRespond(token, resumeID, vacancies); err != nil {
		return fmt.Errorf("ошибка автоотклика: %w", err)
	}

	return nil
}
