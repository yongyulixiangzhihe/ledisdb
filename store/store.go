package store

import (
	"fmt"
	"github.com/siddontang/ledisdb/store/driver"
	"os"
)

const DefaultStoreName = "lmdb"

type Store interface {
	Open(cfg *Config) (driver.IDB, error)
	Repair(cfg *Config) error
}

var dbs = map[string]Store{}

func Register(name string, store Store) {
	if _, ok := dbs[name]; ok {
		panic(fmt.Errorf("db %s is registered", name))
	}

	dbs[name] = store
}

func Open(cfg *Config) (*DB, error) {
	if err := os.MkdirAll(cfg.Path, os.ModePerm); err != nil {
		return nil, err
	}

	if len(cfg.Name) == 0 {
		cfg.Name = DefaultStoreName
	}

	s, ok := dbs[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("db %s is not registered", cfg.Name)
	}

	idb, err := s.Open(cfg)
	if err != nil {
		return nil, err
	}

	db := &DB{idb}

	return db, nil
}

func Repair(cfg *Config) error {
	if len(cfg.Name) == 0 {
		cfg.Name = DefaultStoreName
	}

	s, ok := dbs[cfg.Name]
	if !ok {
		return fmt.Errorf("db %s is not registered", cfg.Name)
	}

	return s.Repair(cfg)
}
