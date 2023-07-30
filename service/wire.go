package service

import (
	"github.com/dominiclet/golang-base/service/session"
	"github.com/dominiclet/golang-base/service/user"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(user.InitUserService, session.InitSessionService)
