package store

import (
	"database/sql"
	"github.com/brietsparks/resumapp-service/app/models"
	"github.com/jmoiron/sqlx"
)

type ProfilesStore struct {
	DB *sqlx.DB
}

func NewProfilesStore(DB *sqlx.DB) *ProfilesStore {
	return &ProfilesStore{DB: DB}
}

func (store *ProfilesStore) GetHandleAvailability(handle string) (bool, error) {
	profile, err := store.GetProfileByHandle(handle)

	isAvailable := profile == nil

	return isAvailable, err
}

func (store *ProfilesStore) GetProfileByUserId(userId string) (*models.Profile, error) {
	profile := models.Profile{}

	err := store.DB.Get(&profile, "SELECT * FROM profiles WHERE user_id = $1 LIMIT 1", userId)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &profile, err
}

func (store *ProfilesStore) GetProfileByHandle(handle string) (*models.Profile, error) {
	profile := models.Profile{}

	err := store.DB.Get(&profile, "SELECT * FROM profiles WHERE handle = $1 LIMIT 1", handle)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &profile, err
}

func (store *ProfilesStore) UpsertProfileByUserId(userId string, profile models.Profile) error {
	m := map[string]interface{}{
		"userId":        profile.UserId,
		"handle":        profile.Handle,
		"firstName":     profile.FirstName,
		"lastName":      profile.LastName,
		"email":         profile.Email,
		"phone":         profile.Phone,
		"profilePicUrl": profile.ProfilePicUrl,
	}

	_, err := store.DB.NamedExec(`
		UPSERT INTO profiles (
			user_id, handle, first_name, last_name, email, phone, profile_pic_url
		) VALUES (
			:userId, :handle, :firstName, :lastName, :email, :phone, :profilePicUrl
		)
	`, m)

	return err
}
