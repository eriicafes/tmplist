package services

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type Flash[T any] string

func (s Flash[T]) Set(w http.ResponseWriter, data T) {
	jsonData, _ := json.Marshal(data)
	cookieData := base64.URLEncoding.EncodeToString(jsonData)

	// set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  string(s),
		Value: cookieData,
	})
}

func (s Flash[T]) Get(w http.ResponseWriter, r *http.Request) (T, bool) {
	var data T

	cookieData, err := r.Cookie(string(s))
	if err != nil {
		return data, false
	}
	jsonString, err := base64.URLEncoding.DecodeString(cookieData.Value)
	if err != nil {
		return data, false
	}
	if err := json.Unmarshal(jsonString, &data); err != nil {
		return data, false
	}

	// clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   string(s),
		MaxAge: -1,
	})
	return data, true
}
