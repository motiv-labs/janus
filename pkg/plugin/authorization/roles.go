package authorization

import (
	"encoding/json"
	"errors"
	"fmt"
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

func NewRoleManager(conf *config.Config) *RoleManager {
	return &RoleManager{
		Roles: map[string]*Role{},
		Conf:  conf,
	}
}

func (rm *RoleManager) FetchRoles() error {
	url := fmt.Sprintf("%s/%s/roles", rm.Conf.RbacURL, rm.Conf.ApiVersion)

	body, err := doGetRequestWithTimeout(url, 3*time.Second)
	if err != nil {
		if errors.Is(err, ErrTimeout) {
			return nil
		}
		return err
	}

	rolesSlice := []*Role{}
	err = json.Unmarshal(body, &rolesSlice)
	if err != nil {
		return err
	}

	rm.Lock()
	defer rm.Unlock()

	rm.Roles = rolesSliceToMap(rolesSlice)

	return nil
}

func rolesSliceToMap(rolesSlice []*Role) map[string]*Role {
	rolesMap := map[string]*Role{}

	for _, role := range rolesSlice {
		rolesMap[role.Name] = role
	}

	return rolesMap
}

func (rm *RoleManager) UpsertRoles(roles []*Role) {
	rm.Lock()
	defer rm.Unlock()

	for _, role := range roles {
		rm.Roles[role.Name] = role
	}
}

func (rm *RoleManager) DeleteRolesByIDs(ids []uint64) {
	rm.Lock()
	defer rm.Unlock()

	for _, id := range ids {
		for key, role := range rm.Roles {
			if role.ID == id {
				delete(rm.Roles, key)
			}
		}
	}
}
