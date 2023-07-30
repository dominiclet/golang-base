package user

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/dominiclet/golang-base/lib/httpresp"
	"github.com/dominiclet/golang-base/lib/resperror"
	"github.com/dominiclet/golang-base/service/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *user.UserService
}

func InitUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{
		userService,
	}
}

// @Summary Create user
// @Description Create new user
// @Tags user
// @Accept json
// @Param req body CreateUserRequest true "Create user data"
// @Produce json
// @Failure 409 {object} httpresp.StandardResponse "User with same email already exists"
// @Success 200 {object} httpresp.StandardDataResponse{data=CreateUserResponse}
// @Router /user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var userReq CreateUserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}

	newUser, err := h.userService.CreateUser(c, userReq.Name,
		userReq.Email, userReq.Password)
	if err != nil {
		httpresp.SendError(c, err)
		return
	}

	httpresp.SendData(c, CreateUserResponse{
		ID: newUser.ID,
	}, http.StatusOK)
}

// @Summary Get basic user information
// @Description Get basic user information (protected endpoint)
// @Tags user,authRequired
// @Param uuid query string true "UUID"
// @Success 200 {object} httpresp.StandardDataResponse{data=User}
// @Failure 404 {object} httpresp.StandardResponse "User not found"
// @Router /user/{uuid} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	uuid := c.Param("uuid")

	decodedUUID, err := url.QueryUnescape(uuid)
	if err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}
	user, err := h.userService.GetUserByUuid(c, decodedUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpresp.SendError(c, resperror.NewError(resperror.UserNotFound))
			return
		}
		httpresp.SendErrorWithStatusCode(c, err, http.StatusInternalServerError)
		return
	}
	respUser := NewUserFromSvcUser(user)

	httpresp.SendData(c, respUser, http.StatusOK)
}

// @Summary Verify email
// @Description Handles verification link for email
// @Tags user
// @Param userUuid query int true "User ID"
// @Param token query string true "Verification token"
// @Router /user/verify/{userUuid}/{token} [get]
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	userUUID := c.Param("userUuid")
	verificationToken := c.Param("token")

	decodedUUID, err := url.QueryUnescape(userUUID)
	if err != nil {
		c.String(http.StatusBadRequest, "Verification failed")
		return
	}
	err = h.userService.VerifyEmail(c, decodedUUID, verificationToken)
	if err != nil {
		c.String(http.StatusForbidden, "Verification failed")
		return
	}
	c.String(http.StatusOK, "Verification successful. Please head back to the portal to login.")
}

// @Summary Resned verification email
// @Description Resends verification email while invalidating previous verification link
// @Tags user
// @Param req body ResendVerificationEmailRequest true "Email and password must be provided to authenticate before sending"
// @Accept json
// @Produce json
// @Success 200 {object} httpresp.StandardResponse
// @Failure 429 {object} httpresp.StandardResponse "Too many requests made in a short period of time (2 mins)"
// @Failure 405 {object} httpresp.StandardResponse "User already verified"
// @Router /user/verify/resend_email [post]
func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}
	err := h.userService.ResendVerificationEmail(c, req.Email, req.Password)
	if err != nil {
		if expectedErr, ok := err.(user.SendVerificationEmailError); ok {
			switch expectedErr.ErrType() {
			case user.ResendDisabled:
				httpresp.SendError(c, resperror.NewError(resperror.TooManyRequests))
				return
			case user.UserVerified:
				httpresp.SendError(c, resperror.NewError(resperror.UserAlreadyVerifiedError))
				return
			}
		}
		// Just return 500 internal server error by default
		httpresp.SendError(c, err)
		return
	}
	httpresp.SendSuccess(c)
}

// @Summary Reset password
// @Description Starts reset password process by generating and sending 6-digit token to provided email if account exists
// @Tags user
// @Accept json
// @Param req body ResetPasswordRequest true "Specify email of account to reset"
// @Produce json
// @Success 200 {object} httpresp.StandardResponse
// @Router /user/reset_password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}

	// Error is intentionally not handled
	h.userService.ResetPassword(c, req.Email)

	// Handler always returns success to mask any errors from FE
	httpresp.SendSuccess(c)
}

// @Summary Reset password auth code exchange
// @Description Verify 6-digit token and exchange token with auth code to set new password
// @Tags user
// @Accept json
// @Param req body ResetPWAuthCodeExchangeRequest true "Specify both email and token"
// @Produce json
// @Failure 401 {object} httpresp.StandardResponse "Token rejected"
// @Success 200 {object} httpresp.StandardDataResponse{data=ResetPWAuthCodeExchangeResponse}
// @Router /user/reset_password/token_exchange [post]
func (h *UserHandler) ResetPasswordAuthCodeExchange(c *gin.Context) {
	var req ResetPWAuthCodeExchangeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}

	authCode, err := h.userService.ResetPwAuthCodeExchange(c, req.Token, req.Email)
	if err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.InvalidTokenError))
		return
	}

	httpresp.SendData(c, ResetPWAuthCodeExchangeResponse{
		AuthCode: authCode,
	}, http.StatusOK)
}

// Final API call for reset password process (sets new password)
// @Summary Set new password
// @Description Final API call for reset password process. Set new password, authenticating using auth code.
// @Tags user
// @Accept json
// @Param req body SetNewPWRequest true "Specify new password for account"
// @Produce json
// @Failure 401 {object} httpresp.StandardResponse "Auth code rejected"
// @Success 200 {object} httpresp.StandardResponse
// @Router /user/reset_password/set_password [post]
func (h *UserHandler) SetNewPassword(c *gin.Context) {
	var req SetNewPWRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.BadRequest))
		return
	}

	err := h.userService.SetNewPassword(c, req.Email, req.AuthCode, req.NewPassword)
	if err != nil {
		httpresp.SendError(c, resperror.NewError(resperror.Unauthorized))
		return
	}
	httpresp.SendSuccess(c)
}
