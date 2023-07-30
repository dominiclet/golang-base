//go:build wireinject
// +build wireinject

package initserver

import (
	"github.com/dominiclet/golang-base/handler"
	"github.com/dominiclet/golang-base/init_server/config"
	"github.com/dominiclet/golang-base/init_server/env"
	"github.com/dominiclet/golang-base/lib"
	"github.com/dominiclet/golang-base/middleware"
	"github.com/dominiclet/golang-base/service"
	"github.com/google/wire"
)

func InitDeps() *RouterService {
	wire.Build(
		config.InitConfig,
		InitGormDB,
		RouterSet,
		env.InitEnvVars,

		middleware.MiddlewareSet,
		handler.HandlerSet,
		service.ServiceSet,
		lib.LibSet,
	)
	return &RouterService{}
}
