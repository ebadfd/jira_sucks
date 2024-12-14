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
