package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadAllowedUsers() (map[int64]bool, error) {
	usersStr := os.Getenv("ALLOWED_USERS")
	if usersStr == "" {
		return nil, errors.New("Nenhum usuário autorizado definido em ALLOWED_USERS")
	}

	users := make(map[int64]bool)
	tokens := strings.Split(usersStr, ",")

	for _, t := range tokens {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("falha de conversão de ID %q: %w", t, err)
		}
		users[id] = true
	}

	if len(users) == 0 {
		return nil, errors.New("Nenhum ID válido encontrado em ALLOWED_USERS")
	}

	return users, nil
}
