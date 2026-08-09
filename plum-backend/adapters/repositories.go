package adapters

import (
	"fmt"
	"sync"
	"time"

	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/domain"
)

type RepositoriesEntry struct {
	Name       domain.WebsiteName
	Repository *Repository
}

type Repositories struct {
	repositories []RepositoriesEntry
	mutex        sync.RWMutex
}

func NewRepositories() *Repositories {
	return &Repositories{
		repositories: nil,
	}
}

func (r *Repositories) Add(name domain.WebsiteName, repo *Repository) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.get(name); ok {
		return fmt.Errorf("repository '%s' already exists", name)
	}

	r.repositories = append(r.repositories, RepositoriesEntry{
		Name:       name,
		Repository: repo,
	})

	return nil
}

func (r *Repositories) Get(name domain.WebsiteName) (app.Repository, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.get(name)
}

func (r *Repositories) get(name domain.WebsiteName) (*Repository, bool) {
	for _, entry := range r.repositories {
		if entry.Name == name {
			return entry.Repository, true
		}
	}

	return nil, false
}

// RemoveOldData discards the data which is older than the retention period in
// all repositories.
func (r *Repositories) RemoveOldData(now time.Time) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, entry := range r.repositories {
		entry.Repository.RemoveOldData(now)
	}
}

// ForEachData calls the provided function for every hour of the data stored in
// all repositories which falls within the provided range.
func (r *Repositories) ForEachData(from, to time.Time, fn func(t time.Time, data *domain.Data)) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, entry := range r.repositories {
		entry.Repository.forEachData(from, to, fn)
	}
}

func (r *Repositories) Names() []domain.WebsiteName {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var names []domain.WebsiteName

	for _, entry := range r.repositories {
		names = append(names, entry.Name)
	}

	return names
}
