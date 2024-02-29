package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func GetSecret(key string) string {
	var secrets map[string]string

	secretContext := os.Getenv(SecretContextKey)

	err := json.Unmarshal([]byte(secretContext), &secrets)

	if err != nil {
		return os.Getenv(key)
	}

	secret, ok := secrets[key]

	log.Println("ok", ok)

	if !ok {
		return os.Getenv(key)
	}

	return secret
}

func GetOwnerRepo() (string, string) {
	repos := GetSecret(GithubRepoKey)
	if repos == "" {
		return "", ""
	}

	result := strings.Split(repos, "/")
	if len(result) != 2 {
		return "", ""
	}

	return result[0], result[1]

}

func GetToken() string {
	_TokenKey := []string{"GITHUB_TOKEN", "token", "GH_PAT"}
	var token string

	for _, key := range _TokenKey {
		token = GetSecret(key)
		if token != "" {
			log.Printf("Token found: %s, %s", key, token)
			break
		}
	}

	return token
}
