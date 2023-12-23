package authorization

import (
	"github.com/hellofresh/janus/pkg/models"
)

func UpsertToken(token *models.JWTToken, tokenManager *models.TokenManager) error {
	tokenManager.Lock()
	defer tokenManager.Unlock()

	tokenManager.Tokens[token.Token] = token
	return nil
}

func DeleteTokenByID(id uint64, tokenManager *models.TokenManager) error {
	tokenManager.Lock()
	defer tokenManager.Unlock()

	for key, token := range tokenManager.Tokens {
		if token.ID == id {
			delete(tokenManager.Tokens, key)
			break
		}
	}

	return nil
}
