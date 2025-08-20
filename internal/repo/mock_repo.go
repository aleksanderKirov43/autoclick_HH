package repo

import (
	"fmt"
)

type MockRepo struct {
	responded map[string]bool
}

func NewMockRepo() *MockRepo {
	return &MockRepo{responded: make(map[string]bool)}
}

func (m *MockRepo) AlreadyResponded(vacancyID, resumeID string) bool {
	key := vacancyID + "_" + resumeID
	return m.responded[key]
}

func (m *MockRepo) SaveResponse(vacancyID, resumeID string) error {
	key := vacancyID + "_" + resumeID
	m.responded[key] = true
	fmt.Printf("Mock сохранение отклика: %s\n", key)
	return nil
}
