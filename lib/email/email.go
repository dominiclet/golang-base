package email

import (
	"fmt"
	"net/url"

	"github.com/dominiclet/golang-base/init_server/config"
	"github.com/dominiclet/golang-base/init_server/env"
	"github.com/dominiclet/golang-base/init_server/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

const (
	emailDisplayName = "Fund Analysis"
)

type EmailService struct {
	config *config.Config
	env    *env.EnvVars
	dialer *gomail.Dialer
	logger *logrus.Entry
	from   string
}

func InitEmailService(config *config.Config, env *env.EnvVars) *EmailService {
	return &EmailService{
		config: config,
		dialer: gomail.NewDialer(config.Email.ServerAddress, config.Email.ServerPort,
			config.Email.EmailAddress, config.Email.AppPassword),
		logger: logger.GetLogger().WithField("module", "email_service"),
		from:   fmt.Sprintf("%s <%s>", emailDisplayName, config.Email.EmailAddress),
		env:    env,
	}
}

func (e *EmailService) SendResetPasswordToken(to string, token string) error {
	e.logger.WithFields(logrus.Fields{
		"to":    to,
		"token": token,
	}).Info("Sending reset password email")
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset Password")

	content := fmt.Sprintf("Your reset token is <b>%s</b>", token)

	m.SetBody("text/html", content)

	if err := e.dialer.DialAndSend(m); err != nil {
		e.logger.WithField("err", err).Error("Error occurred when sending email")
		return err
	}
	return nil
}

func (e *EmailService) SendVerificationEmail(to string, userUUID string, verificationToken string) error {
	e.logger.WithFields(logrus.Fields{
		"to":                to,
		"verificationToken": verificationToken,
	}).Info("Sending verification email")
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify Email")

	protocol := e.env.GetHttpProtocol()
	escapedUUID := url.QueryEscape(userUUID)
	verificationLink := fmt.Sprintf("%s://%s/api/user/verify/%s/%s",
		protocol, e.config.Domain, escapedUUID, verificationToken)
	content := fmt.Sprintf(`Please click <a href="%s">here</a> to verify your email`, verificationLink)

	m.SetBody("text/html", content)
	if err := e.dialer.DialAndSend(m); err != nil {
		e.logger.WithField("err", err).Error("Error occurred when sending email")
		return err
	}
	return nil
}
