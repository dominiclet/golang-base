package session

import (
	"context"
	"errors"
	"time"

	"github.com/dominiclet/golang-base/init_server/logger"
	"github.com/dominiclet/golang-base/lib/resperror"
	"github.com/dominiclet/golang-base/lib/store"
	"github.com/dominiclet/golang-base/service/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Session struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time

	User user.User
}

type SessionService struct {
	userService *user.UserService
	db          *gorm.DB
	logger      *logrus.Entry
	// sessionCache maps session tokens to the Session object it is associated with
	// for faster validation of session token
	sessionCache *store.Store[string, Session]
}

func InitSessionService(userService *user.UserService, db *gorm.DB) *SessionService {
	return &SessionService{
		userService:  userService,
		db:           db,
		sessionCache: store.NewStore[string, Session](),
		logger:       logger.GetLogger().WithField("module", "session_service"),
	}
}

// Verifies user email and password, then generates session for user, returning the user object, session token, and time of expiry
func (a *SessionService) CreateUserSession(ctx context.Context, email string, password string) (*user.User, string, int64, error) {
	user, err := a.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", 0, resperror.NewError(resperror.UserEmailNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		a.logger.WithField("err", err).Error("Incorrect password given")
		return nil, "", 0, resperror.NewError(resperror.UserIncorrectPassword)
	}

	if !user.IsVerified {
		a.logger.WithField("email", email).Error("User not verified")
		return nil, "", 0, resperror.NewError(resperror.UserNotVerifiedError)
	}

	// Check if license of user is valid
	if !a.userService.CheckLicenseValid(user) {
		return nil, "", 0, resperror.NewError(resperror.UserLicenseExpiredError)
	}

	// Remove any existing sessions
	var existingSessions []Session
	result := a.db.Find(&existingSessions)
	if result.Error != nil {
		a.logger.WithField("err", err).Error("Error while querying for existing sessions")
		return nil, "", 0, nil
	}
	// Remove sessions from cache
	if result.RowsAffected != 0 {
		for _, currSession := range existingSessions {
			a.logger.WithField("session_token", currSession.Token).
				Info("Removing session token from cache")
			a.sessionCache.Delete(currSession.Token)
		}
	}
	// Remove sessions from DB
	err = a.db.Where("user_id = ?", user.ID).Delete(&Session{}).Error
	if err != nil {
		a.logger.WithField("err", err).Error("Error occurred while removing existing sessions for user")
		return nil, "", 0, nil
	}

	// Create new session
	token := uuid.NewString()
	if token == "" {
		return nil, "", 0, errors.New("Failed to generate uuid token")
	}

	expiry := time.Now().Add(time.Hour * 24 * DefaultSessionDurationDay)
	newSession := &Session{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: expiry,
	}
	err = a.db.Create(newSession).Error
	if err != nil {
		a.logger.WithField("err", err).Error("Failed to store session token")
		return nil, "", 0, err
	}

	// Store session in cache
	newSession.User = *user
	a.sessionCache.Set(newSession.Token, *newSession)
	a.logger.WithFields(logrus.Fields{
		"session_token": newSession.Token,
		"email":         newSession.User.Email,
	}).Info("Stored session in cache")

	return user, token, expiry.Unix(), nil
}

// Retrieves session user from token. If session is expired, deletes session and returns error
func (a *SessionService) GetSessionByToken(token string) (*user.User, error) {
	a.logger.WithField("token", token).Info("Getting session with token")

	session, err := a.getSession(token)
	if err != nil {
		a.logger.WithField("err", err).Error("Failed to get session")
		return nil, err
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		a.logger.WithFields(logrus.Fields{
			"token":      token,
			"expires_at": session.ExpiresAt,
		}).Error("Session is expired")

		err := a.DeleteSession(session)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Session expired")
	}

	return &session.User, nil
}

// Retrieve session object (first queries cache, then on cache miss, query DB)
func (a *SessionService) getSession(token string) (Session, error) {
	cachedSession, err := a.sessionCache.Get(token)
	if err == nil {
		a.logger.WithFields(logrus.Fields{
			"user_email": cachedSession.User.Email,
			"token":      cachedSession.Token,
		}).Info("Session cache hit")
		return cachedSession, nil
	}
	// Cache miss, query DB
	a.logger.WithField("token", token).Info("Cache miss, querying DB for session")
	var session Session
	err = a.db.Where("token = ?", token).Preload("User").First(&session).Error
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

// Delete session from cache and DB
func (a *SessionService) DeleteSession(session Session) error {
	a.sessionCache.Delete(session.Token)
	return a.deleteDBSessionByID(session.ID)
}

func (a *SessionService) deleteDBSessionByID(id uint) error {
	return a.db.Delete(&Session{ID: id}).Error
}
