package hhclient

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := New("client-id", "client-secret", "username", "password")

	if client == nil {
		t.Fatal("New вернул nil")
	}

	if client.clientID != "client-id" {
		t.Errorf("Ожидался clientID = 'client-id', получен '%s'", client.clientID)
	}

	if client.clientSecret != "client-secret" {
		t.Errorf("Ожидался clientSecret = 'client-secret', получен '%s'", client.clientSecret)
	}

	if client.username != "username" {
		t.Errorf("Ожидался username = 'username', получен '%s'", client.username)
	}

	if client.password != "password" {
		t.Errorf("Ожидался password = 'password', получен '%s'", client.password)
	}
}

func TestVacancyStruct(t *testing.T) {
	vacancy := Vacancy{
		ID:   "123",
		Name: "Golang Developer",
		Area: struct {
			Name string `json:"name"`
		}{
			Name: "Москва",
		},
		Employer: struct {
			Name string `json:"name"`
		}{
			Name: "Tech Company",
		},
	}

	if vacancy.ID != "123" {
		t.Errorf("Ожидался ID = '123', получен '%s'", vacancy.ID)
	}

	if vacancy.Name != "Golang Developer" {
		t.Errorf("Ожидался Name = 'Golang Developer', получен '%s'", vacancy.Name)
	}

	if vacancy.Area.Name != "Москва" {
		t.Errorf("Ожидался Area.Name = 'Москва', получен '%s'", vacancy.Area.Name)
	}

	if vacancy.Employer.Name != "Tech Company" {
		t.Errorf("Ожидался Employer.Name = 'Tech Company', получен '%s'", vacancy.Employer.Name)
	}
}

func TestMockClient_GetToken(t *testing.T) {
	mock := &MockClient{}

	token, err := mock.GetToken()
	if err != nil {
		t.Errorf("MockClient.GetToken вернул ошибку: %v", err)
	}

	if token != "mock-token" {
		t.Errorf("Ожидался токен 'mock-token', получен '%s'", token)
	}
}

func TestMockClient_SearchVacancies(t *testing.T) {
	mock := &MockClient{}

	vacancies, err := mock.SearchVacancies("token", []string{"go", "golang"})
	if err != nil {
		t.Errorf("MockClient.SearchVacancies вернул ошибку: %v", err)
	}

	if len(vacancies) != 3 {
		t.Errorf("Ожидалось 3 вакансии, получено %d", len(vacancies))
	}

	// Проверяем первую вакансию
	if vacancies[0].Name != "Golang Developer" {
		t.Errorf("Ожидалась вакансия 'Golang Developer', получена '%s'", vacancies[0].Name)
	}
}

func TestMockClient_ApplyVacancy(t *testing.T) {
	mock := &MockClient{}

	err := mock.ApplyVacancy("token", "vacancy-123", "resume-456")
	if err != nil {
		t.Errorf("MockClient.ApplyVacancy вернул ошибку: %v", err)
	}
}
