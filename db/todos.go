package db

import (
	"time"
)

type Todo struct {
	Id        int       `db:"id"`
	TopicId   int       `db:"topic_id"`
	Body      string    `db:"body"`
	Done      bool      `db:"done"`
	CreatedAt time.Time `db:"created_at"`
}

func (d DB) GetTodos(topicId int) ([]Todo, error) {
	var todos []Todo
	err := d.db.Select(&todos, `select * from todos where topic_id = $1 order by id desc`, topicId)
	return todos, err
}

func (d DB) GetTodo(id int) (Todo, error) {
	var todo Todo
	err := d.db.Get(&todo, `select * from todos where id = $1`, id)
	return todo, err
}

func (d DB) InsertTodos(todos []Todo) (int64, error) {
	result, err := d.db.NamedExec(`insert into todos (topic_id, body, done) values (:topic_id, :body, :done)`, todos)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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
