package db

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	Id           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func (d DB) GetUser(id int) (User, error) {
	var user User
	err := d.db.Get(&user, `select * from users where id = $1`, id)
	return user, err
}

func (d DB) GetUserByEmail(email string) (User, error) {
	var user User
	err := d.db.Get(&user, `select * from users where email = $1`, email)
	return user, err
}

func (d DB) InsertUser(email string, passwordHash string) (User, error) {
	var user User
	err := d.db.Get(&user, `insert into users (email, password_hash) values ($1, $2) returning *`, email, passwordHash)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return user, ErrDuplicate
			}
		}
		return user, err
	}
	return user, nil
}

func (d DB) DeleteUser(id int) error {
	res, _ := d.db.Exec(`delete from users where id = $1`, id)
	if count, _ := res.RowsAffected(); count < 1 {
		return ErrNotFound
	}
	return nil
}
