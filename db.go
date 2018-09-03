package gitdb

import (
	"github.com/ajnavarro/gitdb/git"
	gogit "gopkg.in/src-d/go-git.v4"
)

type DB struct {
	Name       string
	repository *gogit.Repository
}

func NewDB(repository, dbName string) (*DB, error) {
	repo, err := gogit.PlainOpen(repository)
	if err != nil {
		return nil, err
	}

	return &DB{
		Name:       dbName,
		repository: repo,
	}, nil
}

func (db *DB) Table(name string) (*Table, error) {
	return NewTable(
		name,
		db.Name,
		git.NewRepository(db.repository),
	), nil
}

func (db *DB) Sync() error {
	return nil
}
