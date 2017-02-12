package provider

import (
	"fmt"
	"log"
	"sync"
	"errors"
)

var instance *manager
var once sync.Once

type manager struct {
	container map[string]UserFetcher
}

// Manager returns the concurrency safe singleton
// instance of a providers manager that has
// an overview over the current registered providers
// of user data tables.
func Manager() *manager {
	once.Do(func() {
		instance = newManager()
	})
	return instance
}

func newManager() *manager {
	return &manager{
		container: make(map[string]UserFetcher),
	}
}

// Register associates a data provider with
// a unique identifier name.
func (m *manager) Register(name string, provider UserFetcher) error {
	if _, ok := m.container[name]; ok {
		err := errors.New("A provider was already registered under that name.")
		log.Println(err)
		return err
	}
	m.container[name] = provider
	return nil
}

// Get returns the data provider associated with
// the given name or error if there is none.
func (m *manager) Get(name string) (UserFetcher, error) {
	if p, ok := m.container[name]; !ok {
		err := fmt.Errorf("There is no provider registered under the name: %s.", name)
		log.Println(err)
		return nil, err
	} else {
		return p, nil
	}
}

// All returns all registered providers along
// with their name.
func (m *manager) All() map[string]UserFetcher {
	return m.container
}
