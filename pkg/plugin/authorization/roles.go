package authorization

import (
	"github.com/hellofresh/janus/pkg/models"
)

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
