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

type Role struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Features []Feature
}

type Feature struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
}

type RoleManager struct {
	Roles map[string]*Role
	Conf  *config.Config
	sync.Mutex
}

func (rm *RoleManager) FetchRoles() error {
	url := fmt.Sprintf("%s/%s/roles", rm.Conf.RbacURL, rm.Conf.ApiVersion)

	http.DefaultClient.Timeout = 5 * time.Second
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

	rolesArr := []*Role{}
	rolesMap := map[string]*Role{}

	err = json.Unmarshal(body, &rolesArr)
	if err != nil {
		return err
	}

	for _, role := range rolesArr {
		rolesMap[role.Name] = role
	}

	rm.Lock()
	defer rm.Unlock()

	rm.Roles = rolesMap

	return nil
}
