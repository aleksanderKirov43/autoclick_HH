package main

import (
	"log"

	"autoclick_HH/internal/hhclient"
	"autoclick_HH/internal/repo"
	"autoclick_HH/internal/service"
)

func main() {
	// создаём мок клиента и репозитория
	mockClient := &hhclient.MockClient{}
	mockRepo := repo.NewMockRepo()

	// сервис с лимитом 2 для теста
	autoService := service.NewAutoService(mockClient, mockRepo, 100)

	// получаем "токен" и вакансии
	token, err := mockClient.GetToken()
	if err != nil {
		log.Fatalf("Ошибка получения мок-токена: %v", err)
	}

	keywords := []string{"go", "golang"}

	vacancies, err := mockClient.SearchVacancies(token, keywords)
	if err != nil {
		log.Fatalf("Ошибка поиска вакансий в моках: %v", err)
	}

	// запускаем автоотклик
	err = autoService.AutoRespond(token, "mock-resume-id", vacancies)
	if err != nil {
		log.Fatal(err)
	}
}
