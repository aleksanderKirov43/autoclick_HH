package config

import (
	"fmt"
	"os"
)

func SaveTokens(access, refresh string) error {
	f, err := os.OpenFile(".env", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\nACCESS_TOKEN=%s\nREFRESH_TOKEN=%s\n", access, refresh)
	return err
}
