package initserver

import (
	"net/http"

	"github.com/dominiclet/golang-base/handler/session"
	"github.com/dominiclet/golang-base/handler/user"
	"github.com/dominiclet/golang-base/lib/httpresp"
	"github.com/dominiclet/golang-base/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type RouterService struct {
	middleware *middleware.Middleware

	userHandler    *user.UserHandler
	sessionHandler *session.SessionHandler
}

type Injector struct {
	middleware *middleware.Middleware

	userHandler    *user.UserHandler
	sessionHandler *session.SessionHandler
}

func InitRouterService(inj *Injector) *RouterService {
	return &RouterService{
		inj.middleware,
		inj.userHandler,
		inj.sessionHandler,
	}
}

var RouterSet = wire.NewSet(
	wire.Struct(new(Injector), "*"),
	InitRouterService,
)

func (rs *RouterService) RegisterRoutes(r *gin.Engine) {
	// All routes must start with /api
	apiGroup := r.Group("/api")

	apiGroup.GET("/ping", func(c *gin.Context) {
		httpresp.SendData(c, "pong", http.StatusOK)
	})

	rs.registerUsers(apiGroup)
	rs.registerSessions(apiGroup)
}

func (rs *RouterService) registerUsers(r *gin.RouterGroup) {
	userGroup := r.Group("/user")

	// Verification
	userGroup.GET("/verify/:userUuid/:token", rs.userHandler.VerifyEmail)
	userGroup.POST("/verify/resend_email", rs.userHandler.ResendVerificationEmail)

	userGroup.POST("", rs.userHandler.CreateUser)

	resetPw := userGroup.Group("/reset_password")
	resetPw.POST("", rs.userHandler.ResetPassword)
	resetPw.POST("/token_exchange", rs.userHandler.ResetPasswordAuthCodeExchange)
	resetPw.POST("/set_password", rs.userHandler.SetNewPassword)

	// Protected user endpoints
	protectedUserGroup := userGroup.Group("")
	protectedUserGroup.Use(rs.middleware.AuthRequired())
	protectedUserGroup.GET("/:uuid", rs.userHandler.GetUser)
}

func (rs *RouterService) registerSessions(r *gin.RouterGroup) {
	sessionGroup := r.Group("/session")

	sessionGroup.POST("login", rs.sessionHandler.UserLogin)
}
