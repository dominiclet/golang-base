package session

import (
	"net/http"

	"github.com/dominiclet/golang-base/init_server/config"
	"github.com/dominiclet/golang-base/init_server/env"
	"github.com/dominiclet/golang-base/init_server/logger"
	"github.com/dominiclet/golang-base/lib/httpresp"
	"github.com/dominiclet/golang-base/lib/resperror"
	"github.com/dominiclet/golang-base/service/session"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	sessionService *session.SessionService
	config         *config.Config
	logger         *logrus.Entry
	envVars        *env.EnvVars
}

func InitSessionHandler(sessionService *session.SessionService, config *config.Config, envVars *env.EnvVars) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		config:         config,
		logger:         logger.GetLogger().WithField("module", "session_handler"),
		envVars:        envVars,
	}
}

// @Summary User login
// @Description Create login session for user
// @Tags session
// @Accept json
// @Param req body UserLoginRequest true "Email and password for authentication"
// @Produce json
// @Failure 403 {object} httpresp.StandardResponse "User is not verified"
// @Failure 401 {object} httpresp.StandardResponse "Authentication failed"
// @Success 200 {object} httpresp.StandardDataResponse{data=UserLoginResponse}
// @Router /session/login [post]
func (s *SessionHandler) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}

	user, token, expiry, err := s.sessionService.CreateUserSession(c, req.Email, req.Password)
	if err != nil {
		s.logger.WithField("err", err).Error("Error occurred while creating user session")
		httpresp.SendErrorWithFallback(c, err, resperror.NewError(resperror.Unauthorized))
		return
	}

	var secureCookie bool
	if !s.envVars.IsDev() {
		secureCookie = true
	}
	c.SetCookie(CookieKey, token, daySeconds*session.DefaultSessionDurationDay, "/", s.config.Domain, secureCookie, true)
	httpresp.SendData(c, UserLoginResponse{Uuid: user.Uuid, Expiry: expiry}, http.StatusOK)
}
