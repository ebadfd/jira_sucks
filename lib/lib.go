package lib

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewConfiguration),
	fx.Provide(GetLogger),
)
