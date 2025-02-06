package enhanced

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
)

type Todos []db.Todo

func (t Todos) Tmpl() tmpl.Template {
	return tmpl.Associated("enhanced/topic", "todos", t)
}

func TopicForm(t db.Topic) tmpl.Template {
	return tmpl.Associated("enhanced/topic", "topic-form", t)
}

type Topic struct {
	Layout
	Topic db.Topic
	Todos Todos
}

func (t Topic) Tmpl() tmpl.Template {
	return tmpl.Wrap(&t.Layout, tmpl.Tmpl("enhanced/topic", t))
}

func (t Todos) PendingTodos() []db.Todo {
	var completed []db.Todo
	for _, todo := range t {
		if !todo.Done {
			completed = append(completed, todo)
		}
	}
	return completed
}

func (t Todos) CompletedTodos() []db.Todo {
	var completed []db.Todo
	for _, todo := range t {
		if todo.Done {
			completed = append(completed, todo)
		}
	}
	return completed
}
