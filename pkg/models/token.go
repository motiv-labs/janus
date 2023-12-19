package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hellofresh/janus/pkg/config"
)

var (
	ErrTimeout = errors.New("timed out")
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

func (tm *TokenManager) FetchTokens() error {
	url := fmt.Sprintf("%s/%s/tokens", tm.Conf.UserManagementURL, tm.Conf.ApiVersion)

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

	tokensArr := []*JWTToken{}
	tokensMap := map[string]*JWTToken{}

	err = json.Unmarshal(body, &tokensArr)
	if err != nil {
		return err
	}

	for _, token := range tokensArr {
		tokensMap[token.Token] = token
	}

	tm.Lock()
	defer tm.Unlock()

	tm.Tokens = tokensMap

	return nil
}
