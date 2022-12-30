package admin

import (
	"github.com/bars-squad/ais-user-query-service/jwt"
	"github.com/sirupsen/logrus"
)

type Property struct {
	ServiceName  string
	Logger       *logrus.Logger
	Repository   Repository
	JSONWebToken *jwt.JSONWebToken
	// Session                  session.Session
	// Publisher                pubsub.Publisher
}
