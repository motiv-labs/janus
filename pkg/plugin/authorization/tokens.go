package authorization

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hellofresh/janus/pkg/config"
)

type JWTToken struct {
	ID             uint64    `json:"id"`
	UserID         uint64    `json:"user_id"`
	Token          string    `json:"token"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type TokenManager struct {
	Tokens map[string]*JWTToken
	Conf   *config.Config
	sync.RWMutex
}

func NewTokenManager(conf *config.Config) *TokenManager {
	return &TokenManager{
		Tokens: map[string]*JWTToken{},
		Conf:   conf,
	}
}

func (tm *TokenManager) FetchTokens() error {
	url := fmt.Sprintf("%s/%s/tokens", tm.Conf.UserManagementURL, tm.Conf.ApiVersion)
	body, err := doGetRequestWithTimeout(url, 3*time.Second)
	if err != nil {
		if errors.Is(err, ErrTimeout) {
			return nil
		}
		return err
	}

	tokensSlice := []*JWTToken{}
	err = json.Unmarshal(body, &tokensSlice)
	if err != nil {
		return err
	}

	tm.Lock()
	defer tm.Unlock()

	tm.Tokens = tokenSliceToMap(tokensSlice)

	return nil
}

func tokenSliceToMap(tokensSlice []*JWTToken) map[string]*JWTToken {
	tokensMap := map[string]*JWTToken{}

	for _, token := range tokensSlice {
		tokensMap[token.Token] = token
	}

	return tokensMap
}

func (tm *TokenManager) UpsertTokens(tokens []*JWTToken) {
	tm.Lock()
	defer tm.Unlock()

	for _, token := range tokens {
		tm.Tokens[token.Token] = token
	}
}

func (tm *TokenManager) DeleteTokensByIDs(ids []uint64) {
	tm.Lock()
	defer tm.Unlock()

	for _, id := range ids {
		for key, token := range tm.Tokens {
			if token.ID == id {
				delete(tm.Tokens, key)
			}
		}
	}
}
