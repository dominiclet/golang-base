package handler

import (
	"github.com/dominiclet/golang-base/handler/session"
	"github.com/dominiclet/golang-base/handler/user"
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(user.InitUserHandler, session.InitSessionHandler)
