package classic_pages

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

func (t Topic) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/topic", t.Layout, t).Template()
}
