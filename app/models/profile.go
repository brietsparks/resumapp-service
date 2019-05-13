package models

type Profile struct {
	UserId string `db:"user_id" json:"userId"`
	Handle *string `db:"handle" json:"handle"`
	FirstName *string `db:"first_name" json:"firstName"`
	LastName *string `db:"last_name" json:"lastName"`
	Email *string `db:"email" json:"email"`
	Phone *string `db:"phone" json:"phone"`
	ProfilePicUrl *string `db:"profile_pic_url" json:"profilePicUrl"`
}
