package service

import (
	"log"

	"autoclick_HH/internal/hhclient"
	repo "autoclick_HH/internal/repository"
)

type AutoService struct {
	client   *hhclient.Client
	repo     *repo.PostgresRepo
	maxDaily int
}

func NewAutoService(client *hhclient.Client, repo *repo.PostgresRepo, maxDaily int) *AutoService {
	return &AutoService{
		client:   client,
		repo:     repo,
		maxDaily: maxDaily,
	}
}

func (s *AutoService) AutoRespond(token, resumeID string, vacancies []hhclient.Vacancy) error {
	count := 0

	for _, v := range vacancies {
		if s.repo.AlreadyResponded(v.ID, resumeID) {
			log.Printf("⚠️ Уже откликались на вакансию %s, пропускаем", v.ID)
			continue
		}

		if count >= s.maxDaily {
			log.Printf("⛔ Достигнут лимит (%d откликов)", s.maxDaily)
			break
		}

		if err := s.client.ApplyVacancy(token, v.ID, resumeID); err != nil {
			log.Printf("Ошибка отклика на %s: %v", v.ID, err)
			continue
		}

		if err := s.repo.SaveResponse(v.ID, resumeID); err != nil {
			log.Printf("Ошибка сохранения отклика %s: %v", v.ID, err)
			continue
		}

		log.Printf("✅ Откликнулись на вакансию %s", v.ID)
		count++
	}

	return nil
}
