package db

import (
	"time"

	"github.com/lib/pq"
)

type Topic struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
}

func (d DB) GetTopics(userId int) ([]Topic, error) {
	var topics []Topic
	err := d.db.Select(&topics, `select * from topics where user_id = $1`, userId)
	return topics, err
}

func (d DB) GetTopic(id int) (Topic, error) {
	var topic Topic
	err := d.db.Get(&topic, `select * from topics where id = $1`, id)
	return topic, err
}

func (d DB) InsertTopic(userId int, title string) (Topic, error) {
	var topic Topic
	err := d.db.Get(&topic, `insert into topics (user_id, title) values ($1, $2) returning *`, userId, title)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return topic, ErrDuplicate
			}
		}
		return topic, err
	}
	return topic, nil
}

func (d DB) UpdateTopic(id int, title string) (Topic, error) {
	var topic Topic
	err := d.db.Get(&topic, `update topics set title = $1 where id = $2 returning *`, title, id)
	return topic, err
}

func (d DB) DeleteTopic(id int) error {
	res, _ := d.db.Exec(`delete from topics where id = $1`, id)
	if count, _ := res.RowsAffected(); count < 1 {
		return ErrNotFound
	}
	return nil
}
