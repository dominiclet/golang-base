package user

import "github.com/dominiclet/golang-base/service/user"

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	ID uint `json:"id"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPWAuthCodeExchangeRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResetPWAuthCodeExchangeResponse struct {
	AuthCode string `json:"auth_code"`
}

type SetNewPWRequest struct {
	Email       string `json:"email"`
	AuthCode    string `json:"auth_code"`
	NewPassword string `json:"new_password"`
}

type ResendVerificationEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccountType int    `json:"account_type"`
	IsVerified  bool   `json:"is_verified"`
}

func NewUserFromSvcUser(svcUser *user.User) User {
	return User{
		Uuid:        svcUser.Uuid,
		Name:        svcUser.Name,
		Email:       svcUser.Email,
		AccountType: int(svcUser.AccountType),
		IsVerified:  svcUser.IsVerified,
	}
}
