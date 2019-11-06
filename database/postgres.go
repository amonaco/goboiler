package database

import (
	"github.com/go-pg/pg"
	"github.com/pelletier/go-toml"
)

// DBConn returns a postgres connection pool.
func DBConn(config *toml.Tree) (*pg.DB, error) {

	database_url := config.Get("postgresql.url").(string)

	opts, err := pg.ParseURL(database_url)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)
	if err := checkConn(db); err != nil {
		return nil, err
	}

	return db, nil
}

func checkConn(db *pg.DB) error {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err
}
