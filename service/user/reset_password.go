package user

import (
	"context"
	"errors"
	"time"

	randgenerate "github.com/dominiclet/golang-base/lib/rand_generate"
	"github.com/sirupsen/logrus"
)

const (
	tokenLength      = 6
	authCodeLength   = 16
	tokenValidity    = 5 // No. of minutes where reset pw token is valid
	authCodeValidity = 2 // No. of minutes where reset pw auth code is valid
)

// Sends initial request to reset password for user with email
// Sends an email to the user with reset token if account with email exists
func (u *UserService) ResetPassword(ctx context.Context, email string) error {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	token := randgenerate.GenerateAlphaNumericString(tokenLength)
	u.resetPwTokens.Set(user.Email, token)
	go func() {
		// Remove token after token validity period has passed
		time.Sleep(time.Minute * tokenValidity)
		u.resetPwTokens.Delete(user.Email)
	}()

	// Send reset password token asynchronously
	go func() {
		err = u.emailService.SendResetPasswordToken(user.Email, token)
		if err != nil {
			u.logger.WithField("err", err).Error("Failed to send reset password email")
		}
	}()

	return nil
}

func (u *UserService) ResetPwAuthCodeExchange(ctx context.Context, token string, email string) (string, error) {
	u.logger.WithFields(logrus.Fields{
		"token": token,
		"email": email,
	}).Info("Exchanging reset pw token for auth code")
	storedToken, err := u.resetPwTokens.Get(email)
	if err != nil {
		return "", errors.New("No token stored for provided email")
	}

	// Check if provided token matches stored token
	if storedToken != token {
		u.logger.WithFields(logrus.Fields{
			"want": storedToken,
			"got":  token,
		}).Error("Token mismatch")
		return "", errors.New("Token mismatch")
	}

	u.resetPwTokens.Delete(email)

	authCode, err := randgenerate.GenerateSecureToken(authCodeLength)
	if err != nil {
		u.logger.WithField("err", err).Error("Failed to generate auth code")
		return "", err
	}

	u.resetPwAuthCodes.Set(email, authCode)
	go func() {
		// Remove auth code after auth code validity period has passed
		time.Sleep(time.Minute * authCodeValidity)
		u.resetPwAuthCodes.Delete(email)
	}()

	return authCode, nil
}

func (u *UserService) SetNewPassword(ctx context.Context, email string, authCode string, newPassword string) error {
	u.logger.WithField("email", email).Info("Setting new password")

	storedAuthCode, err := u.resetPwAuthCodes.Get(email)
	if err != nil {
		return errors.New("No auth code stored for provided email")
	}

	// Check if provided auth code matches stored auth code
	if storedAuthCode != authCode {
		u.logger.Error("Auth code mismatch")
		return errors.New("Auth code mismatch")
	}

	u.resetPwAuthCodes.Delete(email)

	// Update user with new password
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	hashedPassword, err := u.hashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return u.db.Model(user).Select("password").Updates(*user).Error
}
