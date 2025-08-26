package hhclient

import "fmt"

type MockClient struct{}

func (m *MockClient) GetToken() (string, error) {
	return "mock-token", nil
}

func (m *MockClient) SearchVacancies(token string, keywords []string) ([]Vacancy, error) {
	// Возвращаем несколько тестовых вакансий
	return []Vacancy{
		{ID: "1", Name: "Golang Developer"},
		{ID: "2", Name: "Go Backend Engineer"},
		{ID: "3", Name: "Frontend Developer"},
	}, nil
}

func (m *MockClient) ApplyVacancy(token, vacancyID, resumeID string) error {
	fmt.Printf("Mock отклик на вакансию %s резюме %s\n", vacancyID, resumeID)
	return nil
}

func (m *MockClient) RefreshToken(refreshToken string) (string, string, error) {
	return "new-access-token", "new-refresh-token", nil
}
