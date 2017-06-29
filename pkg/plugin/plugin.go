package plugin

import (
	"encoding/json"
	"errors"
	"sync"

	"fmt"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	stats "github.com/hellofresh/stats-go"
)

var (
	lock sync.RWMutex

	// plugins is a map of plugin name to Plugin.
	plugins = make(map[string]Plugin)
)

// SetupFunc is used to set up a plugin, or in other words,
// execute a directive. It will be called once per key for
// each server block it appears in.
type SetupFunc func(route *proxy.Route, p Params) error

// Params initialization options.
type Params struct {
	Router      router.Router
	Storage     store.Store
	APIRepo     api.Repository
	OAuthRepo   oauth.Repository
	StatsClient stats.Client
	Config      map[string]interface{}
}

// Plugin defines basic methods for plugins
type Plugin struct {
	Action SetupFunc
}

// RegisterPlugin plugs in plugin. All plugins should register
// themselves, even if they do not perform an action associated
// with a directive. It is important for the process to know
// which plugins are available.
//
// The plugin MUST have a name: lower case and one word.
// If this plugin has an action, it must be the name of
// the directive that invokes it. A name is always required
// and must be unique for the server type.
func RegisterPlugin(name string, plugin Plugin) error {
	lock.Lock()
	defer lock.Unlock()

	if name == "" {
		return errors.New("plugin must have a name")
	}
	if _, dup := plugins[name]; dup {
		return fmt.Errorf("plugin named %s  already registered", name)
	}
	plugins[name] = plugin
	return nil
}

// DirectiveAction gets the action for a plugin
func DirectiveAction(name string) (SetupFunc, error) {
	if plugin, ok := plugins[name]; ok {
		return plugin.Action, nil
	}
	return nil, fmt.Errorf("no action found for plugin '%s' (missing a plugin?)", name)
}

// Decode decodes a map string interface into a struct
// for some reasons mapstructure.Decode() gives empty arrays for all resulting config fields
// this is quick workaround hack t make it work
// FIXME: investigate and fix mapstructure.Decode() behaviour and remove this dirty hack
func Decode(rawConfig map[string]interface{}, obj interface{}) error {
	valJSON, err := json.Marshal(rawConfig)
	if nil != err {
		return err
	}
	err = json.Unmarshal(valJSON, obj)
	if nil != err {
		return err
	}

	return nil
}
