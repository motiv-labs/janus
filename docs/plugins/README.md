# Extending Janus

Janus can be extended with plugins. Plugins add missing functionality to Janus. They are "plugged in" at compile-time.
The plugins are attached to an [API Definition](../api-definition/README.md) and you can enable or disable them at 
any time.

Janus comes with a set of built in plugins that you can add to your API Definitions: 

* [CORS](cors.md)
* [OAuth2](oauth.md)
* [Rate Limit](rate_limit.md)
* [Request Transformer](request_transformer.md)
* [Compression](compression.md)

## How can I create a plugin?

Even though there are different kinds of plugins, the process of creating one is roughly the same for all.

### 1. Create a package and register your plugin.

Start a new Go package with an init function and register your plugin with Janus:

```go
import "github.com/hellofresh/janus/pkg/plugin"

func init() {
	// register a "generic" plugin, like a directive or middleware
	plugin.RegisterPlugin("name", myPlugin)
}
```

Every plugin must have a name and, when applicable, the name must be unique.

### 2. Plug in your plugin.

To plug your plugin into Janus, import it. This is usually done near the top of [loader.go](../../pgk/loader/loader.go):

```go
import _ "your/plugin/package/path/here"
```

### 3. Write Tests!

Write tests. Get good coverage where possible, and make sure your assertions test what you think they are testing! Use go vet and go test -race to ensure your plugin is as error-free as possible.

### 4. Maintain your plugin.

People will use plugins that are useful, clearly documented, easy to use, and maintained by their owner.
And congratulations, you're a Janus plugin author!
