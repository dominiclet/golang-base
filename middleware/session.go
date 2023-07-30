package middleware

import (
	"github.com/dominiclet/golang-base/handler/session"
	ctxwrapper "github.com/dominiclet/golang-base/lib/ctx_wrapper"
	"github.com/dominiclet/golang-base/lib/httpresp"
	"github.com/dominiclet/golang-base/lib/resperror"
	"github.com/gin-gonic/gin"
)

// Check if user is authenticated (has a valid ongoing session)
func (m *Middleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(session.CookieKey)
		if err != nil {
			m.logger.Error("Session cookie not found")
			httpresp.SendError(c, resperror.NewError(resperror.Unauthorized))
			c.Abort()
			return
		}
		user, err := m.sessionService.GetSessionByToken(token)
		if err != nil {
			m.logger.WithField("err", err).Error("Failed to get session")
			httpresp.SendError(c, resperror.NewError(resperror.Unauthorized))
			c.Abort()
			return
		}
		// Inject user object into context
		ctxwrapper.SetUser(c, *user)
	}
}
