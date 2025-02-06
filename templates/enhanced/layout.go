package enhanced

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
)

type toastType string

var (
	ToastWarning = toastType("warning")
	ToastError   = toastType("error")
	ToastSuccess = toastType("success")
)

type Toast struct {
	Message string
	Type    toastType
	Swap    bool
}

func (t Toast) Tmpl() tmpl.Template {
	return tmpl.Tmpl("components/toast", tmpl.Map{
		"message": t.Message,
		"type":    t.Type,
		"swap":    t.Swap,
	})
}

type Layout struct {
	tmpl.Children
	Toast Toast
	Title string
	User  *db.User
}

func (l Layout) Tmpl() tmpl.Template {
	return tmpl.Associated(l.Base(), "enhanced/layout", l)
}
