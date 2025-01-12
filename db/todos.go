package db

import (
	"time"

	"github.com/lib/pq"
)

type Todo struct {
	Id        int       `db:"id"`
	TopicId   int       `db:"topid_id"`
	Body      string    `db:"body"`
	Done      bool      `db:"done"`
	CreatedAt time.Time `db:"created_at"`
}

func (d DB) GetItems(topicId string) ([]Todo, error) {
	var todos []Todo
	err := d.db.Select(&todos, `select * from todos where topic_id = $1`, topicId)
	return todos, err
}

func (d DB) GetTodo(id int) (Todo, error) {
	var todo Todo
	err := d.db.Get(&todo, `select * from todos where id = $1`, id)
	return todo, err
}

func (d DB) InsertTodo(topicId int, body string) (Todo, error) {
	var todo Todo
	err := d.db.Get(&todo, `insert into todos (topic_id, body) values ($1, $2) returning *`, topicId, body)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return todo, ErrDuplicate
			}
		}
		return todo, err
	}
	return todo, nil
}

func (d DB) UpdateTodo(id int, body string, done bool) (Todo, error) {
	var todo Todo
	err := d.db.Get(&todo, `update todos set body = $1, done = $2 where id = $3 returning *`, body, done, id)
	return todo, err
}

func (d DB) DeleteTodo(id int) error {
	res, _ := d.db.Exec(`delete from todos where id = $1`, id)
	if count, _ := res.RowsAffected(); count < 1 {
		return ErrNotFound
	}
	return nil
}
