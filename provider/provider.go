package provider

import (
	"github.com/g-leon/go-data-collector/user"
)

// UserFetcher represents an endpoint containing
// one ore more tables with user data.
type UserFetcher interface {

	// TableNames returns the list of names
	// of the tables that can be found at
	// this specific endpoint.
	TableNames() []string

	// GetTable returns the users data
	// from the table identified by the
	// name given as parameter or an
	// error if the table could not be
	// found.
	GetTable(name string) ([]*user.Model, error)
}
