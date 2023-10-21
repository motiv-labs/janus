package authorization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/models"
)

func FetchInitialTokens(conf *config.Config, tokenManager *models.TokenManager) error {
	url := fmt.Sprintf("%s/%s/tokens", conf.UserManagementURL, conf.ApiVersion)

	http.DefaultClient.Timeout = 3 * time.Second
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return ErrTimeout
		}
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tokensArr := []*models.JWTToken{}
	tokensMap := map[string]*models.JWTToken{}

	err = json.Unmarshal(body, &tokensArr)
	if err != nil {
		return err
	}

	for _, token := range tokensArr {
		tokensMap[token.Token] = token
	}

	tokenManager.Lock()
	defer tokenManager.Unlock()

	tokenManager.Tokens = tokensMap

	return nil
}

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
