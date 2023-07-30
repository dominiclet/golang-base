package user

import (
	"context"
	"errors"
	"time"

	"github.com/dominiclet/golang-base/lib/email"
	randgenerate "github.com/dominiclet/golang-base/lib/rand_generate"
	"github.com/dominiclet/golang-base/lib/resperror"
	"github.com/dominiclet/golang-base/lib/store"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid              string
	Name              string
	Email             string
	Password          string
	AccountType       AccountType
	LicenseExpiry     time.Time
	IsVerified        bool
	VerificationToken string
}

type UserService struct {
	db                  *gorm.DB
	logger              *logrus.Entry
	emailService        *email.EmailService
	resetPwTokens       *store.Store[string, string] // Maps email to generated reset token (reset tokens are tokens sent to email on reset request)
	resetPwAuthCodes    *store.Store[string, string] // Maps email to generate auth codes (auth codes are codes used to authorize a pw change API request)
	resendEmailDisabled *store.Store[uint, bool]     // Set of user IDs that cannot request verification email to be resent
}

func InitUserService(db *gorm.DB, emailService *email.EmailService) *UserService {
	return &UserService{
		db:                  db,
		logger:              logrus.WithField("module", "user_service"),
		emailService:        emailService,
		resetPwTokens:       store.NewStore[string, string](),
		resetPwAuthCodes:    store.NewStore[string, string](),
		resendEmailDisabled: store.NewStore[uint, bool](),
	}
}

// Get user by UUID
func (u *UserService) GetUserByUuid(ctx context.Context, uuid string) (*User, error) {
	var user User
	err := u.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		u.logger.WithField("err", err).Error("Query user DB error")
		return nil, err
	}
	return &user, nil
}

// Creates a user. Will hash provided password before storing into DB
func (u *UserService) CreateUser(ctx context.Context, name string,
	email string, password string) (*User, error) {
	// Check if email already exists
	err := u.db.Where("email = ?", email).First(&User{}).Error
	if err == nil {
		return nil, resperror.NewError(resperror.UserAlreadyExistsError)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	u.logger.WithFields(logrus.Fields{
		"name":  name,
		"email": email,
	}).Info("Creating user")

	// Hash password before storing
	hashedPassword, err := u.hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Generate UUID
	userUuid := uuid.NewString()

	licenseExpiry := time.Now().Add(time.Hour * 24 * TRIAL_DURATION_DAYS)

	newUser := User{
		Name:          name,
		Uuid:          userUuid,
		Email:         email,
		Password:      hashedPassword,
		AccountType:   TestAccount, // TODO: Change test account to trial account after beta test
		IsVerified:    false,       // Newly created user is unverified by default
		LicenseExpiry: licenseExpiry,
	}

	tx := u.db.Create(&newUser)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Send verification email (asynchronously)
	_, err = u.sendVerificationEmail(ctx, newUser.Email, newUser.Uuid, newUser.ID)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// Resend verification email (for users that already exist)
// Will return an error if user does not exist
// Idempotent within a certain period of time
func (u *UserService) ResendVerificationEmail(ctx context.Context, email string, password string) error {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	// Do not resend email if it is temporarily disabled
	if _, err := u.resendEmailDisabled.Get(user.ID); err == nil {
		return NewSendVerificationEmailErr(ResendDisabled,
			"Please wait for a period of time before trying to resend verification email")
	}
	// Check if user is already verified
	if user.IsVerified {
		return NewSendVerificationEmailErr(UserVerified,
			"User is already verified")
	}
	_, err = u.sendVerificationEmail(ctx, user.Email, user.Uuid, user.ID)
	if err != nil {
		return err
	}
	u.disableResendVerification(ctx, user.ID, time.Minute*verificationEmailDisableDuration)

	return nil
}

// Temporarily disable sending email verification for the indicated duration
func (u *UserService) disableResendVerification(ctx context.Context, userId uint, duration time.Duration) {
	u.logger.WithFields(logrus.Fields{
		"user_id":  userId,
		"duration": duration,
	}).Info("Disabling resending email verification for a period of time")
	u.resendEmailDisabled.Set(userId, true)
	go func() {
		time.Sleep(duration)
		u.resendEmailDisabled.Delete(userId)
	}()
}

// Generate verification token and send email
func (u *UserService) sendVerificationEmail(ctx context.Context, email string, userUUID string, userId uint) (string, error) {
	user := User{
		Model: gorm.Model{ID: userId},
	}

	// Generate verification token
	verificationToken, err := randgenerate.GenerateSecureToken(EmailVerificationTokenLength)
	if err != nil {
		u.logger.WithField("err", err).Error("Failed to generate email verification token")
		return "", err
	}

	// Save verification token in db
	err = u.db.Model(&user).Update("verification_token", verificationToken).Error
	if err != nil {
		u.logger.WithField("err", err).Error("Failed to store verification token in DB")
		return "", nil
	}

	// Send verification email
	go func() {
		err := u.emailService.SendVerificationEmail(email, userUUID, verificationToken)
		if err != nil {
			u.logger.WithFields(logrus.Fields{
				"err":   err,
				"email": email,
			}).Error("Failed to send verification email")
		}
	}()

	return verificationToken, nil
}

func (u *UserService) VerifyEmail(ctx context.Context, userUUID string, verificationToken string) error {
	var user User
	err := u.db.Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		u.logger.WithField("err", err).Error("Error querying for user by user UUID")
		return err
	}
	if user.IsVerified {
		return nil
	}
	if user.VerificationToken != verificationToken {
		u.logger.WithFields(logrus.Fields{
			"want": user.VerificationToken,
			"have": verificationToken,
		}).Error("Verification token mismatch")
		return errors.New("Verification token mismatch")
	}
	user.IsVerified = true
	user.VerificationToken = ""
	err = u.db.Save(&user).Error
	if err != nil {
		u.logger.WithField("err", err).Error("Error updating verified status")
		return err
	}
	u.resendEmailDisabled.Delete(user.ID)
	return nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		u.logger.WithField("err", err).Error("Failed to query user by email")
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetUserById(ctx context.Context, id uint) (*User, error) {
	var user User
	err := u.db.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logrus.WithField("err", err).Error("Failed to generate hash from password")
		return "", err
	}
	return string(hashedBytes), nil
}
