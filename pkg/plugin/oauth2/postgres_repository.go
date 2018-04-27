package oauth2

import "sync"

type PostgresRepository struct {
	sync.Mutex
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	return &PostgresRepository{}, nil
}

func (r PostgresRepository) FindAll() ([]*OAuth, error) {
	return []*OAuth{}, nil
}

func (r PostgresRepository) FindByName(name string) (*OAuth, error) {
	return &OAuth{}, nil
}

func (r PostgresRepository) Add(oauth *OAuth) error {
	return nil
}

func (r PostgresRepository) Save(oauth *OAuth) error {
	return nil
}

func (r PostgresRepository) Remove(id string) error {
	return nil
}
