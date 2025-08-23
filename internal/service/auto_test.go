package service

import (
	"errors"
	"testing"

	"autoclick_HH/internal/hhclient"
)

// Mock клиент для тестирования
type mockClient struct {
	tokenError      error
	searchError     error
	applyError      error
	searchVacancies []hhclient.Vacancy
}

func (m *mockClient) GetToken() (string, error) {
	if m.tokenError != nil {
		return "", m.tokenError
	}
	return "test-token", nil
}

func (m *mockClient) SearchVacancies(token string, keywords []string) ([]hhclient.Vacancy, error) {
	if m.searchError != nil {
		return nil, m.searchError
	}
	return m.searchVacancies, nil
}

func (m *mockClient) ApplyVacancy(token, vacancyID, resumeID string) error {
	return m.applyError
}

// Mock репозиторий для тестирования
type mockRepo struct {
	alreadyResponded map[string]bool
	saveError        error
}

func (m *mockRepo) AlreadyResponded(vacancyID, resumeID string) bool {
	key := vacancyID + "_" + resumeID
	return m.alreadyResponded[key]
}

func (m *mockRepo) SaveResponse(vacancyID, resumeID string) error {
	if m.saveError != nil {
		return m.saveError
	}
	key := vacancyID + "_" + resumeID
	m.alreadyResponded[key] = true
	return nil
}

func TestNewAutoService(t *testing.T) {
	client := &mockClient{}
	repo := &mockRepo{alreadyResponded: make(map[string]bool)}

	service := NewAutoService(client, repo, 100)

	if service == nil {
		t.Fatal("NewAutoService вернул nil")
	}

	if service.maxDaily != 100 {
		t.Errorf("Ожидался maxDaily = 100, получен %d", service.maxDaily)
	}
}

func TestAutoRespond_Success(t *testing.T) {
	vacancies := []hhclient.Vacancy{
		{ID: "1", Name: "Golang Developer"},
		{ID: "2", Name: "Go Backend Engineer"},
		{ID: "3", Name: "Frontend Developer"},
	}

	client := &mockClient{
		searchVacancies: vacancies,
	}

	repo := &mockRepo{alreadyResponded: make(map[string]bool)}

	service := NewAutoService(client, repo, 5)

	err := service.AutoRespond("test-token", "resume-1", vacancies)
	if err != nil {
		t.Errorf("AutoRespond вернул ошибку: %v", err)
	}
}

func TestAutoRespond_AlreadyResponded(t *testing.T) {
	vacancies := []hhclient.Vacancy{
		{ID: "1", Name: "Golang Developer"},
		{ID: "2", Name: "Go Backend Engineer"},
	}

	client := &mockClient{
		searchVacancies: vacancies,
	}

	repo := &mockRepo{
		alreadyResponded: map[string]bool{
			"1_resume-1": true, // уже откликались
		},
	}

	service := NewAutoService(client, repo, 5)

	err := service.AutoRespond("test-token", "resume-1", vacancies)
	if err != nil {
		t.Errorf("AutoRespond вернул ошибку: %v", err)
	}
}

func TestAutoRespond_MaxDailyLimit(t *testing.T) {
	vacancies := []hhclient.Vacancy{
		{ID: "1", Name: "Golang Developer"},
		{ID: "2", Name: "Go Backend Engineer"},
		{ID: "3", Name: "Go Developer"},
		{ID: "4", Name: "Backend Go Engineer"},
	}

	client := &mockClient{
		searchVacancies: vacancies,
	}

	repo := &mockRepo{alreadyResponded: make(map[string]bool)}

	// Лимит 2 отклика в день
	service := NewAutoService(client, repo, 2)

	err := service.AutoRespond("test-token", "resume-1", vacancies)
	if err != nil {
		t.Errorf("AutoRespond вернул ошибку: %v", err)
	}
}

func TestAutoRespond_ApplyError(t *testing.T) {
	vacancies := []hhclient.Vacancy{
		{ID: "1", Name: "Golang Developer"},
	}

	client := &mockClient{
		searchVacancies: vacancies,
		applyError:      errors.New("API error"),
	}

	repo := &mockRepo{alreadyResponded: make(map[string]bool)}

	service := NewAutoService(client, repo, 5)

	err := service.AutoRespond("test-token", "resume-1", vacancies)
	if err != nil {
		t.Errorf("AutoRespond вернул ошибку: %v", err)
	}
}

func TestAutoRespond_SaveError(t *testing.T) {
	vacancies := []hhclient.Vacancy{
		{ID: "1", Name: "Golang Developer"},
	}

	client := &mockClient{
		searchVacancies: vacancies,
	}

	repo := &mockRepo{
		alreadyResponded: make(map[string]bool),
		saveError:        errors.New("DB error"),
	}

	service := NewAutoService(client, repo, 5)

	err := service.AutoRespond("test-token", "resume-1", vacancies)
	if err != nil {
		t.Errorf("AutoRespond вернул ошибку: %v", err)
	}
}

func TestAutoRespond_NonGoVacancies(t *testing.T) {
	vacancies := []hhclient.Vacancy{
		{ID: "1", Name: "Frontend Developer"},
		{ID: "2", Name: "Python Developer"},
		{ID: "3", Name: "Java Developer"},
	}

	client := &mockClient{
		searchVacancies: vacancies,
	}

	repo := &mockRepo{alreadyResponded: make(map[string]bool)}

	service := NewAutoService(client, repo, 5)

	err := service.AutoRespond("test-token", "resume-1", vacancies)
	if err != nil {
		t.Errorf("AutoRespond вернул ошибку: %v", err)
	}
}
