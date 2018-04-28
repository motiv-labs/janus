package api

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

const (
	tableName = "apis"
)

// PostgresRepository represents a postgres repository
type PostgresRepository struct {
	sess sqlbuilder.Database
}

// NewPostgresRepository creates a postgres API definition repo
func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	settings, err := postgresql.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	sess, err := postgresql.Open(settings)
	if err != nil {
		return nil, err
	}
	//sess.SetLogging(true)
	return &PostgresRepository{
		sess: sess,
	}, nil
}

// Close terminates the session.  It's a runtime error to use a session
// after it has been closed.
func (r PostgresRepository) Close() error {
	return r.sess.Close()
}

// Listen watches for changes on the configuration
func (r PostgresRepository) Listen(ctx context.Context, cfgChan <-chan ConfigurationMessage) {
	go func() {
		log.Debug("Listening for changes on the provider...")
		for {
			select {
			case cfg := <-cfgChan:
				switch cfg.Operation {
				case AddedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not add the configuration on the provider")
					}
				case UpdatedOperation:
					err := r.update(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not update the configuration on the provider")
					}
				case RemovedOperation:
					err := r.remove(cfg.Configuration.Name)
					if err != nil {
						log.WithError(err).Error("Could not remove the configuration from the provider")
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Watch watches for changes on the database
func (r PostgresRepository) Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged) {
	t := time.NewTicker(time.Second * 5)
	go func(refreshTicker *time.Ticker) {
		defer refreshTicker.Stop()
		log.Debug("Watching Provider...")
		for {
			select {
			case <-refreshTicker.C:
				defs, err := r.FindAll()
				if err != nil {
					log.WithError(err).Error("Failed to get configurations on watch")
					continue
				}

				cfgChan <- ConfigurationChanged{
					Configurations: &Configuration{Definitions: defs},
				}
			case <-ctx.Done():
				return
			}
		}
	}(t)
}

// FindAll fetches all the API definitions available
func (r PostgresRepository) FindAll() ([]*Definition, error) {
	var src []map[string]interface{}
	if err := r.sess.SelectFrom(tableName).All(&src); err != nil {
		return nil, err
	}
	var result []*Definition
	for _, m := range src {
		var dst map[string]interface{}
		if err := postgresql.ScanJSONB(&dst, m["json"]); err != nil {
			return nil, err
		}
		b, err := json.Marshal(dst)
		if err != nil {
			return nil, err
		}
		definition := NewDefinition()
		if err := json.Unmarshal(b, &definition); err != nil {
			return nil, err
		}
		result = append(result, definition)
	}
	return result, nil
}

func (r PostgresRepository) add(definition *Definition) error {
	value, err := postgresql.JSONBValue(definition)
	if err != nil {
		return err
	}
	_, err = r.sess.InsertInto(tableName).Values(value).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresRepository) update(definition *Definition) error {
	value, err := postgresql.JSONBValue(definition)
	if err != nil {
		return err
	}
	_, err = r.sess.Update(tableName).Set("json", value).Where("json->>'name' = ?", definition.Name).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresRepository) remove(name string) error {
	_, err := r.sess.DeleteFrom(tableName).Where("json->>'name' = ?", name).Exec()
	if err != nil {
		return err
	}
	return nil
}
