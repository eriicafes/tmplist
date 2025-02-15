package db

import (
	"time"

	"github.com/lib/pq"
)

type Session struct {
	Id        string    `db:"id" json:"id"`
	UserId    int       `db:"user_id" json:"user_id"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}

func (s Session) GetExpirestAt() time.Time {
	return s.ExpiresAt
}

type SessionStorage struct {
	DB
}

func (s SessionStorage) CreateSession(id string, user User, expiresAt time.Time) (Session, error) {
	session := Session{
		Id:        id,
		UserId:    user.Id,
		ExpiresAt: expiresAt,
	}
	return session, s.InsertSession(session)
}

func (s SessionStorage) UpdateSession(session Session, expiresAt time.Time) (Session, error) {
	session.ExpiresAt = expiresAt
	return session, s.DB.UpdateSession(session)
}

func (d DB) GetSessionAndUser(id string) (Session, User, error) {
	var s Session
	var u User
	row := d.db.QueryRowx(`
	select
		s.id,
		s.user_id,
		s.expires_at,
		u.id as u_id,
		u.email as u_email,
		u.password_hash as u_password_hash,
		u.created_at as u_created_at
	from sessions as s
	join users as u on u.id = s.user_id
	where s.id = $1
	`, id)
	err := row.Scan(&s.Id, &s.UserId, &s.ExpiresAt, &u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return s, u, err
}

func (d DB) InsertSession(session Session) error {
	_, err := d.db.NamedExec(`insert into sessions (id, user_id, expires_at) values (:id, :user_id, :expires_at)`, session)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return ErrDuplicate
			}
		}
		return err
	}
	return nil
}

func (d DB) UpdateSession(session Session) error {
	_, err := d.db.NamedExec(`update sessions set expires_at = :expires_at where id = :id`, session)
	return err
}

func (d DB) DeleteSession(id string) error {
	res, _ := d.db.Exec(`delete from sessions where id = $1`, id)
	if count, _ := res.RowsAffected(); count < 1 {
		return ErrNotFound
	}
	return nil
}
