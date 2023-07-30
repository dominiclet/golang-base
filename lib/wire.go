package lib

import (
	"github.com/dominiclet/golang-base/lib/email"
	"github.com/google/wire"
)

var LibSet = wire.NewSet(email.InitEmailService)
