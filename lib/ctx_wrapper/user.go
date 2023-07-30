package ctxwrapper

import (
	"context"
	"errors"

	"github.com/dominiclet/golang-base/service/user"
	"github.com/gin-gonic/gin"
)

const userKey = "user"

func SetUser(c *gin.Context, user user.User) {
	c.Set(userKey, user)
}

// Gets user model instance from context
// NOTE: User model is only injected in protected endpoints
// so make sure usage can only be reached from protected handlers
func GetUser(ctx context.Context) (user.User, error) {
	v := ctx.Value(userKey)
	if v == nil {
		return user.User{}, errors.New("User object not found in context")
	}
	if user, ok := v.(user.User); ok {
		return user, nil
	}
	return user.User{}, errors.New("Unknown object stored as user in context")
}
