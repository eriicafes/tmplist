package classic

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
)

type Topic struct {
	Layout
	Topic         db.Topic
	Todos         []db.Todo
	LastUpdatedId int
}

func (t Topic) Tmpl() tmpl.Template {
	return tmpl.Wrap(&t.Layout, tmpl.Tmpl("classic/topic", t))
}

func (t Topic) PendingTodos() []db.Todo {
	var completed []db.Todo
	for _, todo := range t.Todos {
		if !todo.Done {
			completed = append(completed, todo)
		}
	}
	return completed
}

func (t Topic) CompletedTodos() []db.Todo {
	var completed []db.Todo
	for _, todo := range t.Todos {
		if todo.Done {
			completed = append(completed, todo)
		}
	}
	return completed
}
