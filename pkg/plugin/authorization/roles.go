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

func FetchInitialRoles(conf *config.Config, roleManager *models.RoleManager) error {
	url := fmt.Sprintf("%s/%s/roles", conf.RbacURL, conf.ApiVersion)

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

	rolesArr := []*models.Role{}
	rolesMap := map[string]*models.Role{}

	err = json.Unmarshal(body, &rolesArr)
	if err != nil {
		return err
	}

	for _, role := range rolesArr {
		rolesMap[role.Name] = role
	}

	roleManager.Lock()
	defer roleManager.Unlock()

	roleManager.Roles = rolesMap

	return nil
}

func UpsertRole(role *models.Role, roleManager *models.RoleManager) error {
	roleManager.Lock()
	defer roleManager.Unlock()

	roleManager.Roles[role.Name] = role

	return nil
}

func DeleteRoleByID(id uint64, roleManager *models.RoleManager) error {
	roleManager.Lock()
	defer roleManager.Unlock()

	for key, role := range roleManager.Roles {
		if role.ID == id {
			delete(roleManager.Roles, key)
			break
		}
	}

	return nil
}
