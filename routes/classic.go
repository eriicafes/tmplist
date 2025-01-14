package routes

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/internal"
	"github.com/eriicafes/tmplist/internal/httperrors"
	"github.com/eriicafes/tmplist/internal/session"
	"github.com/eriicafes/tmplist/schemas"
	classic_pages "github.com/eriicafes/tmplist/templates/classic/pages"
)

var flashMessage session.Flash[classic_pages.FlashMessage] = "flash_message"

func (c Context) Classic(mux internal.Mux) {
	mux = internal.Fallback(mux, c.ClassicErrorHandler())
	auth := internal.Apply(mux, c.authMiddleware())
	guest := internal.Apply(mux, c.guestMiddleware())

	// all topics page
	auth.Route("GET /{$}", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		user, _ := requestUser.Get(r.Context())
		flash, _ := flashMessage.Get(w, r)

		return tr.Render(w, classic_pages.Index{
			Layout: classic_pages.Layout{Flash: flash},
			User:   user,
			Topics: nil,
		})
	})

	// all todos for topic page
	auth.Route("GET /{topicId}", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		topicId, _ := strconv.Atoi(r.PathValue("topicId"))
		topic, err := c.DB.GetTopic(topicId)
		if err != nil {
			return httperrors.New("Topic not found", http.StatusNotFound)
		}
		flash, _ := flashMessage.Get(w, r)

		return tr.Render(w, classic_pages.Topic{
			Layout: classic_pages.Layout{Flash: flash},
			Topic:  topic,
			Todos:  nil,
		})
	})

	// login page
	guest.Route("GET /login", func(w http.ResponseWriter, r *http.Request) error {
		return c.Renderer().Render(w, classic_pages.Login{})
	})

	// login post
	guest.Route("POST /login", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		// validate input
		form := schemas.LoginData{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		errors := schemas.FormErrors(form.Validate())
		if errors != nil {
			return tr.Render(w, classic_pages.Login{
				Form:   form,
				Errors: errors,
			})
		}

		// get user from db
		user, err := c.DB.GetUserByEmail(form.Email)
		if err != nil {
			return tr.Render(w, classic_pages.Login{
				Form:  form,
				Error: "Email address not found",
			})
		}
		// check password
		if !c.Auth.ComparePassword(user.PasswordHash, form.Password) {
			return tr.Render(w, classic_pages.Login{
				Form:  form,
				Error: "Invalid password",
			})
		}
		// create session and set cookie
		token, err := c.Auth.GenerateSessionToken()
		if err != nil {
			return tr.Render(w, classic_pages.Login{
				Form:  form,
				Error: "Failed to login",
			})
		}
		session, err := c.Auth.CreateSession(token, user)
		if err != nil {
			return tr.Render(w, classic_pages.Login{
				Form:  form,
				Error: "Failed to login",
			})
		}
		c.Auth.SetCookie(w, token, session.ExpiresAt)

		// set flash message
		flashMessage.Set(w, classic_pages.FlashMessage{
			Message: "Logged in successfully",
			Success: true,
		})

		// redirect to authenticated page
		http.Redirect(w, r, "/classic/", http.StatusFound)
		return nil
	})

	// register page
	guest.Route("GET /register", func(w http.ResponseWriter, r *http.Request) error {
		return c.Renderer().Render(w, classic_pages.Register{})
	})

	// register post
	guest.Route("POST /register", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		// validate input
		form := schemas.RegisterData{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		errors := schemas.FormErrors(form.Validate())
		if errors != nil {
			return tr.Render(w, classic_pages.Register{
				Form:   form,
				Errors: errors,
			})
		}

		// hash password
		passwordHash, err := c.Auth.HashPassword(form.Password)
		if err != nil {
			return tr.Render(w, classic_pages.Register{
				Form:  form,
				Error: "Failed to create account",
			})
		}
		// insert user in db
		user, err := c.DB.InsertUser(form.Email, passwordHash)
		if err != nil {
			msg := "Failed to create account"
			if err == db.ErrDuplicate {
				msg = "Email address already taken"
			}
			return tr.Render(w, classic_pages.Register{
				Form:  form,
				Error: msg,
			})
		}
		// create session and set cookie
		token, err := c.Auth.GenerateSessionToken()
		if err != nil {
			return tr.Render(w, classic_pages.Register{
				Form:  form,
				Error: "Failed to create account",
			})
		}
		session, err := c.Auth.CreateSession(token, user)
		if err != nil {
			return tr.Render(w, classic_pages.Register{
				Form:  form,
				Error: "Failed to create account",
			})
		}
		c.Auth.SetCookie(w, token, session.ExpiresAt)

		// set flash message
		flashMessage.Set(w, classic_pages.FlashMessage{
			Message: "Account created successfully",
			Success: true,
		})

		// redirect to authenticated page
		http.Redirect(w, r, "/classic/", http.StatusFound)
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
		c.Renderer().Render(w, tmpl.Tmpl("classic/pages/error", tmpl.Map{
			"Title": msg,
		}))
	}
}
