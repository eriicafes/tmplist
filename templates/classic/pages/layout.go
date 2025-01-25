package classic_pages

import "github.com/eriicafes/tmplist/db"

type toastType string

var (
	ToastWarning = toastType("warning")
	ToastError   = toastType("error")
	ToastSuccess = toastType("success")
)

type Toast struct {
	Message string
	Type    toastType
}

type Layout struct {
	Toast Toast
	Title string
	User  *db.User
}
