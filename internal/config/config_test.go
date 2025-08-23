package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalEnv := make(map[string]string)
	envVars := []string{
		"HH_CLIENT_ID", "HH_CLIENT_SECRET", "HH_USERNAME", "HH_PASSWORD",
		"HH_RESUME_ID", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			originalEnv[envVar] = value
		}
	}

	// Устанавливаем тестовые значения
	testEnv := map[string]string{
		"HH_CLIENT_ID":     "test-client-id",
		"HH_CLIENT_SECRET": "test-client-secret",
		"HH_USERNAME":      "test-username",
		"HH_PASSWORD":      "test-password",
		"HH_RESUME_ID":     "test-resume-id",
		"DB_HOST":          "localhost",
		"DB_PORT":          "5432",
		"DB_USER":          "testuser",
		"DB_PASSWORD":      "testpass",
		"DB_NAME":          "testdb",
	}

	for key, value := range testEnv {
		os.Setenv(key, value)
	}

	// Восстанавливаем оригинальные значения после теста
	defer func() {
		for key := range testEnv {
			os.Unsetenv(key)
		}
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	// Загружаем конфигурацию
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load вернул ошибку: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load вернул nil конфигурацию")
	}

	// Проверяем, что все значения загружены корректно
	if cfg.ClientID != "test-client-id" {
		t.Errorf("Ожидался ClientID = 'test-client-id', получен '%s'", cfg.ClientID)
	}

	if cfg.ClientSecret != "test-client-secret" {
		t.Errorf("Ожидался ClientSecret = 'test-client-secret', получен '%s'", cfg.ClientSecret)
	}

	if cfg.Username != "test-username" {
		t.Errorf("Ожидался Username = 'test-username', получен '%s'", cfg.Username)
	}

	if cfg.Password != "test-password" {
		t.Errorf("Ожидался Password = 'test-password', получен '%s'", cfg.Password)
	}

	if cfg.ResumeID != "test-resume-id" {
		t.Errorf("Ожидался ResumeID = 'test-resume-id', получен '%s'", cfg.ResumeID)
	}

	if cfg.DBHost != "localhost" {
		t.Errorf("Ожидался DBHost = 'localhost', получен '%s'", cfg.DBHost)
	}

	if cfg.DBPort != "5432" {
		t.Errorf("Ожидался DBPort = '5432', получен '%s'", cfg.DBPort)
	}

	if cfg.DBUser != "testuser" {
		t.Errorf("Ожидался DBUser = 'testuser', получен '%s'", cfg.DBUser)
	}

	if cfg.DBPassword != "testpass" {
		t.Errorf("Ожидался DBPassword = 'testpass', получен '%s'", cfg.DBPassword)
	}

	if cfg.DBName != "testdb" {
		t.Errorf("Ожидался DBName = 'testdb', получен '%s'", cfg.DBName)
	}
}

func TestLoad_EmptyEnv(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalEnv := make(map[string]string)
	envVars := []string{
		"HH_CLIENT_ID", "HH_CLIENT_SECRET", "HH_USERNAME", "HH_PASSWORD",
		"HH_RESUME_ID", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			originalEnv[envVar] = value
		}
		os.Unsetenv(envVar)
	}

	// Восстанавливаем оригинальные значения после теста
	defer func() {
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	// Загружаем конфигурацию с пустыми переменными
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load вернул ошибку: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load вернул nil конфигурацию")
	}

	// Проверяем, что все значения пустые
	if cfg.ClientID != "" {
		t.Errorf("Ожидался пустой ClientID, получен '%s'", cfg.ClientID)
	}

	if cfg.DBHost != "" {
		t.Errorf("Ожидался пустой DBHost, получен '%s'", cfg.DBHost)
	}
}
