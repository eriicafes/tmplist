package routes

import (
	"net/http"

	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
	"github.com/eriicafes/tmplist/request"
	"github.com/eriicafes/tmplist/schemas"
	"github.com/eriicafes/tmplist/services"
	classic_pages "github.com/eriicafes/tmplist/templates/classic/pages"
)

var authReponseFlash services.Flash[string] = "auth_response"

func (c Context) MountClassic(mux Mux) {
	auth := Route(mux, c.authMiddleware())
	guest := Route(mux, c.guestMiddleware())

	// all topics page
	auth.On("GET /{$}", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()

		user, _ := request.User.FromContext(r.Context())
		flash, _ := authReponseFlash.Get(w, r)

		return tr.Render(w, tmpl.Tmpl("classic/pages/index", tmpl.Map{
			"User":  user,
			"Flash": flash,
		}))
	})

	// all todos for topic page
	auth.On("GET /{topicId}", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("classic/pages/topic"))
	})

	// login page
	guest.On("GET /login", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("classic/pages/login"))
	})

	// login post
	guest.On("POST /login", func(w http.ResponseWriter, r *http.Request) error {
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
		session, err := c.Auth.CreateSession(token, user.Id)
		if err != nil {
			return tr.Render(w, classic_pages.Login{
				Form:  form,
				Error: "Failed to login",
			})
		}
		c.Session.SetCookie(w, token, session.ExpiresAt)

		// set flash message
		authReponseFlash.Set(w, "Logged in successfully")

		// redirect to authenticated page
		http.Redirect(w, r, "/classic/", http.StatusFound)
		return nil
	})

	// register page
	guest.On("GET /register", func(w http.ResponseWriter, r *http.Request) error {
		tr := c.Renderer()
		return tr.Render(w, tmpl.Tmpl("classic/pages/register"))
	})

	// register post
	guest.On("POST /register", func(w http.ResponseWriter, r *http.Request) error {
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
		session, err := c.Auth.CreateSession(token, user.Id)
		if err != nil {
			return tr.Render(w, classic_pages.Register{
				Form:  form,
				Error: "Failed to create account",
			})
		}
		c.Session.SetCookie(w, token, session.ExpiresAt)

		// set flash message
		authReponseFlash.Set(w, "Account created successfully")

		// redirect to authenticated page
		http.Redirect(w, r, "/classic/", http.StatusFound)
		return nil
	})

	// 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.Renderer().Render(w, tmpl.Tmpl("classic/pages/404"))
	})
}
