package middleware

import (
	"github.com/dominiclet/golang-base/init_server/logger"
	"github.com/dominiclet/golang-base/service/session"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	sessionService *session.SessionService
	logger         *logrus.Entry
}

func InitMiddleware(sessionService *session.SessionService) *Middleware {
	return &Middleware{
		sessionService: sessionService,
		logger:         logger.GetLogger().WithField("module", "middleware"),
	}
}
