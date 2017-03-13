package oauth

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"net/url"

	log "github.com/Sirupsen/logrus"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	servers map[string]*OAuth
}

// NewFileSystemRepository creates a mongo country repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := &FileSystemRepository{make(map[string]*OAuth)}
	// Grab json files from directory
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		return nil, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			filePath := filepath.Join(dir, f.Name())
			log.Info("Loading API Specification from ", filePath)
			appConfigBody, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.WithError(err).Error("Couldn't load app configuration file")
				return nil, err
			}

			definition := repo.parseDefinition(appConfigBody)
			err = repo.Add(definition)
			if err != nil {
				log.WithError(err).Error("Can't add the definition to the repository")
				return nil, err
			}
		}
	}

	return repo, nil
}

// FindAll fetches all the countries available
func (r *FileSystemRepository) FindAll() ([]*OAuth, error) {
	var servers []*OAuth
	for _, server := range r.servers {
		servers = append(servers, server)
	}

	return servers, nil
}

// FindBySlug find a country by the iso2code provided
func (r *FileSystemRepository) FindBySlug(slug string) (*OAuth, error) {
	server, ok := r.servers[slug]
	if false == ok {
		return nil, ErrOauthServerNotFound
	}

	return server, nil
}

// Add adds a country to the repository
func (r *FileSystemRepository) Add(server *OAuth) error {
	r.servers[server.Slug] = server

	return nil
}

// Remove removes a country from the repository
func (r *FileSystemRepository) Remove(slug string) error {
	delete(r.servers, slug)
	return nil
}

// FindByTokenURL returns OAuth server records with corresponding token url
func (r *FileSystemRepository) FindByTokenURL(url url.URL) (*OAuth, error) {
	for _, server := range r.servers {
		if server.Endpoints.Token.TargetURL == url.String() {
			return server, nil
		}
	}

	return nil, ErrOauthServerNotFound
}

func (r *FileSystemRepository) parseDefinition(apiDef []byte) *OAuth {
	appConfig := &OAuth{}
	if err := json.Unmarshal(apiDef, appConfig); err != nil {
		log.Error("[RPC] --> Couldn't unmarshal api configuration: ", err)
	}

	return appConfig
}
