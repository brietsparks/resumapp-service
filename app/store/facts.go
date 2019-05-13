package store

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type FactsStore struct {
	DB *sqlx.DB
}

func NewFactsStore(DB *sqlx.DB) *FactsStore {
    return &FactsStore{DB: DB}
}

func (store *FactsStore) GetFactsByUserId(userId string) (string, error) {
	var facts string

	err := store.DB.Get(&facts, "SELECT facts FROM facts WHERE user_id = $1 LIMIT 1", userId)

	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return facts, nil
}

func (store *FactsStore) UpsertFactsByUserId(userId string, facts string) error {
	_, err := store.DB.Exec("UPSERT INTO facts (user_id, facts) VALUES ($1, $2)", userId, facts)

	return err
}
