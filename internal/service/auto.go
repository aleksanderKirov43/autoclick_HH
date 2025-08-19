package service

import (
	"fmt"
	"time"

	"autoclick_HH/internal/hhclient"
	"autoclick_HH/internal/repository"
)

type AutoService struct {
	client *hhclient.Client
	repo   *repo.PostgresRepo
}

const MaxDailyResponses = 20

func NewAutoService(client *hhclient.Client, repo *repo.PostgresRepo) *AutoService {
	return &AutoService{
		client: client,
		repo:   repo,
	}
}

func (s *AutoService) AutoRespond(token, resumeID string, vacancies []hhclient.Vacancy) error {
	todayCount, err := s.repo.CountResponsesToday()
	if err != nil {
		return fmt.Errorf("ошибка подсчета откликов: %w", err)
	}

	for _, v := range vacancies {
		if todayCount >= MaxDailyResponses {
			fmt.Println("Достигнут лимит откликов за сегодня")
			break
		}

		already, err := s.repo.AlreadyResponded(v.ID)
		if err != nil {
			return fmt.Errorf("ошибка проверки вакансии: %w", err)
		}

		if !already {
			if err := s.client.ApplyVacancy(token, v.ID, resumeID); err != nil {
				return fmt.Errorf("ошибка отклика: %w", err)
			}

			if err := s.repo.SaveResponse(v.ID, resumeID); err != nil {
				return fmt.Errorf("ошибка сохранения отклика: %w", err)
			}

			todayCount++
			fmt.Printf("Отклик отправлен на вакансию [%s] %s\n", v.ID, v.Name)
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}
