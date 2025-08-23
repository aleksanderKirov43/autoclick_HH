.PHONY: test test-verbose test-coverage build run clean help

# Переменные
BINARY_NAME=autoclick_HH
MAIN_PATH=./cmd

# Цвета для вывода
GREEN=\033[32m
RED=\033[31m
YELLOW=\033[33m
NC=\033[0m # No Color

# Помощь
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "  $(YELLOW)test$(NC)          - Запустить все тесты"
	@echo "  $(YELLOW)test-verbose$(NC)  - Запустить тесты с подробным выводом"
	@echo "  $(YELLOW)test-coverage$(NC) - Запустить тесты с покрытием кода"
	@echo "  $(YELLOW)build$(NC)         - Собрать приложение"
	@echo "  $(YELLOW)run$(NC)           - Запустить приложение"
	@echo "  $(YELLOW)clean$(NC)         - Очистить собранные файлы"
	@echo "  $(YELLOW)help$(NC)          - Показать эту справку"

# Тесты
test:
	@echo "$(GREEN)Запуск тестов...$(NC)"
	go test ./...

test-verbose:
	@echo "$(GREEN)Запуск тестов с подробным выводом...$(NC)"
	go test -v ./...

test-coverage:
	@echo "$(GREEN)Запуск тестов с покрытием кода...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии сохранен в coverage.html$(NC)"

# Сборка
build:
	@echo "$(GREEN)Сборка приложения...$(NC)"
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Приложение собрано: $(BINARY_NAME)$(NC)"

# Запуск
run:
	@echo "$(GREEN)Запуск приложения...$(NC)"
	go run $(MAIN_PATH)

# Очистка
clean:
	@echo "$(GREEN)Очистка...$(NC)"
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html
	@echo "$(GREEN)Очистка завершена$(NC)"

# Тесты отдельных компонентов
test-service:
	@echo "$(GREEN)Тестирование сервиса...$(NC)"
	go test -v ./internal/service

test-client:
	@echo "$(GREEN)Тестирование HH клиента...$(NC)"
	go test -v ./internal/hhclient

test-repo:
	@echo "$(GREEN)Тестирование репозитория...$(NC)"
	go test -v ./internal/repo

test-config:
	@echo "$(GREEN)Тестирование конфигурации...$(NC)"
	go test -v ./internal/config

test-app:
	@echo "$(GREEN)Тестирование приложения...$(NC)"
	go test -v ./internal/app

# Проверка качества кода
lint:
	@echo "$(GREEN)Проверка качества кода...$(NC)"
	golangci-lint run

# Установка зависимостей
deps:
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	go mod download
	go mod tidy

# Полная проверка проекта
check: deps lint test
	@echo "$(GREEN)Все проверки пройдены успешно!$(NC)"

# Информация о ветках
info:
	@echo "$(GREEN)Информация о проекте:$(NC)"
	@echo "$(YELLOW)Основная ветка: dev$(NC)"
	@echo "$(YELLOW)CI/CD: настроен для ветки dev$(NC)"
	@echo "$(YELLOW)Тесты: make test$(NC)"
	@echo "$(YELLOW)Сборка: make build$(NC)"
