package request

import (
	"context"
	"fmt"

	"github.com/eriicafes/tmplist/db"
)

type ContextValue[T any] string

var (
	User    ContextValue[db.User]    = "user"
	Session ContextValue[db.Session] = "session"
)

func (cv ContextValue[T]) SetContext(ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, cv, val)
}

func (cv ContextValue[T]) FromContext(ctx context.Context) (T, bool) {
	val := ctx.Value(cv)
	fmt.Println("ctx val:", val)
	v, ok := val.(T)
	fmt.Println("cast val:", v, ok)
	return v, ok
}
