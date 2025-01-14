package session

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// Flash sets and gets flash messages with session cookies (cookies that are deleted once the browser session ends).
type Flash[T any] string

// Set writes a cookies flash message.
func (s Flash[T]) Set(w http.ResponseWriter, data T) {
	// marshal json and encode with base64
	jsonData, _ := json.Marshal(data)
	cookieData := base64.URLEncoding.EncodeToString(jsonData)
	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  string(s),
		Value: cookieData,
	})
}

// Get reads a cookies flash message and deletes it immediately.
func (s Flash[T]) Get(w http.ResponseWriter, r *http.Request) (T, bool) {
	var data T
	cookieData, err := r.Cookie(string(s))
	if err != nil {
		return data, false
	}
	// decode base64 and unmarshal json
	jsonString, err := base64.URLEncoding.DecodeString(cookieData.Value)
	if err != nil {
		// delete malformed cookie
		http.SetCookie(w, &http.Cookie{
			Name:   string(s),
			MaxAge: -1,
		})
		return data, false
	}
	if err := json.Unmarshal(jsonString, &data); err != nil {
		// delete malformed cookie
		http.SetCookie(w, &http.Cookie{
			Name:   string(s),
			MaxAge: -1,
		})
		return data, false
	}
	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   string(s),
		MaxAge: -1,
	})
	return data, true
}
