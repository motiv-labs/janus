package oauth

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	sync.Mutex
	servers map[string]*OAuth
}

// NewFileSystemRepository creates a mongo OAuth Server repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := &FileSystemRepository{servers: make(map[string]*OAuth)}
	// Grab json files from directory
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		return nil, err
	}

	for _, f := range files {
		filePath := filepath.Join(dir, f.Name())
		definition := new(OAuth)

		v := viper.New()
		v.SetConfigFile(filePath)
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Debug("OAauth2 configuration changed, reloading...")
			if err := v.Unmarshal(definition); err != nil {
				log.WithError(err).Error("Can't unmarshal the OAauth2 configuration")
			}
		})

		if err := v.ReadInConfig(); err != nil {
			log.WithError(err).Error("Couldn't load the OAauth2 definition file")
			return nil, err
		}

		if err := v.Unmarshal(definition); err != nil {
			return nil, err
		}

		if err = repo.Add(definition); err != nil {
			log.WithError(err).Error("Can't add the definition to the repository")
			return nil, err
		}
	}

	return repo, nil
}

// FindAll fetches all the OAuth Servers available
func (r *FileSystemRepository) FindAll() ([]*OAuth, error) {
	var servers []*OAuth
	for _, server := range r.servers {
		servers = append(servers, server)
	}

	return servers, nil
}

// FindByName find an OAuth Server by name
func (r *FileSystemRepository) FindByName(name string) (*OAuth, error) {
	server, ok := r.servers[name]
	if false == ok {
		return nil, ErrOauthServerNotFound
	}

	return server, nil
}

// Add adds an OAuth Server to the repository
func (r *FileSystemRepository) Add(server *OAuth) error {
	r.Lock()
	defer r.Unlock()

	r.servers[server.Name] = server

	return nil
}

// Remove removes an OAuth Server from the repository
func (r *FileSystemRepository) Remove(name string) error {
	r.Lock()
	defer r.Unlock()

	delete(r.servers, name)
	return nil
}

// FindByTokenURL returns OAuth Server records with corresponding token url
func (r *FileSystemRepository) FindByTokenURL(url url.URL) (*OAuth, error) {
	for _, server := range r.servers {
		if server.Endpoints.Token.UpstreamURL == url.String() {
			return server, nil
		}
	}

	return nil, ErrOauthServerNotFound
}
