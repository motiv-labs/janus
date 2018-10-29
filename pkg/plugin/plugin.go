package plugin

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	lock sync.RWMutex

	// plugins is a map of plugin name to Plugin.
	plugins = make(map[string]Plugin)

	// eventHooks is a map of hook name to Hook. All hooks plugins
	// must have a name.
	eventHooks = make(map[string][]EventHook)
)

// SetupFunc is used to set up a plugin, or in other words,
// execute a directive. It will be called once per key for
// each server block it appears in.
type SetupFunc func(def *proxy.RouterDefinition, rawConfig Config) error

// ValidateFunc validates configuration data against the plugin struct
type ValidateFunc func(rawConfig Config) (bool, error)

// Config initialization options.
type Config map[string]interface{}

// Plugin defines basic methods for plugins
type Plugin struct {
	Action   SetupFunc
	Validate ValidateFunc
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

// EventHook is a type which holds information about a startup hook plugin.
type EventHook func(event interface{}) error

// RegisterEventHook plugs in hook. All the hooks should register themselves
// and they must have a name.
func RegisterEventHook(name string, hook EventHook) error {
	log.WithField("event_name", name).Debug("Event registered")
	lock.Lock()
	defer lock.Unlock()

	if name == "" {
		return errors.New("event hook must have a name")
	}

	if hooks, dup := eventHooks[name]; dup {
		eventHooks[name] = append(hooks, hook)
	} else {
		eventHooks[name] = append([]EventHook{}, hook)
	}

	return nil
}

// EmitEvent executes the different hooks passing the EventType as an
// argument. This is a blocking function. Hook developers should
// use 'go' keyword if they don't want to block Janus.
func EmitEvent(name string, event interface{}) error {
	log.WithField("event_name", name).Debug("Event triggered")

	hooks, found := eventHooks[name]
	if !found {
		return errors.New("Plugin not found")
	}

	for _, hook := range hooks {
		err := hook(event)
		if err != nil {
			log.WithError(err).WithField("event_name", name).Warn("an error occurred when an event was triggered")
		}
	}

	return nil
}

// ValidateConfig validates the plugin configuration data
func ValidateConfig(name string, rawConfig Config) (bool, error) {
	if plugin, ok := plugins[name]; ok {
		if plugin.Validate == nil {
			return true, nil
		}

		result, err := plugin.Validate(rawConfig)
		return result, err
	}
	return false, fmt.Errorf("no validate function found for plugin '%s'", name)
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
