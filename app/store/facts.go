package store

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type FactsStore struct {
	DB *sqlx.DB
}

type Logger interface {
	Fatal(args ...interface{})
}

func NewFactsStore(DB *sqlx.DB, log Logger) *FactsStore {
    return &FactsStore{DB: DB}
}

func (store *FactsStore) GetFactsByUserId(userId string) (string, error) {
	var facts string

	err := store.DB.Get(&facts, "SELECT facts FROM factdatas WHERE user_id = $1 LIMIT 1", userId, )

	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return facts, nil
}

func (store *FactsStore) UpsertFactsByUserId(userId string, facts string) error {
	_, err := store.DB.Exec("UPSERT INTO factdatas (user_id, facts) VALUES ($1, $2)", userId, facts)

	return err
}
