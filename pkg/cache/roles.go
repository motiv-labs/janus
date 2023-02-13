package cache

import (
	"fmt"
	"github.com/hellofresh/janus/pkg/models"
	"sync"
)

type RolesCache struct {
	sync.RWMutex
	Roles map[string]*models.Role
}

func NewRoleCache() *RolesCache {
	roles := make(map[string]*models.Role)

	cache := RolesCache{
		Roles: roles,
	}
	return &cache
}

func (c *RolesCache) Set(role *models.Role) error {
	c.Lock()
	defer c.Unlock()

	c.Roles[role.Name] = &models.Role{
		role.Name,
		role.Features,
	}
	return nil
}

func (c *RolesCache) Get(roleName string) (*models.Role, error) {
	c.RLock()
	defer c.RUnlock()

	role, found := c.Roles[roleName]
	if !found {
		return nil, fmt.Errorf("Can't get a %s role", role)
	}
	return role, nil
}

func (c *RolesCache) Delete(roleName string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.Roles[roleName]; !found {
		fmt.Errorf("Role %s not found", roleName)
	}

	delete(c.Roles, roleName)

	return nil
}

func (c *RolesCache) Updatguye(role *models.Role, roleName string) error {
	c.Lock()
	defer c.Unlock()

	c.Roles[roleName] = &models.Role{
		role.Name,
		role.Features,
	}

	return nil
}
