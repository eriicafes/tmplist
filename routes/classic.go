package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/internal/session"
	"github.com/eriicafes/tmplist/schemas"
	classic_pages "github.com/eriicafes/tmplist/templates/classic/pages"
)

func (c Context) Classic(mux internal.Mux) {
	mux = internal.Fallback(mux, c.ClassicErrorHandler())
	auth := internal.Use(mux, c.authMiddleware())
	guest := internal.Use(mux, c.guestMiddleware())

	toastMessage := session.NewFlash[classic_pages.Toast](session.FlashOptions{
		Cookie: "toast_message",
		Secure: c.Prod,
		Path:   "/",
	})
	lastUpdatedId := session.NewFlash[int](session.FlashOptions{
		Cookie: "last_updated_id",
		Secure: c.Prod,
		Path:   "/",
	})

	// all topics page
	auth.Route("GET ", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		user, _ := requestUser.Get(r.Context())

		// get topics from db
		topics, _ := c.DB.GetTopics(user.Id)

		return tr.Render(w, classic_pages.Index{
			Layout: classic_pages.Layout{
				Toast: toastMessage.Get(w, r),
				Title: "Topics",
				User:  &user,
			},
			Topics: topics,
		})
	})

	// create topic
	auth.Route("POST ", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())

		renderError := func(message string) error {
			toastMessage.Set(w, classic_pages.Toast{
				Message: message,
				Type:    classic_pages.ToastError,
			})
			http.Redirect(w, r, "/classic", http.StatusFound)
			return nil
		}

		// validate input
		form := schemas.TopicData{
			Topic: r.PostFormValue("topic"),
		}
		todos, todosChecked := r.PostForm["todo"], r.PostForm["todo-checked"]
		// check if todos and todos-checked have the same length
		if len(todos) != len(todosChecked) {
			return renderError("Invalid input")
		}
		// add structured todos to form
		for i, v := range slices.Backward(todos) {
			form.Todos = append(form.Todos, schemas.TodoData{
				Text:    v,
				Checked: todosChecked[i] == "on",
			})
		}
		// return form validation errors
		if errors := schemas.FormErrors(form.Validate()); errors != nil {
			return renderError("Invalid input")
		}

		// add topic to db
		topic, err := c.DB.InsertTopic(user.Id, form.Topic)
		if err != nil {
			log.Println(err)
			return renderError("Failed to create topic")
		}
		if len(form.Todos) > 0 {
			var insertTodos []db.Todo
			for _, v := range form.Todos {
				insertTodos = append(insertTodos, db.Todo{
					TopicId: topic.Id,
					Body:    v.Text,
					Done:    v.Checked,
				})
			}
			// add topic todos to db
			_, err = c.DB.InsertTodos(insertTodos)
			if err != nil {
				log.Println(err)
				return renderError("Failed to create todos")
			}
		}

		// set flash message and redirect
		toastMessage.Set(w, classic_pages.Toast{
			Message: "Topic created",
		})
		http.Redirect(w, r, "/classic", http.StatusFound)
		return nil
	})

	// topic page
	auth.Route("GET /{topicId}", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			log.Println(err)
			return httperrors.New("Topic not found", http.StatusNotFound)
		}
		// get todos for topic
		todos, err := c.DB.GetTodos(topic.Id)
		if err != nil {
			log.Println(err)
		}

		return tr.Render(w, classic_pages.Topic{
			Layout: classic_pages.Layout{
				Toast: toastMessage.Get(w, r),
				Title: topic.Title,
				User:  &user,
			},
			Topic:         topic,
			Todos:         todos,
			LastUpdatedId: lastUpdatedId.Get(w, r),
		})
	})

	// update topic
	auth.Route("POST /{topicId}", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))

		renderError := func(message string) error {
			toastMessage.Set(w, classic_pages.Toast{
				Message: message,
				Type:    classic_pages.ToastError,
			})
			http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
			return nil
		}

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			log.Println(err)
			return renderError("Topic not found")
		}

		// validate input
		form := schemas.TopicData{
			Topic: r.PostFormValue("topic"),
		}
		// return form validation errors
		if errors := schemas.FormErrors(form.Validate()); errors != nil {
			return renderError("Invalid input")
		}

		// update todo in db
		if _, err = c.DB.UpdateTopic(topic.Id, form.Topic); err != nil {
			log.Println(err)
			return renderError("Failed to update topic")
		}

		http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
		return nil
	})

	// delete topic
	auth.Route("POST /{topicId}/delete", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))

		renderError := func(message string) error {
			toastMessage.Set(w, classic_pages.Toast{
				Message: message,
				Type:    classic_pages.ToastError,
			})
			http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
			return nil
		}

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			log.Println(err)
			return renderError("Topic not found")
		}

		// delete todo from db
		if err = c.DB.DeleteTopic(topic.Id); err != nil {
			log.Println(err)
			return renderError("Failed to delete topic")
		}

		// set flash message and redirect
		toastMessage.Set(w, classic_pages.Toast{
			Message: "Topic deleted",
		})
		http.Redirect(w, r, "/classic", http.StatusFound)
		return nil
	})

	// create todo
	auth.Route("POST /{topicId}/todos", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))

		renderError := func(message string) error {
			toastMessage.Set(w, classic_pages.Toast{
				Message: message,
				Type:    classic_pages.ToastError,
			})
			http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
			return nil
		}

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			log.Println(err)
			return renderError("Topic not found")
		}

		// validate input
		form := schemas.TodoData{
			Text: r.PostFormValue("todo"),
		}
		// return form validation errors
		if errors := schemas.FormErrors(form.Validate()); errors != nil {
			return renderError("Invalid input")
		}

		// update todo in db
		insertTodos := []db.Todo{{TopicId: topic.Id, Body: form.Text}}
		if _, err = c.DB.InsertTodos(insertTodos); err != nil {
			log.Println(err)
			return renderError("Failed to create todo")
		}

		http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
		return nil
	})

	// update todo
	auth.Route("POST /{topicId}/todos/{todoId}", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))
		todoId, _ := strconv.Atoi(r.PathValue("todoId"))

		renderError := func(message string) error {
			toastMessage.Set(w, classic_pages.Toast{
				Message: message,
				Type:    classic_pages.ToastError,
			})
			http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
			return nil
		}

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			log.Println(err)
			return renderError("Topic not found")
		}
		// check if todo exists and belongs to topic
		todo, err := c.DB.GetTodo(todoId)
		if err != nil || todo.TopicId != topic.Id {
			log.Println(err)
			return renderError("Todo not found")
		}

		// validate input
		form := schemas.TodoData{
			Text:    r.PostFormValue("todo"),
			Checked: r.PostFormValue("todo-checked") == "on",
		}
		// return form validation errors
		if errors := schemas.FormErrors(form.Validate()); errors != nil {
			return renderError("Invalid input")
		}

		// update todo in db
		if _, err = c.DB.UpdateTodo(todo.Id, form.Text, form.Checked); err != nil {
			log.Println(err)
			return renderError("Failed to update todo")
		}

		// set flash message and redirect
		lastUpdatedId.Set(w, todo.Id)
		http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
		return nil
	})

	// delete todo
	auth.Route("POST /{topicId}/todos/{todoId}/delete", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))
		todoId, _ := strconv.Atoi(r.PathValue("todoId"))

		renderError := func(message string) error {
			toastMessage.Set(w, classic_pages.Toast{
				Message: message,
				Type:    classic_pages.ToastError,
			})
			http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
			return nil
		}

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			log.Println(err)
			return renderError("Topic not found")
		}
		// check if todo exists and belongs to topic
		todo, err := c.DB.GetTodo(todoId)
		if err != nil || todo.TopicId != topic.Id {
			log.Println(err)
			return renderError("Todo not found")
		}

		// update todo in db
		if err = c.DB.DeleteTodo(todo.Id); err != nil {
			log.Println(err)
			return renderError("Failed to delete todo")
		}

		http.Redirect(w, r, fmt.Sprintf("/classic/%d", topicId), http.StatusFound)
		return nil
	})

	// login page
	guest.Route("GET /login", func(w http.ResponseWriter, r *http.Request) error {
		return c.Renderer().Render(w, classic_pages.Login{
			Layout: classic_pages.Layout{
				Title: "Login",
			},
		})
	})

	// login post
	guest.Route("POST /login", func(w http.ResponseWriter, r *http.Request) error {
		// prevent other origins from authenticating the user
		if !c.allowOriginForNonSafeRequests(r) {
			return httperrors.New("Login attempt from an unknown site blocked", http.StatusForbidden)
		}

		tr := c.Renderer()

		// validate input
		form := schemas.LoginData{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		renderError := func(message string, details httperrors.Details) error {
			return tr.Render(w, classic_pages.Login{
				Layout: classic_pages.Layout{
					Title: "Login",
				},
				Form:   form,
				Error:  message,
				Errors: details,
			})
		}
		// return form validation errors
		if errors := schemas.FormErrors(form.Validate()); errors != nil {
			return renderError("", errors)
		}

		// get user from db
		user, err := c.DB.GetUserByEmail(form.Email)
		if err != nil {
			return renderError("Email address not found", nil)
		}
		// check password
		if !c.Auth.ComparePassword(user.PasswordHash, form.Password) {
			return tr.Render(w, classic_pages.Login{
				Layout: classic_pages.Layout{
					Title: "Login",
				},
				Form:  form,
				Error: "Invalid password",
			})
		}
		// create session and set cookie
		token, err := c.Auth.GenerateSessionToken()
		if err != nil {
			return renderError("Failed to login", nil)
		}
		session, err := c.Auth.CreateSession(token, user)
		if err != nil {
			return renderError("Failed to login", nil)
		}
		c.Auth.SetCookie(w, token, session.ExpiresAt)

		// redirect to authenticated page
		http.Redirect(w, r, "/classic", http.StatusFound)
		return nil
	})

	// register page
	guest.Route("GET /register", func(w http.ResponseWriter, r *http.Request) error {
		return c.Renderer().Render(w, classic_pages.Register{
			Layout: classic_pages.Layout{
				Title: "Register",
			},
		})
	})

	// register post
	guest.Route("POST /register", func(w http.ResponseWriter, r *http.Request) error {
		// prevent other origins from authenticating the user
		if !c.allowOriginForNonSafeRequests(r) {
			return httperrors.New("Register attempt from an unknown site blocked", http.StatusForbidden)
		}

		tr := c.Renderer()

		// validate input
		form := schemas.RegisterData{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		renderError := func(message string, details httperrors.Details) error {
			return tr.Render(w, classic_pages.Register{
				Layout: classic_pages.Layout{
					Title: "Login",
				},
				Form:   form,
				Error:  message,
				Errors: details,
			})
		}
		// return form validation errors
		if errors := schemas.FormErrors(form.Validate()); errors != nil {
			return renderError("", errors)
		}

		// hash password
		passwordHash, err := c.Auth.HashPassword(form.Password)
		if err != nil {
			return renderError("Failed to create account", nil)
		}
		// insert user in db
		user, err := c.DB.InsertUser(form.Email, passwordHash)
		if err != nil {
			msg := "Failed to create account"
			if err == db.ErrDuplicate {
				msg = "Email address already taken"
			}
			return renderError(msg, nil)
		}
		// create session and set cookie
		token, err := c.Auth.GenerateSessionToken()
		if err != nil {
			return renderError("Failed to create account", nil)
		}
		session, err := c.Auth.CreateSession(token, user)
		if err != nil {
			return renderError("Failed to create account", nil)
		}
		c.Auth.SetCookie(w, token, session.ExpiresAt)

		// set flash message and redirect to authenticated page
		toastMessage.Set(w, classic_pages.Toast{
			Message: "Account created successfully",
			Type:    classic_pages.ToastSuccess,
		})
		http.Redirect(w, r, "/classic", http.StatusFound)
		return nil
	})

	auth.Route("POST /logout", func(w http.ResponseWriter, r *http.Request) error {
		session, _ := requestSession.Get(r.Context())
		if err := c.Auth.InvalidateSession(session.Id); err != nil {
			return err
		}
		c.Auth.DeleteCookie(w)
		http.Redirect(w, r, "/classic/login", http.StatusFound)
		return nil
	})

	// 404
	mux.HandleFunc("/", internal.ErrorHandler(mux, httperrors.New("page not found", http.StatusNotFound)))
}

func (c Context) ClassicErrorHandler() internal.ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		var herr httperrors.HTTPError
		if !errors.As(err, &herr) {
			log.Println("Unexpected error:", err)
			herr = httperrors.New("Something went wrong", http.StatusInternalServerError)
		}
		statusCode, msg, _ := herr.HTTPError()

		// render error page
		w.WriteHeader(statusCode)
		c.Renderer().Render(w, classic_pages.Error{
			Title: msg,
		})
	}
}
