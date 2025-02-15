package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/schemas"
)

func (c Context) Api(mux internal.Mux) {
	mux = internal.Fallback(mux, c.ApiErrorHandler())
	auth := internal.Use(mux, c.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		c.JSONStatus(w, http.StatusUnauthorized, ApiError{Message: "Unauthorized"})
	}))

	// all topics
	auth.Route("GET ", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		search := r.URL.Query().Get("search")

		// get topics from db
		topics, _ := c.DB.GetTopics(user.Id, search)

		return c.JSON(w, topics)
	})

	// create topic
	auth.Route("POST ", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())

		// validate input
		var data schemas.TopicData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
		// return form validation errors
		if details := schemas.FormErrors(data.Validate()); details != nil {
			return httperrors.NewDetails("Invalid input", http.StatusBadRequest, details)
		}

		// insert topic in db
		topic, err := c.DB.InsertTopic(user.Id, data.Topic)
		if err != nil {
			log.Println(err)
			return httperrors.New("Failed to create topic", http.StatusInternalServerError)
		}
		if len(data.Todos) > 0 {
			var insertTodos []db.Todo
			for _, v := range slices.Backward(data.Todos) {
				insertTodos = append(insertTodos, db.Todo{
					TopicId: topic.Id,
					Body:    v.Text,
					Done:    v.Checked,
				})
			}
			// insert topic todos in db
			_, err = c.DB.InsertTodos(insertTodos)
			if err != nil {
				log.Println(err)
				return httperrors.New("Failed to create todos", http.StatusInternalServerError)
			}
		}

		topic.TodosCount = len(data.Todos)
		return c.JSON(w, topic)
	})

	var requestTopic internal.ContextValue[db.Topic] = "topic"
	withTopic := internal.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())
		topicId, _ := strconv.Atoi(r.PathValue("topicId"))

		// check if topic exists and belongs to user
		topic, err := c.DB.GetTopic(topicId)
		if err != nil || topic.UserId != user.Id {
			return httperrors.New("Topic not found", http.StatusNotFound)
		}

		ctx := requestTopic.Set(r.Context(), topic)
		return internal.WithRequest(r.WithContext(ctx))
	})

	// get topic
	auth.Route("GET /{topicId}", withTopic, func(w http.ResponseWriter, r *http.Request) error {
		topic, _ := requestTopic.Get(r.Context())

		// get todos for topic
		todos, err := c.DB.GetTodos(topic.Id)
		if err != nil {
			log.Println(err)
		}

		return c.JSON(w, map[string]any{
			"topic": topic,
			"todos": todos,
		})
	})

	// update topic
	auth.Route("PUT /{topicId}", withTopic, func(w http.ResponseWriter, r *http.Request) error {
		topic, _ := requestTopic.Get(r.Context())

		// validate input
		var data schemas.TopicData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
		// return form validation errors
		if details := schemas.FormErrors(data.Validate()); details != nil {
			return httperrors.NewDetails("Invalid input", http.StatusBadRequest, details)
		}

		// update todo in db
		topic, err := c.DB.UpdateTopic(topic.Id, data.Topic)
		if err != nil {
			log.Println(err)
			return httperrors.New("Failed to update topic", http.StatusInternalServerError)
		}

		return c.JSON(w, topic)
	})

	// delete topic
	auth.Route("DELETE /{topicId}", withTopic, func(w http.ResponseWriter, r *http.Request) error {
		topic, _ := requestTopic.Get(r.Context())

		// delete todo from db
		if err := c.DB.DeleteTopic(topic.Id); err != nil {
			log.Println(err)
			return httperrors.New("Failed to delete topic", http.StatusInternalServerError)
		}

		return c.JSON(w, topic)
	})

	// create todo
	auth.Route("POST /{topicId}/todos", withTopic, func(w http.ResponseWriter, r *http.Request) error {
		topic, _ := requestTopic.Get(r.Context())

		// validate input
		var data schemas.TodoData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
		// return form validation errors
		if details := schemas.FormErrors(data.Validate()); details != nil {
			return httperrors.NewDetails("Invalid input", http.StatusBadRequest, details)
		}

		// insert todo in db
		insertTodos := []db.Todo{{TopicId: topic.Id, Body: data.Text}}
		if _, err := c.DB.InsertTodos(insertTodos); err != nil {
			log.Println(err)
			return httperrors.New("Failed to create todo", http.StatusInternalServerError)
		}

		// get todos for topic
		todos, err := c.DB.GetTodos(topic.Id)
		if err != nil {
			log.Println(err)
		}

		return c.JSON(w, todos)
	})

	var requestTodo internal.ContextValue[db.Todo] = "todo"
	withTodo := internal.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		topic, _ := requestTopic.Get(r.Context())
		todoId, _ := strconv.Atoi(r.PathValue("todoId"))

		// check if todo exists and belongs to topic
		todo, err := c.DB.GetTodo(todoId)
		if err != nil || todo.TopicId != topic.Id {
			return httperrors.New("Todo not found", http.StatusNotFound)
		}

		ctx := requestTodo.Set(r.Context(), todo)
		return internal.WithRequest(r.WithContext(ctx))
	})

	// update todo
	auth.Route("PUT /{topicId}/todos/{todoId}", withTopic, withTodo, func(w http.ResponseWriter, r *http.Request) error {
		todo, _ := requestTodo.Get(r.Context())

		// validate input
		var data schemas.TodoData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
		// return form validation errors
		if details := schemas.FormErrors(data.Validate()); details != nil {
			return httperrors.NewDetails("Invalid input", http.StatusBadRequest, details)
		}

		// update todo in db
		todo, err := c.DB.UpdateTodo(todo.Id, data.Text, data.Checked)
		if err != nil {
			log.Println(err)
			return httperrors.New("Failed to update todo", http.StatusInternalServerError)
		}

		return c.JSON(w, todo)
	})

	// delete todo
	auth.Route("DELETE /{topicId}/todos/{todoId}", withTopic, withTodo, func(w http.ResponseWriter, r *http.Request) error {
		todo, _ := requestTodo.Get(r.Context())

		// delete todo from db
		if err := c.DB.DeleteTodo(todo.Id); err != nil {
			log.Println(err)
			return httperrors.New("Failed to delete todo", http.StatusInternalServerError)
		}

		return c.JSON(w, todo)
	})

	// profile
	auth.Route("GET /profile", func(w http.ResponseWriter, r *http.Request) error {
		user, _ := requestUser.Get(r.Context())

		return c.JSON(w, user)
	})

	// login
	mux.Route("POST /login", func(w http.ResponseWriter, r *http.Request) error {
		// prevent other origins from authenticating the user
		if !c.allowOriginForNonSafeRequests(r) {
			return httperrors.New("Login attempt from an unknown site blocked", http.StatusForbidden)
		}

		// validate input
		var data schemas.LoginData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
		// return form validation errors
		if details := schemas.FormErrors(data.Validate()); details != nil {
			return httperrors.NewDetails("", http.StatusBadRequest, details)
		}

		// get user from db
		user, err := c.DB.GetUserByEmail(data.Email)
		if err != nil {
			return httperrors.New("Email address not found", http.StatusBadRequest)
		}
		// check password
		if !c.Auth.ComparePassword(user.PasswordHash, data.Password) {
			return httperrors.New("Invalid password", http.StatusBadRequest)
		}
		// create session and set cookie
		token, err := c.Auth.GenerateSessionToken()
		if err != nil {
			return httperrors.New("Failed to login", http.StatusInternalServerError)
		}
		session, err := c.Auth.CreateSession(token, user)
		if err != nil {
			return httperrors.New("Failed to login", http.StatusInternalServerError)
		}
		c.Auth.SetCookie(w, token, session.ExpiresAt)

		return c.JSON(w, user)
	})

	// register
	mux.Route("POST /register", func(w http.ResponseWriter, r *http.Request) error {
		// prevent other origins from authenticating the user
		if !c.allowOriginForNonSafeRequests(r) {
			return httperrors.New("Register attempt from an unknown site blocked", http.StatusForbidden)
		}

		// validate input
		var data schemas.RegisterData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
		// return form validation errors
		if details := schemas.FormErrors(data.Validate()); details != nil {
			return httperrors.NewDetails("", http.StatusBadRequest, details)
		}

		// hash password
		passwordHash, err := c.Auth.HashPassword(data.Password)
		if err != nil {
			return httperrors.New("Failed to create account", http.StatusInternalServerError)
		}
		// insert user in db
		user, err := c.DB.InsertUser(data.Email, passwordHash)
		if err != nil {
			if err == db.ErrDuplicate {
				return httperrors.New("Email address already taken", http.StatusConflict)
			}
			return httperrors.New("Failed to create account", http.StatusInternalServerError)
		}
		// create session and set cookie
		token, err := c.Auth.GenerateSessionToken()
		if err != nil {
			return httperrors.New("Failed to create account", http.StatusInternalServerError)
		}
		session, err := c.Auth.CreateSession(token, user)
		if err != nil {
			return httperrors.New("Failed to create account", http.StatusInternalServerError)
		}
		c.Auth.SetCookie(w, token, session.ExpiresAt)

		return c.JSON(w, map[string]any{
			"message": "Account created successfully",
			"profile": user,
		})
	})

	// logout
	auth.Route("POST /logout", func(w http.ResponseWriter, r *http.Request) error {
		session, _ := requestSession.Get(r.Context())
		if err := c.Auth.InvalidateSession(session.Id); err != nil {
			return err
		}
		c.Auth.DeleteCookie(w)

		return c.JSON(w, map[string]any{
			"message": "Logged out successfully",
		})
	})
	mux.HandleFunc("/", internal.ErrorHandler(mux, httperrors.New("route not found", http.StatusNotFound)))
}

type ApiError struct {
	Message string             `json:"message"`
	Errors  httperrors.Details `json:"errors"`
}

func (c Context) ApiErrorHandler() internal.ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		var herr httperrors.HTTPError
		if !errors.As(err, &herr) {
			log.Println("Unexpected error:", err)
			herr = httperrors.New("Something went wrong", http.StatusInternalServerError)
		}
		statusCode, msg, details := herr.HTTPError()
		c.JSONStatus(w, statusCode, ApiError{msg, details})
	}
}

func (c Context) JSON(w http.ResponseWriter, data any) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func (c Context) JSONStatus(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
