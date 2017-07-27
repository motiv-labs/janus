package api

// Repository defines the behavior of a proxy specs repository
type Repository interface {
	FindAll() ([]*Definition, error)
	FindByName(name string) (*Definition, error)
	FindByListenPath(path string) (*Definition, error)
	Exists(def *Definition) (bool, error)
	Add(app *Definition) error
	Remove(name string) error
	FindValidAPIHealthChecks() ([]*Definition, error)
}

func exists(r Repository, def *Definition) (bool, error) {
	_, err := r.FindByName(def.Name)
	if nil != err && err != ErrAPIDefinitionNotFound {
		return false, err
	} else if err != ErrAPIDefinitionNotFound {
		return true, ErrAPINameExists
	}

	_, err = r.FindByListenPath(def.Proxy.ListenPath)
	if nil != err && err != ErrAPIDefinitionNotFound {
		return false, err
	} else if err != ErrAPIDefinitionNotFound {
		return true, ErrAPIListenPathExists
	}

	return false, nil
}
