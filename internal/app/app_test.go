package app

import (
	"fmt"
	"testing"

	"autoclick_HH/internal/hhclient"
	"autoclick_HH/internal/repo"
	"autoclick_HH/internal/service"
)

// Mock компоненты для тестирования
type mockAppDeps struct {
	client hhclient.HHClient
	repo   hhclient.Repo
}

func newMockAppDeps() *mockAppDeps {
	return &mockAppDeps{
		client: &hhclient.MockClient{},
		repo:   repo.NewMockRepo(),
	}
}

func TestAppIntegration(t *testing.T) {
	deps := newMockAppDeps()

	// Тестируем создание сервиса
	service := service.NewAutoService(deps.client, deps.repo, 10)
	if service == nil {
		t.Fatal("Не удалось создать AutoService")
	}

	// Тестируем получение токена
	token, err := deps.client.GetToken()
	if err != nil {
		t.Fatalf("Ошибка получения токена: %v", err)
	}

	if token == "" {
		t.Error("Получен пустой токен")
	}

	// Тестируем поиск вакансий
	keywords := []string{"go", "golang"}
	vacancies, err := deps.client.SearchVacancies(token, keywords)
	if err != nil {
		t.Fatalf("Ошибка поиска вакансий: %v", err)
	}

	if len(vacancies) == 0 {
		t.Error("Не найдено вакансий")
	}

	// Тестируем автоотклик
	resumeID := "test-resume-id"
	err = service.AutoRespond(token, resumeID, vacancies)
	if err != nil {
		t.Fatalf("Ошибка автоотклика: %v", err)
	}
}

func TestAppWithErrors(t *testing.T) {
	// Создаем мок клиент, который возвращает ошибки
	errorClient := &errorMockClient{}
	repo := repo.NewMockRepo()

	service := service.NewAutoService(errorClient, repo, 5)

	// Тестируем обработку ошибок
	_, err := errorClient.GetToken()
	if err == nil {
		t.Error("Ожидалась ошибка получения токена")
	}

	// Даже с ошибкой токена, сервис должен создаться
	if service == nil {
		t.Fatal("Не удалось создать AutoService с error клиентом")
	}
}

// Mock клиент, который всегда возвращает ошибки
type errorMockClient struct{}

func (m *errorMockClient) GetToken() (string, error) {
	return "", &mockError{message: "token error"}
}

func (m *errorMockClient) SearchVacancies(token string, keywords []string) ([]hhclient.Vacancy, error) {
	return nil, &mockError{message: "search error"}
}

func (m *errorMockClient) ApplyVacancy(token, vacancyID, resumeID string) error {
	return &mockError{message: "apply error"}
}

// Простая структура ошибки для тестирования
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

// Тест производительности
func BenchmarkAutoService(b *testing.B) {
	deps := newMockAppDeps()
	service := service.NewAutoService(deps.client, deps.repo, 100)

	// Создаем много вакансий для тестирования
	vacancies := make([]hhclient.Vacancy, 1000)
	for i := 0; i < 1000; i++ {
		vacancies[i] = hhclient.Vacancy{
			ID:   fmt.Sprintf("vacancy-%d", i),
			Name: "Golang Developer",
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		service.AutoRespond("test-token", "test-resume", vacancies)
	}
}
