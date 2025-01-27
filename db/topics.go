package db

import (
	"time"
)

type Topic struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`

	// Aggregated fields

	TodosCount int `db:"todos_count"`
}

func (t Topic) FormatCreatedAt() string {
	return t.CreatedAt.Format("Jan _2, 2006")
}

func (d DB) GetTopics(userId int, search string) ([]Topic, error) {
	var topics []Topic

	err := d.db.Select(&topics, `
	select topics.*, count(todos.topic_id) as todos_count
	from topics
	left join todos on todos.topic_id = topics.id
	where topics.user_id = $1 and topics.title ilike '%' || $2 || '%'
	group by topics.id
	order by topics.id desc
	`, userId, search)
	return topics, err
}

func (d DB) GetTopic(id int) (Topic, error) {
	var topic Topic
	err := d.db.Get(&topic, `
	select topics.*, count(todos.topic_id) as todos_count
	from topics
	left join todos on todos.topic_id = topics.id
	where topics.id = $1
	group by topics.id
	`, id)
	return topic, err
}

func (d DB) InsertTopic(userId int, title string) (Topic, error) {
	var topic Topic
	err := d.db.Get(&topic, `insert into topics (user_id, title) values ($1, $2) returning *`, userId, title)
	if err != nil {
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
