package oauth2

import (
	"encoding/json"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

const (
	tableName = "oauth"
)

// PostgresRepository represents a postgres repository
type PostgresRepository struct {
	sess sqlbuilder.Database
}

// NewPostgresRepository creates a postgres OAuth Server repo
func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	settings, err := postgresql.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	sess, err := postgresql.Open(settings)
	if err != nil {
		return nil, err
	}
	// sess.SetLogging(true)
	return &PostgresRepository{
		sess: sess,
	}, nil
}

// FindAll fetches all the OAuth Servers available
func (r PostgresRepository) FindAll() ([]*OAuth, error) {
	var src []map[string]interface{}
	if err := r.sess.SelectFrom(tableName).All(&src); err != nil {
		return nil, err
	}
	var result []*OAuth
	for _, m := range src {
		var dst map[string]interface{}
		if err := postgresql.ScanJSONB(&dst, m["json"]); err != nil {
			return nil, err
		}
		b, err := json.Marshal(dst)
		if err != nil {
			return nil, err
		}
		oauth := NewOAuth()
		if err := json.Unmarshal(b, &oauth); err != nil {
			return nil, err
		}
		result = append(result, oauth)
	}
	return result, nil
}

// FindByName find an OAuth Server by name
func (r PostgresRepository) FindByName(name string) (*OAuth, error) {
	var src map[string]interface{}
	if err := r.sess.SelectFrom(tableName).Where("json->>'name' = ?", name).One(&src); err != nil {
		return nil, err
	}
	var dst map[string]interface{}
	if err := postgresql.ScanJSONB(&dst, src["json"]); err != nil {
		return nil, err
	}
	b, err := json.Marshal(dst)
	if err != nil {
		return nil, err
	}
	oauth := NewOAuth()
	if err := json.Unmarshal(b, &oauth); err != nil {
		return nil, err
	}
	return oauth, nil
}

// Add add a new OAuth Server to the repository
func (r PostgresRepository) Add(oauth *OAuth) error {
	value, err := postgresql.JSONBValue(oauth)
	if err != nil {
		return err
	}
	_, err = r.sess.InsertInto(tableName).Values(value).Exec()
	if err != nil {
		return err
	}
	return nil
}

// Save saves OAuth Server to the repository
func (r PostgresRepository) Save(oauth *OAuth) error {
	value, err := postgresql.JSONBValue(oauth)
	if err != nil {
		return err
	}
	_, err = r.sess.Update(tableName).Set("json", value).Where("json->>'name' = ?", oauth.Name).Exec()
	if err != nil {
		return err
	}
	return nil
}

// Remove removes an OAuth Server from the repository
func (r PostgresRepository) Remove(name string) error {
	_, err := r.sess.DeleteFrom(tableName).Where("json->>'name' = ?", name).Exec()
	if err != nil {
		return err
	}
	return nil
}
