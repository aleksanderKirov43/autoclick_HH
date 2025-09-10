package config

import (
	"bytes"
	"log"
	"os"
	"strings"
)

func SaveTokens(access, refresh string) error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	var buf bytes.Buffer
	for _, line := range lines {
		if strings.HasPrefix(line, "HH_ACCESS_TOKEN=") || strings.HasPrefix(line, "HH_REFRESH_TOKEN=") {
			continue // удаляем старые токены
		}
		buf.WriteString(line + "\n")
	}
	buf.WriteString("HH_ACCESS_TOKEN=" + access + "\n")
	buf.WriteString("HH_REFRESH_TOKEN=" + refresh + "\n")
	log.Println("✅ Токены успешно обновлены в .env")
	return os.WriteFile(".env", buf.Bytes(), 0644)
}
