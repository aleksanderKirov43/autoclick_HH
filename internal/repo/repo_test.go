package repo

import (
	"testing"
)

func TestNewMockRepo(t *testing.T) {
	repo := NewMockRepo()

	if repo == nil {
		t.Fatal("NewMockRepo вернул nil")
	}

	if repo.responded == nil {
		t.Fatal("Поле responded не инициализировано")
	}
}

func TestMockRepo_AlreadyResponded(t *testing.T) {
	repo := NewMockRepo()

	// Изначально не откликались
	if repo.AlreadyResponded("vacancy-1", "resume-1") {
		t.Error("Ожидалось, что не откликались на вакансию")
	}

	// Сохраняем отклик
	repo.SaveResponse("vacancy-1", "resume-1")

	// Теперь должны быть откликнувшимися
	if !repo.AlreadyResponded("vacancy-1", "resume-1") {
		t.Error("Ожидалось, что уже откликались на вакансию")
	}
}

func TestMockRepo_SaveResponse(t *testing.T) {
	repo := NewMockRepo()

	// Сохраняем отклик
	err := repo.SaveResponse("vacancy-123", "resume-456")
	if err != nil {
		t.Errorf("SaveResponse вернул ошибку: %v", err)
	}

	// Проверяем, что сохранилось
	if !repo.AlreadyResponded("vacancy-123", "resume-456") {
		t.Error("Отклик не был сохранен")
	}
}

func TestMockRepo_UniqueKeys(t *testing.T) {
	repo := NewMockRepo()

	// Разные комбинации vacancy + resume должны быть уникальными
	repo.SaveResponse("vacancy-1", "resume-1")
	repo.SaveResponse("vacancy-1", "resume-2")
	repo.SaveResponse("vacancy-2", "resume-1")

	// Проверяем все комбинации
	if !repo.AlreadyResponded("vacancy-1", "resume-1") {
		t.Error("vacancy-1 + resume-1 не найдено")
	}

	if !repo.AlreadyResponded("vacancy-1", "resume-2") {
		t.Error("vacancy-1 + resume-2 не найдено")
	}

	if !repo.AlreadyResponded("vacancy-2", "resume-1") {
		t.Error("vacancy-2 + resume-1 не найдено")
	}

	// Проверяем, что несуществующая комбинация не найдена
	if repo.AlreadyResponded("vacancy-2", "resume-2") {
		t.Error("vacancy-2 + resume-2 не должно существовать")
	}
}

// Тесты для PostgresRepo (если нужно тестировать с реальной БД)
func TestPostgresRepo_Integration(t *testing.T) {
	// Пропускаем интеграционные тесты, если нет подключения к БД
	t.Skip("Интеграционные тесты с PostgreSQL пропущены")

	// Здесь можно добавить тесты с test container или реальной БД
	// для тестирования PostgresRepo
}
