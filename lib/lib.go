package lib

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewConfiguration),
	fx.Provide(GetLogger),
)

func Render(w http.ResponseWriter, status int, template templ.Component) error {
	w.WriteHeader(status)
	return template.Render(context.Background(), w)
}

func IsHTMXRequest(r *http.Request) bool {
	hxRequest := r.Header.Get("HX-Request")

	if hxRequest == "true" {
		return true
	}

	return false
}
