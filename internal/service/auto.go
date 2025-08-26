package service

import (
	"log"
	"strings"

	"autoclick_HH/internal/hhclient"
)

type AutoService struct {
	client   hhclient.HHClient
	repo     hhclient.Repo
	maxDaily int
}

func NewAutoService(client hhclient.HHClient, repo hhclient.Repo, maxDaily int) *AutoService {
	return &AutoService{
		client:   client,
		repo:     repo,
		maxDaily: maxDaily,
	}
}

func (s *AutoService) AutoRespond(token, resumeID string, vacancies []hhclient.Vacancy) error {
	count := 0

	// Список исключений 
	excluded := []string{"senior", "администратор", "курьер", "кассир", "автор", "менеджер", "руководитель", "бариста", "архитектор"}

	for _, v := range vacancies {
		nameLower := strings.ToLower(v.Name)

		// откликаемся только на Go/Golang
		if !strings.Contains(nameLower, "go") && !strings.Contains(nameLower, "golang") {
			continue
		}

		// проверка по исключениям
		skip := false
		for _, word := range excluded {
			if strings.Contains(nameLower, word) {
				log.Printf("🚫 Пропускаем позицию (%s): %s", word, v.Name)
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// проверка дубля
		if s.repo.AlreadyResponded(v.ID, resumeID) {
			log.Printf("⚠️ Уже откликались на вакансию %s, пропускаем", v.ID)
			continue
		}

		// проверка лимита
		if count >= s.maxDaily {
			log.Printf("⛔ Достигнут лимит (%d откликов)", s.maxDaily)
			break
		}

		// откликаемся
		if err := s.client.ApplyVacancy(token, v.ID, resumeID); err != nil {
			log.Printf("Ошибка отклика на %s: %v", v.ID, err)
			continue
		}

		// сохраняем факт отклика
		if err := s.repo.SaveResponse(v.ID, resumeID); err != nil {
			log.Printf("Ошибка сохранения отклика %s: %v", v.ID, err)
			continue
		}

		log.Printf("✅ Откликнулись на вакансию #: %s, Должность: %s", v.ID, v.Name)
		count++
	}

	return nil
}

