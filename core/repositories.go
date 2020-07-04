package core

import (
	"fmt"
)

type RepositoriesEntry struct {
	Name       string
	Repository *Repository
}

type Repositories struct {
	repositories []RepositoriesEntry
}

func NewRepositories() *Repositories {
	return &Repositories{
		repositories: nil,
	}
}

func (r *Repositories) Add(name string, repo *Repository) error {
	_, ok := r.Get(name)
	if ok {
		return fmt.Errorf("repository '%s' already exists", name)
	}

	r.repositories = append(r.repositories, RepositoriesEntry{
		Name:       name,
		Repository: repo,
	})

	return nil
}

func (r *Repositories) Get(name string) (*Repository, bool) {
	for _, entry := range r.repositories {
		if entry.Name == name {
			return entry.Repository, true
		}
	}

	return nil, false
}

func (r *Repositories) Names() []string {
	var names []string

	for _, entry := range r.repositories {
		names = append(names, entry.Name)
	}

	return names
}
