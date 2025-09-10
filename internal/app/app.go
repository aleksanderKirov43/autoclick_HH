package app

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"autoclick_HH/internal/config"
	"autoclick_HH/internal/hhclient"
	"autoclick_HH/internal/repo"
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

	// Лимит на запуск — передаём явно (200)
	const maxDaily = 200
	autoService := service.NewAutoService(client, pgRepo, maxDaily)

	// Токен: приоритет — access → refresh → password grant
	var token string
	if cfg.AccessToken != "" {
		token = cfg.AccessToken
		if len(token) > 8 {
			log.Printf("Используем выданный access_token, prefix=%q", token[:8])
		} else {
			log.Printf("Используем выданный access_token")
		}
	} else if cfg.RefreshToken != "" {
		log.Printf("Обновляем токен по refresh_token…")
		newAccess, newRefresh, rerr := client.RefreshToken(cfg.RefreshToken)
		if rerr != nil {
			return fmt.Errorf("ошибка обновления токена: %w", rerr)
		}
		token = newAccess
		// Сообщаем префиксы для ручного обновления .env пользователем
		if len(newAccess) > 8 {
			log.Printf("Новый access_token prefix=%q", newAccess[:8])
		}
		if len(newRefresh) > 8 {
			log.Printf("Новый refresh_token prefix=%q (обновите .env при необходимости)", newRefresh[:8])
		}
	} else {
		var gerr error
		token, gerr = client.GetToken()
		if gerr != nil {
			return fmt.Errorf("ошибка получения токена: %w", gerr)
		}
		log.Println("Успешная авторизация по логину/паролю")
	}

	keywords := []string{"go", "golang", "golang developer", "go разработчик"}

	vacancies, err := client.SearchVacancies(token, keywords)
	if err != nil && strings.Contains(err.Error(), "403") {
		log.Println("🔄 Токен истёк, пробуем обновить…")

		time.Sleep(2 * time.Second)

		newAccess, newRefresh, rerr := client.RefreshToken(cfg.RefreshToken)
		if rerr != nil {
			return fmt.Errorf("ошибка повторного обновления токена: %w", rerr)
		}
		token = newAccess
		if err := config.SaveTokens(newAccess, newRefresh); err != nil {
			log.Printf("⚠️ Не удалось сохранить токены: %v", err)
		}
		vacancies, err = client.SearchVacancies(token, keywords)
	}
	if err != nil {
		return fmt.Errorf("ошибка поиска вакансий: %w", err)
	}

	log.Printf("Найдено вакансий: %d", len(vacancies))
	for i := 0; i < len(vacancies) && i < 3; i++ {
		log.Printf("Пример #%d: %s", i+1, vacancies[i].Name)
	}

	if err := autoService.AutoRespond(token, cfg.ResumeID, vacancies); err != nil {
		return fmt.Errorf("ошибка автооткликов: %w", err)
	}

	log.Println("✅ Работа завершена. Перезапусти приложение через 24 часа.")
	return nil

}
