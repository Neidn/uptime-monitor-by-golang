package config

import (
	"encoding/json"
	"os"
	"strings"
)

func GetSecret(key string) string {
	var secrets map[string]string

	secretContext := os.Getenv(SecretContextKey)

	err := json.Unmarshal([]byte(secretContext), &secrets)
	if err != nil {
		return ""
	}

	secret, ok := secrets[key]

	if !ok {
		return os.Getenv(key)
	}

	return secret
}

func GetOwnerRepo() map[string]string {
	repos := GetSecret(GithubRepoKey)
	if repos == "" {
		return nil
	}

	result := strings.Split(repos, "/")
	if len(result) != 2 {
		return nil
	}

	return map[string]string{
		"owner": result[0],
		"repo":  result[1],
	}

}
