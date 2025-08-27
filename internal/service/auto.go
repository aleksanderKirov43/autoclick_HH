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

	// Список исключений по должностям
	excluded := []string{"senior", "администратор", "курьер", "кассир", "автор", "менеджер", "руководитель", "бариста", "архитектор", "дизайнер"}

	// Список компаний, на которые не нужно откликаться
	excludedCompanies := []string{"ozon", "wildberries", "т-банк", "мтс", "магнит", "суточно.ру"}

	for _, v := range vacancies {
		nameLower := strings.ToLower(v.Name)

		// Проверка на Go/Golang
		if !strings.Contains(nameLower, "go") && !strings.Contains(nameLower, "golang") {
			continue
		}

		// Проверка по исключениям в названии вакансии
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

		// Проверка по исключениям компаний
		// Для этого нужно получить название компании. Предполагаем, что оно есть в v.Employer.Name
		companyName := ""
		if v.Employer.Name != "" {
			companyName = strings.ToLower(v.Employer.Name)
		}
		for _, exclCompany := range excludedCompanies {
			if strings.Contains(companyName, exclCompany) {
				log.Printf("🚫 Пропускаем вакансию от компании (%s): %s", exclCompany, v.Employer.Name)
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Проверка дубля
		if s.repo.AlreadyResponded(v.ID, resumeID) {
			log.Printf("⚠️ Уже откликались на вакансию %s, пропускаем", v.ID)
			continue
		}

		// Проверка лимита
		if count >= s.maxDaily {
			log.Printf("⛔ Достигнут лимит (%d откликов)", s.maxDaily)
			break
		}

		// Откликаемся
		if err := s.client.ApplyVacancy(token, v.ID, resumeID); err != nil {
			log.Printf("Ошибка отклика на %s: %v", v.ID, err)
			continue
		}

		// Сохраняем факт отклика
		if err := s.repo.SaveResponse(v.ID, resumeID); err != nil {
			log.Printf("Ошибка сохранения отклика %s: %v", v.ID, err)
			continue
		}

		log.Printf("✅ Откликнулись на вакансию #: %s, Должность: %s", v.ID, v.Name)
		count++
	}

	return nil
}
