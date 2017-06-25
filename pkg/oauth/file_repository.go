package oauth

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	*InMemoryRepository
	sync.Mutex
}

// NewFileSystemRepository creates a mongo OAuth Server repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := &FileSystemRepository{InMemoryRepository: NewInMemoryRepository()}

	// Grab json files from directory
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		return nil, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			filePath := filepath.Join(dir, f.Name())
			oauthServerRaw, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.WithError(err).WithField("path", filePath).Error("Couldn't load the oauth server file")
				return nil, err
			}

			oauthServer := repo.parseOAuthServer(oauthServerRaw)
			if err = repo.Add(oauthServer); err != nil {
				log.WithError(err).Error("Can't add the definition to the repository")
				return nil, err
			}
		}
	}

	return repo, nil
}

func (r *FileSystemRepository) parseOAuthServer(oauthServerRaw []byte) *OAuth {
	oauthServer := new(OAuth)
	if err := json.Unmarshal(oauthServerRaw, oauthServer); err != nil {
		log.WithError(err).Error("[RPC] --> Couldn't unmarshal oauth server configuration")
	}

	return oauthServer
}
